package model

import (
	"cheatppt/model/sql"
	"sync"
)

// collect models from SQL in setup procedure
func Setup() {
	cacheModelMgr = &cacheModelManager{
		models: make(map[string]*cacheModel),
		mu:     &sync.RWMutex{},
	}

	var models []sql.Model
	db := sql.NewSQLClient()

	if err := db.Find(&models).Error; err != nil {
		panic(err)
	}

	for _, model := range models {
		CacheAdd(&model)
	}
}
