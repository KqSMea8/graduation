package main

import (
	"encoding/json"
	"github.com/g10guang/graduation/consumer/handler"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// 将图片转化为 jpeg 格式
func jpegCompress(message *nsq.Message) error {
	msg := parsePostFileEventMsg(message.Body)
	if msg == nil {
		return errors.New("message error")
	}
	h := handler.NewJpegCompressHandler(msg)
	if err := h.Handle(tools.NewCtxWithLogID()); err != nil {
		logrus.Errorf("JpegCompressHandler Error: %s", err)
		return err
	}
	logrus.Infof("JpegCompressHandler Success")
	return nil
}

// 将图片转化为 png 格式
func pngCompress(message *nsq.Message) error {
	return nil
}

func parsePostFileEventMsg(body []byte) *model.PostFileEvent {
	m := new(model.PostFileEvent)
	if err := json.Unmarshal(body, m); err != nil {
		logrus.Errorf("PostFileEvent message Error: %s", err)
	}
	return m
}
