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

func (r *FileInfoRedis) MGet(fids []int64) (metas map[int64]*model.File, missFids []int64, err error) {
	fidsKey := make([]string, len(fids))
	for i, fid := range fids {
		fidsKey[i] = r.genKey(fid)
	}
	result, err := r.conn.MGet(fidsKey...).Result()
	if err != nil {
		// redis error
		logrus.Errorf("redis mget Error: %s", err)
		// 如果 redis 发生错误，则所有 fids 都 miss
		return nil, fids, err
	}
	for i, v := range result {
		if v == nil {
			// cache not found
			missFids = append(missFids, fids[i])
			logrus.Debugf("fid: %d redis cache not found", fids[i])
			continue
		}
		m := new(model.File)
		if err = json.NewDecoder(strings.NewReader(v.(string))).Decode(m); err != nil {
			missFids = append(missFids, fids[i])
			logrus.Errorf("unmarshal FileMetas Error: %s", err)
			continue
		}
		// cache hit
		logrus.Debugf("fid: %d redis cache hit", fids[i])
		metas[fids[i]] = m
	}
	return metas, missFids, nil
}

func (r *FileInfoRedis) Del(fid int64) error {
	if _, err := r.conn.Del(r.genKey(fid)).Result(); err != nil {
		logrus.Errorf("delete redis cache Error: %s", err)
		return err
	}
	return nil
}

// 因为 MSet 不支持超时
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

func (r *FileInfoRedis) MSet(files []*model.File) error {
	pipe := r.conn.Pipeline()
	for _, meta := range files {
		b, err := json.Marshal(meta)
		if err != nil {
			logrus.Errorf("json Marshal Error: %s", err)
			continue
		}
		pipe.Set(r.genKey(meta.Fid), string(b), time.Minute * 5)
	}
	_, err := pipe.Exec()
	if err != nil {
		logrus.Errorf("Multi Set redis pipeline Error: %s", err)
		return err
	}
	return nil
}
