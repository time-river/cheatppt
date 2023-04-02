package sql

import (
	"cheatppt/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func dbConnect() *gorm.DB {
	conf := config.GlobalCfg.DB
	db, err := gorm.Open(sqlite.Open(conf.Addr), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	if _, err := db.DB(); err != nil {
		panic(err)
	}
	return db
}
