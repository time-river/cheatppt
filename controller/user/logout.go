package user

import "cheatppt/model/redis"

func Logout(token string) {
	rds := redis.NewRedisCient()
	rds.TokenRevoke(token)
}
