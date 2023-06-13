package sql

import "time"

type Model struct {
	ID          uint   `gorm:"primaryKey"`
	DisplayName string `gorm:"unique"`
	ModelName   string
	Provider    string `gorm:"index"`
	LeastCoins  int
	Comment     string `gorm:"type:text"`
	Activated   bool
	CreatedAt   time.Time
}

type ModelSwitchMapping struct {
	ID          uint  `gorm:"primaryKey"`
	FromModelID uint  `gorm:"index"`
	ToModelID   uint  `gorm:"index"`
	FromModel   Model `gorm:"foreignKey:FromModelID"`
	ToModel     Model `gorm:"foreignKey:ToModelID"`
}

func modelTableInit() {
	db := NewSQLClient()
	if err := db.AutoMigrate(&Model{}); err != nil {
		panic(err)
	}
}
