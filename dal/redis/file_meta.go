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
	r := &FileInfoRedis{}
	r.conn = redis.NewClient(&redis.Options{
		Addr:     "10.8.118.15:6379",
		Password: "",
		DB:       0,
	})
	return r
}


func (r *FileInfoRedis) genKey(fid int64) string {
	return fmt.Sprintf("f%d", fid)
}

func (r *FileInfoRedis) genKeys(fids []int64) []string {
	keys := make([]string, len(fids))
	for i, fid := range fids {
		keys[i] = r.genKey(fid)
	}
	return keys
}

// 只要获取 redis 失败就认为 cache not found
// 不管是网络超时还是确实 key 不存在
// TODO 区分清楚是 redis key 不存在的情况，免得压垮下游
func (r *FileInfoRedis) Get(fid int64) (meta model.File, err error) {
	s, err := r.conn.Get(r.genKey(fid)).Result()
	if err != nil {
		logrus.Errorf("redis Get Fid: %d not found", fid)
		return meta, err
	}
	if err = json.NewDecoder(strings.NewReader(s)).Decode(&meta); err != nil {
		logrus.Errorf("unmarshal redis cache Error: %s", err)
	}
	return
}

func (r *FileInfoRedis) MGet(fids []int64) (metas map[int64]*model.File, missFids []int64, err error) {
	fidsKey := r.genKeys(fids)
	result, err := r.conn.MGet(fidsKey...).Result()
	if err != nil {
		// redis error
		logrus.Errorf("redis mget Error: %s", err)
		// 如果 redis 发生错误，则所有 fids 都 miss
		return nil, fids, err
	}
	metas = make(map[int64]*model.File, len(fids))
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

func (r *FileInfoRedis) Del(fids []int64) error {
	fidsKey := r.genKeys(fids)
	if _, err := r.conn.Del(fidsKey...).Result(); err != nil {
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
