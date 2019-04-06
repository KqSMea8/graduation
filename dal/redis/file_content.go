package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"time"
)

// 图片内容的缓存
type FileContentRedis struct {
	conn *redis.Client
}

func NewFileContentRedis() *FileContentRedis {
	r := &FileContentRedis{}
	r.conn = redis.NewClient(&redis.Options{
		Addr:     "10.8.118.15:6379",
		Password: "",
		DB:       1,
	})
	return r
}

func (r *FileContentRedis) genKey(fid int64) string {
	return fmt.Sprintf("c_%d", fid)
}

func (r *FileContentRedis) genKeys(fids []int64) []string {
	keys := make([]string, len(fids))
	for i, fid := range fids {
		keys[i] = r.genKey(fid)
	}
	return keys
}

func (r *FileContentRedis) Set(fid int64, data []byte) error {
	if _, err := r.conn.Set(r.genKey(fid), data, time.Minute).Result(); err != nil {
		logrus.Errorf("redis Set fid %d content Error: %s", fid, err)
		return err
	}
	return nil
}

func (r *FileContentRedis) Get(fid int64) ([]byte, error) {
	b, err := r.conn.Get(r.genKey(fid)).Bytes()
	if err != nil {
		logrus.Errorf("redis Get fid %d content Error: %s", fid, err)
		return nil, err
	}
	return b, nil
}

func (r *FileContentRedis) Del(fids []int64) error {
	if err := r.conn.Del(r.genKeys(fids)...).Err(); err != nil {
		logrus.Errorf("redis Del fids: %v content Error: %s", fids, err)
		return err
	}
	return nil
}
