package store

import (
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type HdfsStorage struct {
	username     string
	hdfsPath     string
	dataNodeIp   string
	dataNodePort int
	webHdfsIp    string
	webHdfsPort  int
	client       *http.Client
}

func NewHdfsStorage(hdfsUsername, hdfsPath, dataNodeIp, webHdfsIp string, dataNodePort, webHdfsPort int) *HdfsStorage {
	s := &HdfsStorage{
		username:     hdfsUsername,
		hdfsPath:     hdfsPath,
		dataNodeIp:   dataNodeIp,
		dataNodePort: dataNodePort,
		webHdfsIp:    webHdfsIp,
		webHdfsPort:  webHdfsPort,
		client:       &http.Client{},
	}
	s.hdfsPath = strings.Trim(s.hdfsPath, "/")
	s.init()
	// 判断 hdfsPath 是否存在，如果不存在，则创建
	return s
}

func (s *HdfsStorage) init() {
	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s?op=MKDIRS&user.naem=%s", s.mkUrlPrefix(), s.username), nil)
	if err != nil {
		panic(err)
	}
	rsp, err := s.client.Do(req)
	logrus.Debugf("webhdfs mkdirs req: %+v resp: %+v", req, rsp)
	if err != nil {
		panic(err)
	}
	if rsp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("http status code: %d not OK", rsp.StatusCode))
	}
}

func (s *HdfsStorage) mkUrlPrefix() string {
	return fmt.Sprintf("http://%s:%d/webhdfs/v1/%s", s.webHdfsIp, s.webHdfsPort, s.hdfsPath)
}

//type Storage interface {
//	Write(fid int64, data []byte) error
//	WriteWithFormat(fid int64, format constdef.ImageFormat, data []byte) error
//	Read(fid int64) (data []byte, err error)
//	ReadWithFormat(fid int64, format constdef.ImageFormat) (data []byte, err error)
//	Delete(fid int64) error
//}

// 这里需要先进行一次 redirect，然后再将数据写入 data node
func (s *HdfsStorage) write(filePath string, reader io.Reader) error {
	url := fmt.Sprintf("%s/%s?user.name=%s&op=CREATE&overwrite=false&replicatoin=1", s.mkUrlPrefix(), filePath, s.username)
	r, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		panic(err)
	}
	rsp, err := s.client.Do(r)
	logrus.Debugf("webhdfs create 1 req: %+v resp: %+v", r, rsp)
	if err != nil {
		logrus.Errorf("webhdfs create 1 Error: %s", err)
		return err
	}
	newUrl := rsp.Header.Get("Location")
	r2, err := http.NewRequest(http.MethodPut, newUrl, reader)
	if err != nil {
		panic(err)
	}
	newRsp, err := s.client.Do(r2)
	if err != nil {
		logrus.Errorf("webhdfs create 2 file Error: %s", err)
		return err
	}
	logrus.Debugf("webhdfs create 2 req: %+v resp: %+v", r, newRsp)
	if newRsp.StatusCode != http.StatusOK {
		logrus.Errorf("webhdfs create 2 statusCode: %d", http.StatusOK)
		return fmt.Errorf("webhdfs create 2 status code: %d", newRsp.StatusCode)
	}
	return nil
}

func (s *HdfsStorage) read(filePath string) (io.Reader, error) {
	url := fmt.Sprintf("%s/%s?user.name=%s&op=OPEN", s.mkUrlPrefix(), filePath, s.username)
	r, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		panic(err)
	}
	rsp, err := s.client.Do(r)
	logrus.Debugf("webhdfs open 1 req: %+v resp: %+v", r, rsp)
	if err != nil {
		logrus.Errorf("webhdfs open 1 Error: %s", err)
		return nil, err
	}
	newUrl := rsp.Header.Get("Location")
	r2, err := http.NewRequest(http.MethodGet, newUrl, nil)
	if err != nil {
		panic(err)
	}
	newRsp, err := s.client.Do(r2)
	logrus.Debugf("webhdfs open 2 req: %+v resp: %+v", r, newRsp)
	if err != nil {
		logrus.Errorf("webhdfs open 2 file Error: %s", err)
		return nil, err
	}
	if newRsp.StatusCode != http.StatusOK {
		logrus.Errorf("webhdfs open 2 status code: %d", http.StatusOK)
		return nil, fmt.Errorf("webhdfs open 2 status code: %d", newRsp.StatusCode)
	}
	return newRsp.Body, nil
}

func (s *HdfsStorage) delete_(filePath string) error {
	url := fmt.Sprintf("%s/%s?user.name=%s&op=DEELTE", s.mkUrlPrefix(), filePath, s.username)
	r, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		panic(err)
	}
	rsp, err := s.client.Do(r)
	logrus.Debugf("webhdfs delete 1 req: %+v resp: %+v", r, rsp)
	if err != nil {
		logrus.Errorf("webhdfs delete 1 Error: %s", err)
		return err
	}
	if rsp.StatusCode != http.StatusOK {
		logrus.Errorf("webhdfs delete status code: %d", rsp.StatusCode)
		return fmt.Errorf("webhdfs delete status code: %d", rsp.StatusCode)
	}
	return nil
}

func (s *HdfsStorage) genFilePath(fid int64, format ...constdef.ImageFormat) string {
	return path.Join("/oss/picture/", s.genFileName(fid, format...))
}

func (s *HdfsStorage) genFileName(fid int64, format ...constdef.ImageFormat) string {
	if len(format) == 0 {
		return strconv.FormatInt(fid, 10)
	}
	return fmt.Sprintf("%d_%d", fid, format[0])
}
