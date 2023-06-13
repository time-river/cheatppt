package sql

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
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

type UserRechargeRecord struct {
	ID            uint    `gorm:"primaryKey"`
	UserID        uint    `gorm:"index"`
	Amount        float64 `gorm:"type:decimal(10,2)"`
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
}
