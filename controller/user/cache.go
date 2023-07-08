package user

import (
	"fmt"
	"sync"

	"cheatppt/model/sql"
)

type CacheUser struct {
	ID       uint
	Username string
	Email    string
	Level    int
	Coins    int64
}

// TODO: size
type cacheUserManager struct {
	users map[uint]*CacheUser

	mu *sync.RWMutex
}

var cacheUserMgr *cacheUserManager

func (cm *cacheUserManager) add(user *sql.User) *CacheUser {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	if val, ok := cm.users[user.ID]; !ok {
		val = &CacheUser{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Level:    user.Level,
			Coins:    user.Coins,
		}
		cm.users[user.ID] = val

		return val
	} else {
		return val
	}
}

func (cm *cacheUserManager) put(userId uint) {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	delete(cm.users, userId)
}

func (cm *cacheUserManager) get(userId uint) *CacheUser {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	if user, ok := cm.users[userId]; ok {
		return user
	} else {
		return nil
	}
}

func CacheFind(userId uint) (*CacheUser, error) {
	cacheUser := cacheUserMgr.get(userId)
	if cacheUser == nil {
		var user sql.User
		db := sql.NewSQLClient()

		if err := db.First(&user, userId).Error; err != nil {
			return nil, fmt.Errorf("内部错误")
		} else if !user.Activated {
			return nil, fmt.Errorf("用户不存在")
		}

		return cacheUserMgr.add(&user), nil
	} else {
		return cacheUser, nil
	}
}
