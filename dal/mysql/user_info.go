package mysql

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"os"
)

type UserInfoMySql struct {
	Conn *gorm.DB
}

func NewUserInfoMySql() *UserInfoMySql {
	var err error
	h := &UserInfoMySql{}
	h.Conn, err = gorm.Open(getUserInfoMySqlConfig())
	if err != nil {
		logrus.Panicf("Create UserInfo MySQL connection Error: %s", err)
		panic(err)
	}
	// In test env, print SQL gorm execute
	if _, exists := os.LookupEnv(constdef.ENV_ProductEnv); !exists {
		h.Conn = h.Conn.Debug()
	}
	return h
}

func getUserInfoMySqlConfig() (string, string) {
	return "mysql", "g10guang:hello@tcp(10.8.118.15:3306)/oss_meta?charset=utf8mb4&parseTime=True&loc=Local"
}

func (h *UserInfoMySql) Save(conn *gorm.DB, user *model.User) (err error) {
	if conn == nil {
		conn = h.Conn
	}

	if err = conn.Save(user).Error; err != nil {
		logrus.Errorf("Save user %+v info Error: %s", user, err)
	}
	return
}

func (h *UserInfoMySql) Get(uid int64) (user *model.User, err error) {
	user = new(model.User)
	if err = h.Conn.Where("uid IN (?)", uid).Find(user).Error; err != nil {
		logrus.Errorf("Get uid: %d Error: %s", uid, err)
	}
	return
}
