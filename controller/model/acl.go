package model

import (
	log "github.com/sirupsen/logrus"

	"cheatppt/controller/user"
	"cheatppt/model/sql"
)

type AclDetail struct {
	Username  uint
	ModelName string
	Provider  string
}

func Allow(coins int, token string) bool {
	claims := user.TokenParse(token)
	if claims == nil {
		return false
	}

	var result sql.User
	db := sql.NewSQLClient()
	err := db.First(&result, claims.ID).Error
	if err != nil {
		log.Errorf("MODEL Allow First ERROR: %s\n", err.Error())
		return false
	}

	return true
}
