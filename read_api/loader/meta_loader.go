package loader

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
)

const LoaderName_FileMeta = "file_meta_loader"

type FileMetaLoader struct {
	fids []int64
}

func NewFileMetaLoader(fids []int64) *FileMetaLoader {
	l := &FileMetaLoader{
		fids: fids,
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
	metas, missFids, err := redis.FileRedis.MGet(l.fids)
	if err == nil {
		// redis 出错尝试 mysql
		logrus.Debugf("redis cache hit")
	}

	if len(missFids) == 0 {
		// 全部 cache 命中
		return metas, nil
	}

	metasFromMySQL, err := mysql.FileMySQL.MGet(missFids)
	if err != nil {
		return nil, err
	}

	go l.saveRedisCache(metasFromMySQL)
	for _, m := range metasFromMySQL {
		metas[m.Fid] = m
	}

	return metas, nil
}

func (l *FileMetaLoader) saveRedisCache(metas []*model.File) {
	if len(metas) == 0 {
		return
	}
	redis.FileRedis.MSet(metas)
}
