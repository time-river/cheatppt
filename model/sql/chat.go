package sql

import "time"

type ChatSession struct {
	ID     uint `gorm:"primaryKey"`
	UserID uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserID"`
}

type ChatMessage struct {
	ID               uint `gorm:"primaryKey"`
	ChatSessionID    uint `gorm:"index"`
	ModelDisplayName string
	Sender           string
	Message          string      `gorm:"type:text"`
	CreatedAt        time.Time   `gorm:"index"`
	ChatSession      ChatSession `gorm:"foreignKey:ChatSessionID"`
	Model            Model       `gorm:"foreignKey:ModelDisplayName"`
	Deleted          bool
}
