package store

import (
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// concurrency safe
type LocalStorage struct {
	dirPath string }

func NewLocalStorage(dir string) *LocalStorage {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err = os.Mkdir(dir, 0666); err != nil {
			panic(err)
		}
	}
	h := &LocalStorage{
		dirPath: dir,
	}
	return h
}

func (h *LocalStorage) Write(fid int64, data []byte) error {
	filePath := path.Join(h.dirPath, strconv.Itoa(int(fid)))
	return h.write(filePath, data)
}

func (h *LocalStorage) Read(fid int64) (data []byte, err error) {
	filePath := path.Join(h.dirPath, strconv.Itoa(int(fid)))
	return h.read(filePath)
}

func (h *LocalStorage) Delete(fid int64) error {
	filePath := path.Join(h.dirPath, strconv.Itoa(int(fid)))
	if err := os.Remove(filePath); err != nil {
		logrus.Errorf("delete %s Error: %s", filePath, err)
		return err
	}
	return nil
}

func (h *LocalStorage) WriteWithFormat(fid int64, format constdef.ImageFormat, data []byte) error {
	filepath := path.Join(h.dirPath, fmt.Sprintf("%d_%d", fid, format))
	return h.write(filepath, data)
}

func (h *LocalStorage) ReadWithFormat(fid int64, format constdef.ImageFormat) (data []byte, err error) {
	filepath := path.Join(h.dirPath, fmt.Sprintf("%d_%d", fid, format))
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
