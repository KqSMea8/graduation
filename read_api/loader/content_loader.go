package loader

import (
	"github.com/g10guang/graduation/store"
	"sync"
)

const Loader_FileContent = "file_content_loader"

type FileContentLoader struct {
	fids    []int64
	storage store.Storage
}

func NewFileContentLoader(fids []int64, storage store.Storage) *FileContentLoader {
	return &FileContentLoader{
		fids:    fids,
		storage: storage,
	}
}

func (l *FileContentLoader) GetName() string {
	return Loader_FileContent
}

func (l *FileContentLoader) Run() (interface{}, error) {
	result := make(map[int64][]byte, len(l.fids))
	var wg sync.WaitGroup
	wg.Add(len(l.fids))
	var mutex sync.Mutex
	for _, id := range l.fids {
		go func() {
			b, err := l.storage.Read(id)
			if err != nil {
				return
			}
			mutex.Lock()
			result[id] = b
			defer mutex.Unlock()
		}()
	}
	wg.Wait()

	return result, nil
}
