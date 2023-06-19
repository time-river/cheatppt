package model

import (
	log "github.com/sirupsen/logrus"

	"cheatppt/model/sql"
)

func Find(model, provider string) *Model {
	result := sql.Model{
		DisplayName: model,
		Provider:    provider,
	}

	db := sql.NewSQLClient()
	if err := db.First(&result).Error; err != nil {
		log.Errorf("MODEL find ERROR: %s\n", err.Error())
		return nil
	}

	return &Model{
		Id:          int(result.ID),
		DisplayName: result.DisplayName,
		ModelName:   result.ModelName,
		Provider:    result.Provider,
		InputCoins:  result.InputCoins,
		OutputCoins: result.OutputCoins,
		Activated:   result.Activated,
	}
}
