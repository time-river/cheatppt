package model

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"cheatppt/model/sql"
)

type AddDetail struct {
	Id          uint
	DisplayName string
	ModelName   string
	Provider    string
	InputCoins  int
	OutputCoins int
	Activated   bool
}

func Add(detail *AddDetail, create bool) error {
	var err error
	model := sql.Model{
		DisplayName: detail.DisplayName,
		ModelName:   detail.ModelName,
		Provider:    detail.Provider,
		InputCoins:  detail.InputCoins,
		OutputCoins: detail.OutputCoins,
		Activated:   detail.Activated,
	}

	db := sql.NewSQLClient()
	if create {
		err = db.Model(&sql.Model{}).Create(&model).Error
	} else {
		err = db.Model(&sql.Model{}).Where("id = ?", detail.Id).Updates(&model).Error
	}

	if err != nil {
		log.Errorf("MODEL Add ERROR: %s\n", err.Error())
		return fmt.Errorf("内部错误")
	}

	return nil
}
