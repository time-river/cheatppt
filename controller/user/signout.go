package user

import "cheatppt/model/redis"

func SignOut(token string) {
	rds := redis.NewRedisCient()
	rds.RevokeToken(token)
}
