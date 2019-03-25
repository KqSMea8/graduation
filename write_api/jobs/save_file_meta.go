package jobs

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

const JobName_SaveFileMeta = "save_file_meta"

type SaveFileMetaJob struct {
	file *model.File
	db   *gorm.DB
}

func NewSaveFileMetaJob(file *model.File, conn *gorm.DB) *SaveFileMetaJob {
	j := &SaveFileMetaJob{
		file: file,
		db:   conn,
	}
	return j
}

func (j *SaveFileMetaJob) GetName() string {
	return JobName_SaveFileMeta
}

func (j *SaveFileMetaJob) Run() (interface{}, error) {
	if err := mysql.FileMySQL.Save(j.db, j.file); err != nil {
		logrus.Errorf("Save FileMeta: %+v to mysql Error: %s", j.file, err)
		return nil, err
	}
	return nil, nil
}
