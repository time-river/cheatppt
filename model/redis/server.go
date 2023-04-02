package redis

import (
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"cheatppt/config"
)

type Redis struct {
	client *redis.Client
	lease  time.Duration
}

var onceConf sync.Once
var rds *Redis

func RedisCtxCreate() *Redis {
	conf := config.GlobalCfg.Redis
	options := &redis.Options{
		Addr:     conf.Addr,
		Password: conf.Passwd,
		DB:       conf.Db,
	}
	rdb := redis.NewClient(options)

	if rds == nil {
		onceConf.Do(func() {
			rds = &Redis{
				client: rdb,
				lease:  time.Duration(60 * int64(time.Second)), // TODO: seconds?
			}
			// TODO: connect test
		})
	}
	return rds
}
