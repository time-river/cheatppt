package sql

import "time"

type ChatSession struct {
	ID     uint `gorm:"primaryKey,AUTO_INCREMENT=100"`
	UserID uint `gorm:"index"`
	User   User `gorm:"foreignKey:UserID"`
}

type ChatMessage struct {
	ID               uint `gorm:"primaryKey,AUTO_INCREMENT=100"`
	ChatSessionID    uint `gorm:"index"`
	ModelDisplayName string
	Prompt           string      `gorm:"type:text"`
	Message          string      `gorm:"type:text"`
	CreatedAt        time.Time   `gorm:"index"`
	ChatSession      ChatSession `gorm:"foreignKey:ChatSessionID"`
	Model            Model       `gorm:"foreignKey:ModelDisplayName"`
	Deleted          bool
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
