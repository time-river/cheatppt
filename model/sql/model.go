package sql

import "time"

type Model struct {
	// TODO: ensure id is from 100
	ID          uint   `gorm:"primaryKey,AUTO_INCREMENT=100"`
	DisplayName string `gorm:"unique"`
	ModelName   string `gorm:"index"`
	Provider    string `gorm:"index"`
	InputCoins  int    // virtual coins
	OutputCoins int    // virtual coins
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

	if err := db.AutoMigrate(&ModelSwitchMapping{}); err != nil {
		panic(err)
	}
}
