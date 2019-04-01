package loader

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
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
// 3、异步设置 redis 缓存
func (l *FileMetaLoader) Run() (interface{}, error) {
	if meta, err := redis.FileRedis.Get(l.fid); err == nil {
		logrus.Debugf("redis cache hit")
		return meta, err
	}

	if meta, err := mysql.FileMySQL.Get(l.fid); err != nil {
		return nil, err
	} else {
		l.saveRedisCache(meta)
		return meta, nil
	}
}

func (l *FileMetaLoader) saveRedisCache(meta model.File) {
	go redis.FileRedis.Set(&meta)
}
