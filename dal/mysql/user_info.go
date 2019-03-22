package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type UserInfoMySql struct {
	conn *gorm.DB
}

func NewUserInfoMySql() *UserInfoMySql {
	var err error
	h := &UserInfoMySql{}
	h.conn, err = gorm.Open(getUserInfoMySqlConfig())
	if err != nil {
		logrus.Panicf("Create UserInfo MySQL connection Error: %s", err)
		panic(err)
	}
	return h
}

func getUserInfoMySqlConfig() (string, string) {
	return "mysql", "g10guang:workhard@/oos_meta?charset=utf8mb4&parseTime=True&loc=Local"
}
