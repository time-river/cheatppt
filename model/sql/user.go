package sql

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey,AUTO_INCREMENT=100"`
	Username  string `gorm:"uniqueIndex"`
	Email     string `gorm:"uniqueIndex"`
	Password  []byte
	Level     int
	Coins     int
	Activated bool
	Deleted   bool
	CreatedAt time.Time
	DeletedAt time.Time
}

type UserVoiceRecord struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint `gorm:"index"`
	Amount        int
	Coins         int
	CreatedAt     time.Time
	PaymentMethod string `gorm:"TINYTEXT"`
	User          User   `gorm:"foreignKey:UserID"`
}

type DailyFree struct {
	ID     uint      `gorm:"primaryKey"`
	UserID uint      `gorm:"index"`
	Date   time.Time `gorm:"index"`
	Coins  int
	User   User `gorm:"foreignKey:UserID"`
}

func userTableInit() {
	db := NewSQLClient()
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&UserVoiceRecord{}); err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&DailyFree{}); err != nil {
		panic(err)
	}
}
