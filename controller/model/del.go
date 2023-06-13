package model

import (
	"fmt"

	"cheatppt/log"
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
	err := db.Delete(&model).Error
	if err != nil {
		log.Errorf("MODEL Add ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}
