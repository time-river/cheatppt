package sql

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"UniqueIndex"`
	Email     string
	Password  []byte
	Activated bool
	Deleted   bool
	Level     int
	VipEnd    time.Time
}

func userTableInit() {
	db := NewSQLClient()
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}
