package model

import (
	"cheatppt/log"
	"cheatppt/model/sql"
	"fmt"
)

type AddDetail struct {
	Id          uint
	DisplayName string
	ModelName   string
	Provider    string
	LeastCoins  int
	Activated   bool
}

func Add(detail *AddDetail, create bool) error {
	var err error
	model := sql.Model{
		DisplayName: detail.DisplayName,
		ModelName:   detail.ModelName,
		Provider:    detail.Provider,
		LeastCoins:  detail.LeastCoins,
		Activated:   detail.Activated,
	}

	db := sql.NewSQLClient()
	if create {
		model.ID = detail.Id
		err = db.Model(&sql.Model{}).Create(&model).Error
	} else {
		err = db.Model(&sql.Model{}).Updates(&model).Error
	}

	if err != nil {
		log.Errorf("MODEL Add ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}
