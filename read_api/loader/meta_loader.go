package loader

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
)

const LoaderName_FileMeta = "file_meta_loader"

type FileMetaLoader struct {
	fid int64
}

func NewFileMetaLoader(fid int64) *FileMetaLoader {
	l := &FileMetaLoader{
		fid: fid,
	}
	return l
}

func (l *FileMetaLoader) GetName() string {
	return LoaderName_FileMeta
}

// 1、尝试从 redis 缓存中获取
// 2、如果缓存没有命中，访问 mysql
func (l *FileMetaLoader) Run() (interface{}, error) {
	if meta, err := redis.FileRedis.Get(l.fid); err == nil {
		return meta, err
	}

	if meta, err := mysql.FileMySQL.Get(l.fid); err != nil {
		return nil, err
	} else {
		return meta, nil
	}
}
