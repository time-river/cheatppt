package auth

import (
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	secret []byte // HMAC secret
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

type Digest struct {
	salt string
}

func (t *Token) generate(username *string) (*string, error) {
	now := time.Now()
	// max allow time is one weak
	expireTime := now.Add(7 * 24 * time.Hour)

	claims := &Claims{
		Username: *username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(t.secret)
	if err != nil {
		return nil, err
	}

	return &tokenString, nil
}

/*
func (t *Token) parse(tokenString *string) *Claims {
	token, _ := jwt.ParseWithClaims(*tokenString, &Claims{},
		func(token *jwt.Token) (interface{}, error) {
			return t.secret, nil
		})

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		// TODO: check time
		return claims
	}

	return nil
}
*/

func (d *Digest) digest(text *string) []byte {
	plain := *text + d.salt
	result := sha3.Sum256([]byte(plain))

	secret := make([]byte, len(result))
	copy(secret, result[:])

	return secret
}
