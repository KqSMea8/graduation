package redis

var FileRedis *FileInfoRedis
var UserRedis *UserInfoRedis

func init() {
	FileRedis = NewFileInfoRedis()
	UserRedis = NewUserInfoRedis()
}
