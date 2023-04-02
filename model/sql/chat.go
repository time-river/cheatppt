package sql

type Model struct {
	ID         uint   `gorm:"primarykey"`
	DislayName string `gorm:"UniqueIndex"`
	ModelName  string
	Provider   string
	Valid      bool
}

func modelTableInit() {
	db := NewSQLClient()
	if err := db.AutoMigrate(&User{}); err != nil {
		panic(err)
	}
}
