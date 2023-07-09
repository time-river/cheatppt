package billing

import (
	"fmt"

	"cheatppt/model/sql"
)

type UsageItem struct {
	Amount  float32 `json:"coins"`
	Date    int64   `json:"date"`
	Comment string  `json:"comment"`
}

func GetUsages(rng *Range) ([]UsageItem, error) {
	var records []sql.UserUsage
	var data = make([]UsageItem, rng.Length)

	db := sql.NewSQLClient()

	if err := db.Offset(rng.From).Limit(rng.Length).Find(&records).Error; err != nil {
		return nil, fmt.Errorf("内部错误")
	}

	for i, r := range records {
		data[i] = UsageItem{
			Amount:  sql.Coins2RMB(r.Coins),
			Date:    r.CreatedAt.Unix(),
			Comment: r.Comment,
		}
	}

	return data[0:], nil
}
