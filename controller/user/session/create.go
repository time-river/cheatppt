package user

import (
	"fmt"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"

	"cheatppt/model/sql"
)

func CreateSession(userId int) (string, error) {
	db := sql.NewSQLClient()

	id, err := uuid.NewRandom()
	if err != nil {
		log.Warn(err)
		return "", fmt.Errorf("内部错误")
	}

	session := sql.ChatSession{
		UUID:   id,
		UserID: uint(userId),
	}

	if err := db.Save(&session).Error; err != nil {
		log.Warn(err)
		return "", fmt.Errorf("内部错误")
	}

	return id.String(), nil
}
