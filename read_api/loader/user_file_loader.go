package loader

import (
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/dal/redis"
	"github.com/g10guang/graduation/model"
)

const LoaderName_UserFile = "user_file_loader"

type UserFileLoader struct {
	uid    int64
	offset int64
	limit  int64
}

func NewUserFileLoader(uid, offset, limit int64) *UserFileLoader {
	l := &UserFileLoader{
		uid:    uid,
		offset: offset,
		limit:  limit,
	}
	return l
}

func (l *UserFileLoader) GetName() string {
	return LoaderName_UserFile
}

func (l *UserFileLoader) Run() (interface{}, error) {
	metas, err := redis.FileRedis.GetPageCache(l.uid, l.offset, l.limit)
	if err == nil && metas != nil {
		return metas, nil
	}

	fileMetas, err := mysql.FileMySQL.GetFileByUid(l.uid, l.offset, l.limit)
	if err != nil {
		return nil, err
	}

	go l.saveRedisCache(metas)
	return fileMetas, nil
}

func (l *UserFileLoader) saveRedisCache(metas []*model.File) error {
	return redis.FileRedis.SetPageCache(l.uid, l.offset, l.limit, metas)
}
