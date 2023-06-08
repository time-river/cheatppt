package user

import (
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"cheatppt/config"
	"cheatppt/log"
	"cheatppt/model/redis"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

const tokenValidHour = 14 * 24 // validity period, hours

var secret []byte // HMAC secret
var onceConf sync.Once

func tokenGenerate(username string) (*string, error) {
	now := time.Now()
	expire := time.Now().Add(tokenValidHour * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expire),
		},
	}

	if len(secret) == 0 {
		onceConf.Do(func() {
			secret = []byte(config.Server.Secret)
		})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Errorf("TOKEN SignedString ERROR: %s\n", err.Error())
		return nil, err
	}

	return &tokenString, nil
}

func tokenParse(tokenString string) *Claims {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

	if claims, ok := token.Claims.(*Claims); ok {
		return claims
	}
	return nil
}

func newToken(username string) (*string, error) {
	token, err := tokenGenerate(username)
	if err != nil {
		return nil, err
	}

	rds := redis.NewRedisCient()
	if err := rds.SetToken(*token, username, tokenValidHour); err != nil {
		log.Errorf("TOKEN SetToken ERROR: %s\n", err.Error())
		return nil, err
	}

	return token, nil
}

func ValidToken(token string) (bool, error) {
	now := time.Now()

	claims := tokenParse(token)
	if claims == nil || claims.ExpiresAt == nil || claims.ExpiresAt.Before(now) {
		return false, nil
	}

	rds := redis.NewRedisCient()
	exist, err := rds.ExistsToken(token)
	if err != nil {
		log.Errorf("TOKEN ExistsToken ERROR: %s\n", err.Error())
		return false, fmt.Errorf("内部错误")
	}

	return exist > 0, nil
}
