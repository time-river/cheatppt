package user

import (
	"fmt"

	"cheatppt/log"
	"cheatppt/model/sql"
	"cheatppt/utils"
)

func SignUp(username, passwd string) error {
	db := sql.NewSQLClient()

	info := map[string]interface{}{
		"password": utils.Digest(passwd),
		"level":    1000,
		"deleted":  false,
	}

	err := db.Model(&sql.User{}).Where("username = ?", username).Updates(info).Error
	if err != nil {
		log.Errorf("SignUp ERROR: %s", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}
