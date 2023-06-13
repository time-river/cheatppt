package user

import (
	"bytes"
	"fmt"

	"github.com/kr/pretty"

	"cheatppt/log"
	"cheatppt/model/sql"
	"cheatppt/utils"
)

type SignInData struct {
	Email string
	Token string
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

	token, err := newToken(username)
	if err != nil {
		return nil, fmt.Errorf("内部错误")
	}

	data := &SignInData{
		Email: user.Email,
		Token: *token,
	}

	log.Trace(pretty.Sprint(data))

	return data, nil
}
