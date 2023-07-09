package billing

import (
	"sync/atomic"
	"time"

	"cheatppt/controller/user"
	"cheatppt/model/sql"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Consumer struct {
	UserId    int
	UserLevel int
	Free      bool
	Coins     int
	user      *user.CacheUser
}

// price: virtual coins
// force: force comsume (e.g. )
func GetComsumer(userId int) (*Consumer, error) {
	cacheUser, err := user.CacheFind(uint(userId))
	if err != nil {
		return nil, err
	}

	return &Consumer{
		UserId:    userId,
		UserLevel: cacheUser.Level,
		user:      cacheUser,
		Free:      atomic.LoadInt64(&cacheUser.Coins) <= 0,
	}, nil
}

func (c *Consumer) Comsume(price int) {
	if c.Free {
		return
	}

	c.Coins += price
	atomic.AddInt64(&c.user.Coins, int64(-price))
}

func (c *Consumer) Commit(comment string) {
	if c.Free {
		return
	}

	db := sql.NewSQLClient()
	db.Transaction(func(tx *gorm.DB) error {
		var user sql.User

		if err := tx.First(&user, c.UserId).Error; err != nil {
			return err
		}

		user.Coins -= int64(c.Coins)
		tx.Save(&user)

		usage := sql.UserUsage{
			UserID:        c.user.ID,
			ChatMessageId: uuid.Must(uuid.NewRandom()),
			Coins:         int64(c.Coins),
			CreatedAt:     time.Now(),
			Comment:       comment,
		}
		tx.Save(&usage)
		return nil
	})
}

func (c *Consumer) Rollback() {
	atomic.AddInt64(&c.user.Coins, int64(c.Coins))
}
