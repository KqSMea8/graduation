package jobs

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/jinzhu/gorm"
)

const JobName_DeleteFileMeta = "delete_file_meta"
const retryTime = 2

type DeleteFileMetaJob struct {
	fid int64
	db  *gorm.DB
}

func NewDeleteFileMetaJob(fid int64, conn *gorm.DB) *DeleteFileMetaJob {
	j := &DeleteFileMetaJob{
		fid: fid,
		db:  conn,
	}
	return j
}

func (j *DeleteFileMetaJob) GetName() string {
	return JobName_DeleteFileMeta
}

func (j *DeleteFileMetaJob) Run() (interface{}, error) {
	var err error
	for i := 0; i < retryTime; i++ {
		if err = mysql.FileMySQL.Delete(j.db, j.fid); err == nil {
			break
		}
	}
	if err == nil {
		// delete redis cache
		err = redis.FileRedis.DeleteFileMeta(j.fid)
	}
	return nil, err
}
