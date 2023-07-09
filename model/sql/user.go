package sql

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uint   `gorm:"primaryKey,AUTO_INCREMENT=100"`
	Username  string `gorm:"uniqueIndex"`
	Email     string `gorm:"uniqueIndex"`
	Password  []byte
	Level     int
	Coins     int64 // virtual coins
	Activated bool
	Deleted   bool
	CreatedAt time.Time
	DeletedAt time.Time
}

type UserBilling struct {
	UUID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID        uint      `gorm:"index"`
	Coins         int64     // virtual coins
	Status        int
	PaymentMethod string `gorm:"TINYTEXT"`
	CreatedAt     time.Time
	PaiedAt       time.Time
	Comment       string
	User          User `gorm:"foreignKey:UserID"`
}

type UserUsage struct {
	UUID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	Coins         int64     // virtual coins
	Comment       string    // the voice detail
	CreatedAt     time.Time
	UserID        uint        `gorm:"index"`
	User          User        `gorm:"foreignKey:UserID"`
	ChatMessageId uuid.UUID   `gorm:"index;type:uuid"`
	ChatMessage   ChatMessage `gorm:"foreignKey:ChatMessageId"`
}

func userTableInit() {
	db := NewSQLClient()
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&UserBilling{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&UserUsage{}); err != nil {
		panic(err)
	}
}
