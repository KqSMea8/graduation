package store

import (
	"errors"
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
)

// concurrency safe
type LocalStorage struct {
	dirPath string
}

func NewLocalStorage() *LocalStorage {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic(errors.New("GOPATH not exists in env"))
	}
	dir := path.Join(goPath, "src/github.com/g10guang/graduation/oss")
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0777); err != nil {
			panic(err)
		}
	}
	h := &LocalStorage{
		dirPath: dir,
	}
	return h
}

func (h *LocalStorage) Write(fid int64, data []byte) error {
	filePath := h.genFilePath(fid)
	return h.write(filePath, data)
}

func (h *LocalStorage) Read(fid int64) (data []byte, err error) {
	filePath := h.genFilePath(fid)
	return h.read(filePath)
}

// 删除需要将其他格式一并删除
func (h *LocalStorage) Delete(fid int64) error {
	go os.Remove(h.genFilePath(fid))
	for _, f := range constdef.ImageFormatList {
		go os.Remove(h.genFilePath(fid, f))
	}
	return nil
}

func (h *LocalStorage) WriteWithFormat(fid int64, format constdef.ImageFormat, data []byte) error {
	filepath := h.genFilePath(fid, format)
	return h.write(filepath, data)
}

func (h *LocalStorage) ReadWithFormat(fid int64, format constdef.ImageFormat) (data []byte, err error) {
	filepath := h.genFilePath(fid, format)
	return h.read(filepath)
}

func (h *LocalStorage) write(path string, b []byte) error {
	err := ioutil.WriteFile(path, b, 0666)
	if err != nil {
		logrus.Errorf("write %s Error: %s", path, err)
	}
	return err
}

func (h *LocalStorage) read(path string) ([]byte, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		logrus.Errorf("read %s Error: %s", path, err)
	}
	return b, err
}

func (h *LocalStorage) genFileName(fid int64, format ...constdef.ImageFormat) string {
	if len(format) == 0 {
		return fmt.Sprintf("%d", fid)
	} else {
		return fmt.Sprintf("%d_%d", fid, format[0])
	}
}

func (h *LocalStorage) genFilePath(fid int64, format ...constdef.ImageFormat) string {
	return path.Join(h.dirPath, h.genFileName(fid, format...))
}
