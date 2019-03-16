package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserInfoMySql struct {
	conn *gorm.DB
}

func NewUserInfoMySql() *UserInfoMySql {
	var err error
	h := &UserInfoMySql{}
	h.conn, err = gorm.Open(getUserInfoMySqlConfig())
	return h
}

func getUserInfoMySqlConfig() (string, string) {
	return "mysql", "g10guang:workhard@/oos_meta?charset=utf8mb4&parseTime=True&loc=Local"
}
