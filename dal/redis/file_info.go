package redis

import "github.com/go-redis/redis"

type FileInfoRedis struct {
	conn *redis.Client
}

func NewFileInfoRedis() *FileInfoRedis {
	h := &FileInfoRedis{}
	h.conn = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return h
}
