package user

import (
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"cheatppt/config"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var secret []byte // HMAC secret
var onceConf sync.Once

func tokenGenerate(username string) (string, error) {
	now := time.Now()
	expire := time.Now().Add(14 * 24 * time.Hour)

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
		return "", err
	}

	return tokenString, nil
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

func ValidToken(token string) bool {
	now := time.Now()

	claims := tokenParse(token)
	if claims == nil || claims.ExpiresAt == nil || claims.ExpiresAt.Before(now) {
		return false
	}

	return true
}
