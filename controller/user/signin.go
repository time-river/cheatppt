package user

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/kr/pretty"

	"cheatppt/controller/chat/model"
	"cheatppt/log"
	"cheatppt/model/redis"
	"cheatppt/model/sql"
	"cheatppt/utils"
)

type SignInData struct {
	Email        string
	Token        string
	ModelSetting model.ModelSetting
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
		log.Trace("wrong password")
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

	modelSetting := model.GetModelSetting(user.Level)
	data := &SignInData{
		Email:        user.Email,
		Token:        token,
		ModelSetting: modelSetting,
	}

	log.Trace(pretty.Sprint(data))

	return data, nil
}
