package mysql

import (
	"os"

	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
)

type FileInfoMySql struct {
	Conn *gorm.DB
}

func NewFileInfoMySql() *FileInfoMySql {
	var err error
	h := &FileInfoMySql{}
	h.Conn, err = gorm.Open(getFileInfoMySql())
	if err != nil {
		logrus.Panicf("Create UserInfo MySQL connection Error: %s", err)
		panic(err)
	}
	// In test env, print SQL gorm execute
	if _, exists := os.LookupEnv(constdef.ENV_TestEnv); exists {
		h.Conn = h.Conn.Debug()
	}
	return h
}

func getFileInfoMySql() (string, string) {
	return "mysql", "g10guang:workhard@/oos_meta?charset=utf8mb4&parseTime=True&loc=Local"
}

// 写需要涉及到事务，所以由外部传递 connection
func (h *FileInfoMySql) SaveFileMeta(conn *gorm.DB, file *model.File) error {
	err := h.Conn.Create(file).Error
	if err != nil {
		logrus.Error("SaveFileMeta Error: %s", err)
	}
	return err
}

// 删需要涉及到事务，所以由外部传递 connection
func (h *FileInfoMySql) DeleteFileMeta(conn *gorm.DB, fid int64) (err error) {
	if err = conn.Delete("fid = ?", fid).Error; err != nil {
		logrus.Errorf("DeleteFileMeta Error: %s", err)
	}
	return
}
