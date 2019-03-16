package redis

import "github.com/go-redis/redis"

type UserInfoRedis struct {
	conn *redis.Client
}

func NewUserInfoRedis() *UserInfoRedis {
	h := &UserInfoRedis{}
	h.conn = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return h
}
