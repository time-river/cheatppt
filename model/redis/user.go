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

	num, err := rdb.conn.ZCard(ctx, key).Result()
	if err != nil {
		return err
	}

	if num >= 3 {
		pop := num - 3
		rdb.conn.ZPopMin(ctx, key, pop)
	}

	/* ignore, all of adding new or refreshing score will be regarded as success */
	if err := rdb.conn.ZAdd(ctx, key, val).Err(); err != nil {
		return err
	}

	return nil
}

func (rdb *Redis) sessionPop(key string, member string) {
	zsetLock.Lock()
	defer zsetLock.Unlock()

	rdb.conn.ZRem(ctx, key, member)
}

func (rdb *Redis) sessionRefresh(key string, member string) error {
	zsetLock.Lock()
	defer zsetLock.Unlock()

	now := time.Now().Unix()
	val := redis.Z{
		Score:  float64(now),
		Member: member,
	}

	num, err := rdb.conn.ZAdd(ctx, key, val).Result()
	if err != nil {
		return err
	} else if num != 0 {
		rdb.conn.ZRem(ctx, key, member)
		return errors.New("Expect refresh member score but add member")
	}

	return nil
}

func (rdb *Redis) TokenLease(token string, username string) error {
	if err := rdb.sessionCreate(username, token); err != nil {
		return err
	}

	return rdb.conn.Set(ctx, token, username, rdb.lease).Err()
}

func (rdb *Redis) TokenRevoke(token string) {
	/*
	 * 1. Get(token): none stands for the token has been expired
	 * 2. Del(token)
	 * 3. pop ZSet(token)
	 */
	username, err := rdb.conn.Get(ctx, token).Result()
	if err != nil {
		return
	}

	err = rdb.conn.Del(ctx, token).Err()
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
	username, err := rdb.conn.Get(ctx, token).Result()
	if err != nil {
		return nil, err
	}

	if err := rdb.sessionRefresh(username, token); err != nil {
		return nil, err
	}

	done, err := rdb.conn.Expire(ctx, token, rdb.lease).Result()
	if err != nil {
		/* EXPIRE failed, but still effective
		 * TODO: warning here
		 */
	} else if !done {
		/* no one, add it */
		if err := rdb.conn.Set(ctx, token, username, rdb.lease).Err(); err != nil {
			return nil, err
		}
	}

	return &username, nil
}

func (rdb *Redis) TokenValue(token string) (string, error) {
	return rdb.conn.Get(ctx, token).Result()
}

// Verification Code op
func (rdb *Redis) SetCode(key, val string, m int) error {
	return rdb.conn.Set(ctx, key, val, time.Duration(m)*time.Minute).Err()
}

func (rdb *Redis) GetCode(key string) (string, error) {
	return rdb.conn.Get(ctx, key).Result()
}

func (rdb *Redis) DelCode(key string) error {
	return rdb.conn.Del(ctx, key).Err()
}
