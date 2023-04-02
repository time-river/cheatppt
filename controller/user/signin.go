package user

import (
	"bytes"
	"errors"
	"fmt"

	"cheatppt/controller/chat"
	"cheatppt/log"
	"cheatppt/model/redis"
	"cheatppt/model/sql"
	"cheatppt/utils"

	"github.com/kr/pretty"
)

type SignInData struct {
	Email        string
	Token        string
	ModelSetting chat.ModelSetting
}

func SignIn(username, passwd string) (*SignInData, error) {
	var user sql.User

	db := sql.NewSQLClient()
	err := db.Model(&sql.User{}).Where("username = ?", username).First(&user).Error
	if err != nil {
		log.Error(err.Error())
		return nil, fmt.Errorf("用户名或密码错误")
	}

	digest := utils.Digest(passwd)
	if !bytes.Equal(user.Password, digest) {
		log.Debug("wrong password")
		return nil, fmt.Errorf("用户名或密码错误")
	}

	token, err := tokenGenerate(username)
	if err != nil {
		log.Warnf("tokenGenerate ERROR: %s\n", err.Error())
		return nil, errors.New("内部错误")
	}

	rds := redis.NewRedisCient()
	if err := rds.TokenLease(token, username); err != nil {
		log.Errorf("TokenLease ERROR: %s\n", err.Error())
		return nil, errors.New("内部错误")
	}

	modelSetting := chat.GetModelSetting(user.Level)
	data := &SignInData{
		Email:        user.Email,
		Token:        token,
		ModelSetting: modelSetting,
	}

	log.Debug(pretty.Sprint(data))

	return data, nil
}
