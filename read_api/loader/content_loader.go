package loader

import (
	"github.com/g10guang/graduation/store"
)

const LoaderName_FileContent = "file_content_loader"

type FileContentLoader struct {
	fid     int64
	storage store.Storage
}

func NewFileContentLoader(fid int64, storage store.Storage) *FileContentLoader {
	return &FileContentLoader{
		fid:     fid,
		storage: storage,
	}
}

func (l *FileContentLoader) GetName() string {
	return LoaderName_FileContent
}

func (l *FileContentLoader) Run() (interface{}, error) {
	return l.storage.Read(l.fid)
}
