package user

import (
	"fmt"

	"cheatppt/log"
	"cheatppt/model/sql"
	"cheatppt/utils"
)

func ResetPassword(username, passwd string) error {
	db := sql.NewSQLClient()

	blob := utils.Digest(passwd)
	if err := db.Model(&sql.User{}).Where("username = ?", username).Update("password", blob).Error; err != nil {
		log.Tracef("ResetPassword ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}
