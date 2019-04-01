package constdef

// request 请求中的字段常量定义
const (
	Param_Uid      = "uid"
	Param_Fid      = "fid"
	Param_File     = "file"
	Param_Filename = "filename"
)

// OS environment variable key
const (
	ENV_TestEnv    = "test_env"
	ENV_ProductEnv = "product_env"
)

// NSQ 消息队列相关常量
const (
	// Picture Post Event Topic
	PostFileEventTopic = "post_file"
	// Delete Picture Event Topic
	DeleteFileEventTopic = "delete_file"
)

type ImageFormat int16

// 支持的图片格式
const (
	InvalidImageFormat ImageFormat = iota
	Jpeg
	Png
)

var ImageFormatList = []ImageFormat{Jpeg, Png}

// NSQ 配置信息
const (
	NsqLookupdAddr = "10.8.118.15:4161"
	NsqdAddr       = "10.8.118.15:4150"
)

// HDFS 配置
const (
	WebHdfsAddr = "10.8.118.15:50070"
	WebHdfsUser = "root"
	WebHdfsDir  = "/oss/image"
)
