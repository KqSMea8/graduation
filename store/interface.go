package store

import "github.com/g10guang/graduation/constdef"

// define the interface for storage layer
type Storage interface {
	Write(fid int64, data []byte) error
	WriteWithFormat(fid int64, format constdef.ImageFormat, data []byte) error
	Read(fid int64) (data []byte, err error)
	ReadWithFormat(fid int64, format constdef.ImageFormat) (data []byte, err error)
	Delete(fid int64) error
}
