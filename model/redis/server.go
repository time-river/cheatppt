package redis

import (
	"runtime"
	"sync"

	"github.com/redis/go-redis/v9"

	"cheatppt/config"
)

type Redis struct {
	conn *redis.Client
}

var onceConf sync.Once
var rds *Redis

func NewRedisCient() *Redis {
	conf := config.GlobalCfg.Redis
	options := &redis.Options{
		Addr:     conf.Addr,
		Password: conf.Passwd,
		DB:       conf.Db,
		// 59 connections per every available CPU.
		PoolSize: 50 * runtime.NumCPU(),
	}
	rdb := redis.NewClient(options)

	if rds == nil {
		onceConf.Do(func() {
			rds = &Redis{
				conn: rdb,
			}
		})
	}
	return rds
}
