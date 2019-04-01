package handler

import (
	"bufio"
	"bytes"
	"context"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"strings"
	"sync"
)

type CompressHandler struct {
	msg *model.PostFileEvent
}

func NewCompressHandler(msg *model.PostFileEvent) *CompressHandler {
	h := &CompressHandler{
		msg: msg,
	}
	return h
}

func  (h *CompressHandler) Handle(ctx context.Context) error {
	reader, err := storage.Read(h.msg.Fid)
	if err != nil {
		return err
	}
	meta, err := mysql.FileMySQL.Get(h.msg.Fid)
	if err != nil {
		return err
	}

	jpegBuffer, pngBuffer := &bytes.Buffer{}, &bytes.Buffer{}
	jpegWriter := bufio.NewWriter(jpegBuffer)
	pngWriter := bufio.NewWriter(pngBuffer)
	var imageFormat constdef.ImageFormat

	switch {
	case strings.HasSuffix(meta.Name, "jpeg") || strings.HasSuffix(meta.Name, "jpg"):
		imageFormat = constdef.Jpeg
	case strings.HasSuffix(meta.Name, "png"):
		imageFormat = constdef.Png
	default:
		imageFormat = constdef.InvalidImageFormat
	}
	if imageFormat == constdef.InvalidImageFormat {
		logrus.Errorf("invalid image format. image name: %s", meta.Name)
		return nil
	}
	if err = tools.ImageCompress(reader, jpegWriter, pngWriter, imageFormat); err != nil {
		logrus.Errorf("ImageCompress Error: %s", err)
		return err
	}

	var wg sync.WaitGroup
	wg.Add(2)
	// FIXME 如何处理失败的场景
	go func() {
		defer wg.Done()
		if err := storage.WriteWithFormat(h.msg.Fid, constdef.Jpeg, jpegBuffer); err != nil {
			logrus.Errorf("Write Fid: %d jpeg format Error: %s", err)
		}
	}()
	go func() {
		defer wg.Done()
		if err := storage.WriteWithFormat(h.msg.Fid, constdef.Png, pngBuffer); err != nil {
			logrus.Errorf("Write Fid: %d png format Error: %s", err)
		}
	}()
	wg.Wait()
	return nil
}
