package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()
var zsetLock = sync.Mutex{}

/*
 * Once get, refresh expire time immediately.
 */

func (rdb *Redis) sessionCreate(key string, member string) error {
	zsetLock.Lock()
	defer zsetLock.Unlock()

	now := time.Now().Unix()
	val := redis.Z{
		Score:  float64(now),
		Member: member,
	}

	num, err := rdb.client.ZCard(ctx, key).Result()
	if err != nil {
		return err
	}

	if num >= 3 {
		pop := num - 3
		rdb.client.ZPopMin(ctx, key, pop)
	}

	/* ignore, all of adding new or refreshing score will be regarded as success */
	if err := rdb.client.ZAdd(ctx, key, val).Err(); err != nil {
		return err
	}

	return nil
}

func (rdb *Redis) sessionPop(key string, member string) {
	zsetLock.Lock()
	defer zsetLock.Unlock()

	rdb.client.ZRem(ctx, key, member)
}

func (rdb *Redis) sessionRefresh(key string, member string) error {
	zsetLock.Lock()
	defer zsetLock.Unlock()

	now := time.Now().Unix()
	val := redis.Z{
		Score:  float64(now),
		Member: member,
	}

	num, err := rdb.client.ZAdd(ctx, key, val).Result()
	if err != nil {
		return err
	} else if num != 0 {
		rdb.client.ZRem(ctx, key, member)
		return errors.New("Expect refresh member score but add member")
	}

	return nil
}

func (rdb *Redis) TokenLease(token string, username string) error {
	if err := rdb.sessionCreate(username, token); err != nil {
		return err
	}

	return rdb.client.Set(ctx, token, username, rdb.lease).Err()
}

func (rdb *Redis) TokenRevoke(token string) {
	/*
	 * 1. Get(token): none stands for the token has been expired
	 * 2. Del(token)
	 * 3. pop ZSet(token)
	 */
	username, err := rdb.client.Get(ctx, token).Result()
	if err != nil {
		return
	}

	err = rdb.client.Del(ctx, token).Err()
	if err != nil {
		// TODO: warning here
	}
	rdb.sessionPop(username, token)
}

func (rdb *Redis) TokenVerify(token string) (*string, error) {
	/*
	 * 1. Get(token)
	 * 2. refresh ZSet(token) timestamp
	 * 3. Reset Expire(token) time
	 */
	username, err := rdb.client.Get(ctx, token).Result()
	if err != nil {
		return nil, err
	}

	if err := rdb.sessionRefresh(username, token); err != nil {
		return nil, err
	}

	done, err := rdb.client.Expire(ctx, token, rdb.lease).Result()
	if err != nil {
		/* EXPIRE failed, but still effective
		 * TODO: warning here
		 */
	} else if !done {
		/* no one, add it */
		if err := rdb.client.Set(ctx, token, username, rdb.lease).Err(); err != nil {
			return nil, err
		}
	}

	return &username, nil
}
