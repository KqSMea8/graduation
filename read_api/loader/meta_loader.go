package loader

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
)

const Loader_FileMeta = "file_meta_loader"

type FileMetaLoader struct {
	fid int64
}

func (l *FileMetaLoader) GetName() string {
	return Loader_FileMeta
}

// 1、尝试从 redis 缓存中获取
// 2、如果缓存没有命中，访问 mysql
func (l *FileMetaLoader) Run() (interface{}, error) {
	redis.FileRedis
}
