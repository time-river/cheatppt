package sql

import (
	"time"
)

type ChatGPTAccount struct {
	ID              uint `gorm:"primarykey"`
	Activated       bool
	Email           string `gorm:"uniqueIndex"`
	Password        []byte
	IV              []byte
	CreatedAt       time.Time
	UpdatedAt       time.Time
	AccessToken     []byte
	AccessTokenRfAt time.Time
	Puid            []byte
	PuidRfAt        time.Time
}

type ChatGPTConversationMapping struct {
	ID             uint   `gorm:"primarykey"`
	ConversationId string `gorm:"uniqueIndex"`
	AccountEmail   string
	ChatGPTAccount ChatGPTAccount `gorm:"foreignKey:AccountEmail"`
}

func initChatGPT() {
	db := NewSQLClient()

	if err := db.AutoMigrate(&ChatGPTAccount{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate((&ChatGPTConversationMapping{})); err != nil {
		panic(err)
	}
}
