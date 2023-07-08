package model

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"cheatppt/model/sql"
)

type DelDetail struct {
	Id uint
}

func Del(detail *DelDetail) error {
	model := sql.Model{
		ID: detail.Id,
	}

	db := sql.NewSQLClient()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&sql.Model{}).First(&model, detail.Id).Error; err != nil {
			return err
		}

		CacheDel(BuildCacheKey(model.Provider, model.ModelName))

		return tx.Delete(&model).Error
	})
	if err != nil {
		log.Errorf("MODEL Add ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}
