package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (rdb *Redis) GetClient() *redis.Client {
	return rdb.conn
}

// verification code op:
// 1. signup
// 2. password reset
//
// (key, val):
//  - signup: (`$username-signup`, `$code $email`)
//  - password reset: (`$username-reset`, `$code`)
func (rdb *Redis) SetCode(key, val string, m int) error {
	var ctx = context.Background()
	return rdb.conn.Set(ctx, key, val, time.Duration(m)*time.Minute).Err()
}

func (rdb *Redis) GetCode(key string) (string, error) {
	var ctx = context.Background()
	return rdb.conn.Get(ctx, key).Result()
}

func (rdb *Redis) DelCode(key string) error {
	var ctx = context.Background()
	return rdb.conn.Del(ctx, key).Err()
}

// authorization code op
// (key, val): (token, username)
func (rdb *Redis) SetToken(key, val string, h int) error {
	var ctx = context.Background()
	return rdb.conn.Set(ctx, key, val, time.Duration(h)*time.Hour).Err()
}

func (rdb *Redis) ExistsToken(key string) (int64, error) {
	var ctx = context.Background()
	return rdb.conn.Exists(ctx, key).Result()
}

func (rdb *Redis) RevokeToken(key string) error {
	var ctx = context.Background()
	// only revoke token except session
	return rdb.conn.Del(ctx, key).Err()
}

// session op
// (key, val): (session, token)
func (rdb *Redis) SetSession(key, val string, h int) error {
	var ctx = context.Background()
	return rdb.conn.Set(ctx, key, val, time.Duration(h)*time.Hour).Err()
}

// return value: `(*string, error)``
// `(xxx, nil)`: key exists and doesn't expire
// `(xxx, err)`: key exists and doesn't expire, but TTL error
// `(nil, err)`: key don't exist (expire or not set)
func (rdb *Redis) ValidThenRenewSession(key string, h int) (*string, error) {
	var ctx = context.Background()

	// err is `nil` is key don't exists or expire
	val, err := rdb.conn.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	} else {
		expire := time.Duration(h) * time.Hour
		err := rdb.conn.Expire(ctx, key, expire).Err()
		return &val, err
	}
}
