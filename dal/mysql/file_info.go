package mysql

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type FileInfoMySql struct {
	conn *gorm.DB
}

func NewFileInfoMySql() *FileInfoMySql {
	var err error
	h := &FileInfoMySql{}
	h.conn, err = gorm.Open(getFileInfoMySql())
	if err != nil {
		panic(err)
	}
	return h
}

func getFileInfoMySql() (string, string) {
	return "mysql", "g10guang:workhard@/oos_meta?charset=utf8mb4&parseTime=True&loc=Local"
}
