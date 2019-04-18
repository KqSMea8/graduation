package store

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/vladimirvivien/gowfs"
)

func TestHDFSCreate(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	storage := NewHdfsStorage("10.224.12.131:50070", "root", "/oss/image")
	fp, err := os.Open("/Users/g10guang/Public/output.jpeg")
	if err != nil {
		panic(err)
	}
	defer fp.Close()
	if err = storage.Write(1112599362007408640, fp); err != nil {
		panic(err)
	}
}

func TestHDFSOpen(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	storage := NewHdfsStorage("10.224.12.131:50070", "root", "/oss/image")
	reader, err := storage.Read(1112599362007408640)
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
	storage := NewHdfsStorage("10.224.12.131:50070", "root", "/oss/image")
	storage.Delete(1)
	time.Sleep(time.Second)
}

func TestHDFSSum(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	storage := NewHdfsStorage("10.224.12.131:50070", "root", "/oss/image")
	sum, err := storage.client.GetFileChecksum(gowfs.Path{Name: "/oss/image/1112619134648320000"})
	if err != nil {
		panic(err)
	}
	logrus.Infof("sum: %+v", sum)
}
