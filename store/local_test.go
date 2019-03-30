package store

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestLocalStorage(t *testing.T) {
	s := NewLocalStorage()
	fakeFid := int64(1)
	err := s.Write(fakeFid, []byte("hello world"))
	if err != nil {
		logrus.Errorf("Cannot write file Error: %s", err)
	} else {
		logrus.Infof("Write file success")
	}

	b, err := s.Read(fakeFid)
	if err != nil {
		logrus.Errorf("Cannot read file Error: %s", err)
	} else {
		logrus.Infof("Read file success %s", string(b))
	}

	err = s.Delete(fakeFid)
	if err != nil {
		logrus.Errorf("Delete file Error: %s", err)
	} else {
		logrus.Infof("Delete file success")
	}
}
