package user

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"cheatppt/controller/mail"
	"cheatppt/model/redis"
	"cheatppt/model/sql"
)

// triplet: (username, email, code)
const codeValidMin = 30 // validity period, minutes
const codeLength = 6

func generateCode(length int) string {
	rand.Seed(time.Now().UnixNano())

	const lettersAndDigits = `abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789`
	code := make([]byte, length)
	for i := range code {
		code[i] = lettersAndDigits[rand.Intn(len(lettersAndDigits))]
	}

	return string(code)
}

func GenerateResetCode(username string) error {
	var user sql.User
	db := sql.NewSQLClient()

	if result := db.Model(&sql.User{}).Where("username = ?", username).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// user doesn't exist, but ignore error, it will make the user insensible
			return nil
		}

		return fmt.Errorf("内部错误")
	} else if !user.Activated {
		// user is creating (not activate), but not reports error
		return nil
	}

	mailCtx := mail.CodeCtx{
		Username: username,
		Email:    user.Email,
		Code:     generateCode(codeLength),
		ValidMin: codeValidMin,
	}

	key := fmt.Sprintf("%s-reset", mailCtx.Username)
	val := mailCtx.Code
	if err := sendCode(&mailCtx, key, val); err != nil {
		return fmt.Errorf("内部错误")
	}

	return nil
}

var dbMutex sync.Mutex

func prepareUser(username, email string) error {

	user := sql.User{
		Username:  username,
		Email:     email,
		Activated: false,
	}

	db := sql.NewSQLClient()

	result := db.Create(&user)
	if result.Error == nil {
		return nil
	}
	// TODO: other errors?

	// the lock ensures two SQL are atomic
	dbMutex.Lock()
	defer dbMutex.Unlock()

	err := db.Model(&sql.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Error(err.Error())
		return fmt.Errorf("内部错误")
	} else if user.Activated || ((time.Since(user.CreatedAt) < 2*time.Hour) && user.Email != email) {
		// user has been created and activated, or created but not activated
		return fmt.Errorf("用户已存在")
	}

	info := map[string]interface{}{
		"email":      email,
		"created_at": time.Now(),
	}
	if err := db.Model(&sql.User{}).Where("username = ?", username).Updates(info).Error; err != nil {
		log.Error(err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}

func GenerateSignUpCode(username, email string) error {
	if err := prepareUser(username, email); err != nil {
		return err
	}

	mailCtx := mail.CodeCtx{
		Username: username,
		Email:    email,
		Code:     generateCode(codeLength),
		ValidMin: codeValidMin,
	}

	/* only one valid code exists for one user */
	key := fmt.Sprintf("%s-signup", mailCtx.Username)
	value := fmt.Sprintf("%s %s", mailCtx.Code, mailCtx.Email)
	if err := sendCode(&mailCtx, key, value); err != nil {
		return fmt.Errorf("内部错误")
	}

	return nil
}

func sendCode(ctx *mail.CodeCtx, key, value string) error {

	rds := redis.NewRedisCient()
	if err := rds.SetCode(key, value, ctx.ValidMin); err != nil {
		log.Errorf("SetCode ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	if err := mail.SendCode(ctx); err != nil {
		log.Errorf("SendMail ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}

func validateCode(key, code string) (bool, error) {
	rds := redis.NewRedisCient()
	value, err := rds.GetCode(key)
	if err != nil {
		log.Errorf("GetCode ERROR: %s\n", err.Error())
		return false, fmt.Errorf("内部错误")
	}

	rds.DelCode(key)

	return value == code, nil
}

func ValidateSignUpCode(username, email, code string) (bool, error) {
	key := fmt.Sprintf("%s-signup", username)
	val := fmt.Sprintf("%s %s", code, email)
	valid, err := validateCode(key, val)
	if err != nil {
		return false, err
	} else if valid {
		db := sql.NewSQLClient()

		if err := db.Model(&sql.User{}).Where("username = ?", username).Update("activated", true).Error; err != nil {
			log.Errorf("ValidateSignUpCode ERROR: %s", err.Error())
			return false, fmt.Errorf("内部错误")
		}
	}

	return valid, nil
}

func ValidateResetCode(username, code string) (bool, error) {
	key := fmt.Sprintf("%s-reset", username)
	val := code
	valid, err := validateCode(key, val)
	if err != nil {
		return false, err
	}

	return valid, nil
}
