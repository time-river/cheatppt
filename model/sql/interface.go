package sql

import (
	"sync"

	"gorm.io/gorm"
)

var onceConf sync.Once
var db *gorm.DB

func NewSQLClient() *gorm.DB {
	if db == nil {
		onceConf.Do(func() {
			db = dbConnect()
		})
	}

	return db
}

func DatabaseInit() {
	userTableInit()
	modelTableInit()
}
