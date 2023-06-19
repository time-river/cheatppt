package user

import (
	"fmt"

	"cheatppt/model/sql"
)

type Information struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Level     int    `json:"level"`
	Coins     int    `json:"coins"`
	CreatedAt int64  `json:"createAt"`
}

func Detail(id int) (*Information, error) {
	var user sql.User
	db := sql.NewSQLClient()

	if err := db.First(&user, id).Error; err != nil {
		return nil, fmt.Errorf("内部错误")
	}
	if !user.Activated || user.Deleted {
		return nil, fmt.Errorf("用户不存在")
	}

	info := Information{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Level:     user.Level,
		Coins:     user.Coins,
		CreatedAt: user.CreatedAt.Unix(),
	}

	return &info, nil
}
