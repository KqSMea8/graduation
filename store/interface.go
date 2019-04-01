package store

import (
	"github.com/g10guang/graduation/constdef"
	"io"
)

// define the interface for storage layer
type Storage interface {
	Write(fid int64, reader io.Reader) error
	WriteWithFormat(fid int64, format constdef.ImageFormat, reader io.Reader) error
	Read(fid int64) (reader io.Reader, err error)
	ReadWithFormat(fid int64, format constdef.ImageFormat) (reader io.Reader, err error)
	Delete(fid int64) error
}
