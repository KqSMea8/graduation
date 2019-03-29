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

type UserInfoRedis struct {
	conn *redis.Client
}

func NewUserInfoRedis() *UserInfoRedis {
	h := &UserInfoRedis{}
	h.conn = redis.NewClient(&redis.Options{
		Addr:     "10.8.118.15:6379",
		Password: "",
		DB:       0,
	})
	return h
}

func (r *UserInfoRedis) genKey(uid int64) string {
	return fmt.Sprintf("u_%d", uid)
}

func (r *UserInfoRedis) Set(user *model.User) error {
	b, _ := json.Marshal(user)
	_, err := r.conn.Set(r.genKey(user.Uid), string(b), time.Hour).Result()
	if err != nil {
		logrus.Errorf("redis set User %+v Error: %s", user, err)
	}
	return err
}

func (r *UserInfoRedis) Get(uid int64) (user *model.User, err error) {
	str, err := r.conn.Get(r.genKey(uid)).Result()
	if err != nil {
		logrus.Errorf("redis get uid: %d Error: %s", uid, err)
		return nil, err
	}
	user = new(model.User)
	if err = json.NewDecoder(strings.NewReader(str)).Decode(user); err != nil {
		logrus.Errorf("json unmarshal user model Error: %s", err)
		return nil, err
	}
	return user, nil
}
