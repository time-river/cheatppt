package sql

import (
	"sync"

	"gorm.io/gorm"
)

type Sql struct {
	db *gorm.DB
}

var onceConf sync.Once
var sql *Sql

func SQLCtxCreate() *Sql {
	if sql == nil {
		onceConf.Do(func() {
			sql = &Sql{
				db: dbConnect(),
			}
		})
	}

	return sql
}
