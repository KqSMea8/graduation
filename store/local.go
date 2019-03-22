package store

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// concurrency safe
type LocalStorage struct {
	dirPath string
}

func NewLocalStorage(dir string) *LocalStorage {
	h := &LocalStorage{
		dirPath: dir,
	}
	return h
}

func (h *LocalStorage) Write(fid int64, data []byte) error {
	logrus.Debugf("Write fid: %d", fid)
	filePath := path.Join(h.dirPath, strconv.Itoa(int(fid)))
	fp, err := os.Create(filePath)
	if err != nil {
		logrus.Errorf("Create FileBytes %s")
		return err
	}
	n, err := fp.Write(data)
	logrus.Infof("write %d bytes into %s", n, filePath)
	if err != nil {
		logrus.Errorf("write to %s Error: %s", filePath, err)
	}
	return err
}

func (h *LocalStorage) Read(fid int64) (data []byte, err error) {
	logrus.Debugf("Read fid: %d", fid)
	filePath := path.Join(h.dirPath, strconv.Itoa(int(fid)))
	data, err = ioutil.ReadFile(filePath)
	if err != nil {
		logrus.Errorf("read %s Error: %s", filePath, err)
	}
	return data, err
}

func (h *LocalStorage) Delete(fid int64) error {
	logrus.Debugf("Delete fid: %d", fid)
	filePath := path.Join(h.dirPath, strconv.Itoa(int(fid)))
	if err := os.Remove(filePath); err != nil {
		logrus.Errorf("delete %s Error: %s", filePath, err)
		return err
	}
	return nil
}
