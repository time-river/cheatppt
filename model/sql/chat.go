package sql

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ChatSession struct {
	UUID   uuid.UUID `gorm:"primaryKey;type:uuid"`
	UserID uint      `gorm:"index"`
	User   User      `gorm:"foreignKey:UserID"`
}

type ChatMessage struct {
	UUID          uuid.UUID `gorm:"primaryKey;type:uuid"`
	ModelName     string
	Provider      string
	Messages      datatypes.JSON
	CreatedAt     time.Time
	UpdatedAt     time.Time
	ChatSessionId uuid.UUID   `gorm:"type:uuid"`
	ChatSession   ChatSession `gorm:"foreignKey:ChatSessionId"`
	Deleted       bool
}

func chatTableInit() {
	db := NewSQLClient()

	if err := db.AutoMigrate(&ChatSession{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&ChatMessage{}); err != nil {
		panic(err)
	}
}
