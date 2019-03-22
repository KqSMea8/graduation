package store

// define the interface for storage layer
type Storage interface {
	Write(fid int64, data []byte) error
	Read(fid int64) (data []byte, err error)
	Delete(fid int64) error
}
