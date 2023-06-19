package user

import (
	"encoding/base64"
	"math/rand"

	log "github.com/sirupsen/logrus"

	"cheatppt/model/redis"
)

const sessionValidHour = 12 // validity period, hours

func NewSession(token string) *string {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		log.Errorf("SESSION NewSession rand ERROR: %s\n", err.Error())
		return nil
	}

	session := base64.URLEncoding.EncodeToString(b)

	rds := redis.NewRedisCient()
	if err := rds.SetSession(session, token, sessionValidHour); err != nil {
		log.Errorf("SESSION NewSession ERROR: %s\n", err.Error())
		return nil
	}

	return &session
}

func ValidSession(session string) (string, bool) {
	rds := redis.NewRedisCient()
	token, err := rds.ValidThenRenewSession(session, sessionValidHour)
	if token != nil && err != nil {
		log.Errorf("SESSION Expire ERROR: %s\n", err.Error())
	}

	if token == nil {
		return "", false
	} else {
		return *token, true
	}
}
