package db

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username      string `gorm:"UniqueIndex"`
	Email         string
	Password      []byte
	Level         int
	EmailVerified bool
}
