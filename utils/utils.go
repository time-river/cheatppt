package utils

import (
	"cheatppt/config"
	"regexp"

	"golang.org/x/crypto/sha3"
)

func Must[T interface{}](v T, err error) T {
	if err != nil {
		panic(err)
	}

	return v
}

func Digest(text string) []byte {
	salt := config.Server.Secret
	plain := text + salt
	result := sha3.Sum256([]byte(plain))

	secret := make([]byte, len(result))
	copy(secret, result[:])

	return secret
}

func UsernameCheck(username string) bool {
	const pattern = `^\S{4,16}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(username)
}

func EmailCheck(email string) bool {
	const pattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(email)
}

func PasswordCheck(passwd string) bool {
	const pattern = `^.{6,}$`
	regex := regexp.MustCompile(pattern)

	return regex.MatchString(passwd)
}
