package redis

import (
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"cheatppt/config"
)

type Redis struct {
	conn  *redis.Conn
	lease time.Duration
}

var onceConf sync.Once
var rds *Redis

func NewRedisCient() *Redis {
	conf := config.GlobalCfg.Redis
	options := &redis.Options{
		Addr:     conf.Addr,
		Password: conf.Passwd,
		DB:       conf.Db,
	}
	rdb := redis.NewClient(options).Conn()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(err.Error())
	}

	if rds == nil {
		onceConf.Do(func() {
			rds = &Redis{
				conn:  rdb,
				lease: time.Duration(60 * int64(time.Second)),
			}
		})
	}
	return rds
}
