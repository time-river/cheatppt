package model

import (
	"fmt"

	log "github.com/sirupsen/logrus"

	"cheatppt/model/sql"
)

type Model struct {
	Id          int    `json:"id"`
	DisplayName string `json:"displayName"`
	ModelName   string `json:"modelName"`
	Provider    string `json:"provider"`
	InputCoins  int    `json:"inputCoins"`
	OutputCoins int    `json:"outputCoins"`
	Activated   bool   `json:"activated"`
	Comment     string `json:"comment,omitempty"`
	CreatedAt   int64  `json:"createAt,omitempty"`
}

func ListAvailable() ([]Model, error) {
	var models []sql.Model
	var data = make([]Model, 0)

	db := sql.NewSQLClient()
	result := db.Model(&sql.Model{}).Find(&models)
	if result.Error != nil {
		log.Errorf("MODEL ListAvailable ERROR: %s\n", result.Error.Error())
		return nil, fmt.Errorf("内部错误")
	}

	for _, model := range models {
		if !model.Activated {
			continue
		}

		m := Model{
			Id:          int(model.ID),
			DisplayName: model.DisplayName,
			ModelName:   model.ModelName,
			Provider:    model.Provider,
			InputCoins:  model.InputCoins,
			OutputCoins: model.OutputCoins,
			Activated:   model.Activated,
		}
		data = append(data, m)
	}

	return data, nil
}

func ListAll() ([]Model, error) {
	var models []sql.Model
	var data = make([]Model, 0)

	db := sql.NewSQLClient()
	result := db.Model(&sql.Model{}).Find(&models)
	if result.Error != nil {
		log.Errorf("MODEL ListAll ERROR: %s\n", result.Error.Error())
		return nil, fmt.Errorf("内部错误")
	}

	for _, model := range models {
		m := Model{
			Id:          int(model.ID),
			DisplayName: model.DisplayName,
			ModelName:   model.ModelName,
			Provider:    model.Provider,
			InputCoins:  model.InputCoins,
			OutputCoins: model.OutputCoins,
			Activated:   model.Activated,
			Comment:     model.Comment,
			CreatedAt:   model.CreatedAt.Unix(),
		}
		data = append(data, m)
	}

	return data, nil
}
