package store

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"testing"
	"time"
)

func TestHDFSCreate(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	storage := NewHdfsStorage("10.8.118.15:50070", "root", "/oss/image")
	fp, err := os.Open("/Users/g10guang/Public/output.jpeg")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	if err = storage.Write(2, fp); err != nil {
		panic(err)
	}
}

func TestHDFSOpen(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	storage := NewHdfsStorage("10.8.118.15:50070", "root", "/oss/image")
	reader, err := storage.Read(2)
	if err != nil {
		panic(err)
	}
	fp, err := os.OpenFile("./tset.jpeg", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	n, err := io.Copy(fp, reader)
	if err != nil {
		panic(err)
	}
	logrus.Infof("%d bytes written", n)
}

func TestHDFSDelete(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	storage := NewHdfsStorage("10.8.118.15:50070", "root", "/oss/image")
	storage.Delete(1)
	time.Sleep(time.Second)
}
