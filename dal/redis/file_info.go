package redis

import (
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/model"
	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type FileInfoRedis struct {
	conn *redis.Client
}

func NewFileInfoRedis() *FileInfoRedis {
	h := &FileInfoRedis{}
	h.conn = redis.NewClient(&redis.Options{
		Addr:     "10.8.118.15:6379",
		Password: "",
		DB:       0,
	})
	return h
}


func (r *FileInfoRedis) genKey(fid int64) string {
	return fmt.Sprintf("f%d", fid)
}

func (r *FileInfoRedis) Get(fid int64) (meta model.File, err error) {
	s, err := r.conn.Get(r.genKey(fid)).Result()
	if err == redis.Nil {
		logrus.Debugf("redis Get Fid: %d not found", fid)
		return
	}
	if err != nil {
		logrus.Errorf("redis Get Error: %s", err)
		return meta, err
	}
	if err = json.NewDecoder(strings.NewReader(s)).Decode(&meta); err != nil {
		logrus.Errorf("unmarshal redis cache Error: %s", err)
	}
	return
}

func (r *FileInfoRedis) Del(fid int64) error {
	if _, err := r.conn.Del(r.genKey(fid)).Result(); err != nil {
		logrus.Errorf("delete redis cache Error: %s", err)
		return err
	}
	return nil
}

func (r *FileInfoRedis) Set(file *model.File) error {
	b, err := json.Marshal(file)
	logrus.Debugf("set redis cache fid: %d json: %s", file.Fid, string(b))
	if err != nil {
		logrus.Errorf("json Marshal Error: %s", err)
		return err
	}
	if _, err = r.conn.Set(r.genKey(file.Fid), string(b), time.Minute*5).Result(); err != nil {
		logrus.Errorf("redis Set file: %+v", file)
	}
	return err
}
