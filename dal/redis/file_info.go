package redis

import (
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/model"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"strings"
)

type FileInfoRedis struct {
	Conn *redis.Client
}

func NewFileInfoRedis() *FileInfoRedis {
	h := &FileInfoRedis{}
	h.Conn = redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return h
}

func (r *FileInfoRedis) GetFileMeta(fid int64) (meta model.File, err error) {
	s, err := r.Conn.Get(r.GenFileMetaRedisKey(fid)).Result()
	if err != nil {
		logrus.Errorf("redis GetFileMeta Error: %s", err)
		return nil, err
	}
	json.NewDecoder(strings.NewReader(s)).Decode(&meta)
}

func (r *FileInfoRedis) GenFileMetaRedisKey(fid int64) string {
	return fmt.Sprintf("f%d", fid)
}