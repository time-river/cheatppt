package user

import (
	"cheatppt/config"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var secret []byte // HMAC secret
var onceConf sync.Once

func tokenGenerate(username string) (string, error) {
	now := time.Now()

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
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
		return "", err
	}

	return tokenString, nil
}

func tokenParse(tokenString string) *string {
	token, _ := jwt.ParseWithClaims(tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

	if claims, ok := token.Claims.(*Claims); ok {
		return &claims.Username
	}
	return nil
}
