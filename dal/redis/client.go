package redis

var FileRedis *FileInfoRedis

func init() {
	FileRedis = NewFileInfoRedis()
}
