package mysql

import (
	"errors"
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
func (h *FileInfoMySql) Save(conn *gorm.DB, file *model.File) error {
	err := h.Conn.Create(file).Error
	if err != nil {
		logrus.Error("Save Error: %s", err)
	}
	return err
}

// 删需要涉及到事务，所以由外部传递 connection
func (h *FileInfoMySql) Delete(conn *gorm.DB, fid int64) (err error) {
	if conn == nil {
		conn = h.Conn
	}
	if err = conn.Delete("fid = ?", fid).Error; err != nil {
		logrus.Errorf("Delete Error: %s", err)
	}
	return
}

func (h *FileInfoMySql) Get(fid int64) (meta model.File, err error) {
	if err = h.Conn.Where("fid = ?", fid).Find(&meta).Error; err != nil {
		logrus.Errorf("Get Error: %s", err)
	}
	return
}

func (h *FileInfoMySql) MultiGet(fids []int64) (metas []*model.File, err error) {
	if len(fids) == 0 {
		return nil, errors.New("len(fids) == 0")
	}
	if err = h.Conn.Where("fid IN (?)", fids).Find(&metas).Error; err != nil {
		logrus.Errorf("MultiGet Error: %vs", err)
	}
	return
}
