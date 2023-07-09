package user

import (
	"github.com/kr/pretty"
	log "github.com/sirupsen/logrus"

	"cheatppt/model/sql"
)

func CreateSession(userId int) (any, error) {
	db := sql.NewSQLClient()

	session := sql.ChatSession{
		UserID: uint(userId),
	}

	db.Save(&session)

	log.Trace(pretty.Sprint(session))
	return session, nil
}
