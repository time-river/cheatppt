package charge

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"cheatppt/model/sql"
)

type Consumer struct {
	UserId    int
	UserLevel int
	Coins     int
	Free      bool
}

func Comsume(userId int, price int, force bool) (*Consumer, error) {
	var user sql.User
	free := true

	db := sql.NewSQLClient()
	err := db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userId).Error; err != nil {
			return err
		} else if !user.Activated {
			return fmt.Errorf("用户不存在")
		}

		if !force && user.Coins <= 0 {
			return nil
		}

		user.Coins -= price
		if err := tx.Save(&user).Error; err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("CHARGE Comsume ERROR: %s\n", err.Error())
		return nil, fmt.Errorf("内部错误")
	}

	return &Consumer{
		UserId:    userId,
		UserLevel: user.Level,
		Coins:     user.Coins,
		Free:      free,
	}, nil
}
