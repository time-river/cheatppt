package user

import "sync"

func Setup() {
	cacheUserMgr = &cacheUserManager{
		users: make(map[uint]*CacheUser),
		mu:    &sync.RWMutex{},
	}
}
