package redis

import (
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
	"testing"
	"time"
)

func TestFileSet(t *testing.T) {
	logrus.SetReportCaller(true)
	err := FileRedis.Set(&model.File{
		Fid:        1,
		Uid:        1,
		Name:       "test",
		Size:       10,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		Md5:        "5d41402abc4b2a76b9719d911017c592",
		Extra:      "hello",
		Status:     0,
	})
	if err != nil {
		logrus.Errorf("Redis Set Error: %s", err)
	} else {
		logrus.Infof("Redis Set Success")
	}
}

func TestFileGet(t *testing.T) {
	file, err := FileRedis.Get(1)
	if err != nil {
		logrus.Errorf("Redis Get Error: %s", err)
	} else {
		logrus.Infof("file: %+v Error: %s", file, err)
	}
}

func TestUserSet(t *testing.T) {
	err := UserRedis.Set(&model.User{
		Uid:    1,
		Status: 1,
		Extra:  "hello world",
	})
	if err != nil {
		logrus.Errorf("Redis Set Error: %s", err)
	} else {
		logrus.Infof("Redis Set Success")
	}
}

func TestUserGet(t *testing.T) {
	user, err := UserRedis.Get(1)
	if err != nil {
		logrus.Errorf("Redis Get Error: %s", err)
	} else {
		logrus.Infof("Redis Get User: %+v", user)
	}
}
