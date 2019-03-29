package main

import (
	"encoding/json"
	"github.com/g10guang/graduation/consumer/post_event/handler"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// 将图片转化为 jpeg/png 格式
func compress(message *nsq.Message) error {
	msg := parsePostFileEventMsg(message.Body)
	if msg == nil {
		return errors.New("message error")
	}
	h := handler.NewCompressHandler(msg)
	if err := h.Handle(tools.NewCtxWithLogID()); err != nil {
		logrus.Errorf("CompressHandler Error: %s", err)
		return err
	}
	logrus.Infof("CompressHandler Success")
	return nil
}

func parsePostFileEventMsg(body []byte) *model.PostFileEvent {
	m := new(model.PostFileEvent)
	if err := json.Unmarshal(body, m); err != nil {
		logrus.Errorf("PostFileEvent message Error: %s", err)
	}
	return m
}
