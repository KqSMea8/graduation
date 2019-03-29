package main

import (
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/consumer/delete_event/handler"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

func delete_(message *nsq.Message) error {
	msg := parseDeleteFileEventMsg(message.Body)
	if msg.Fid == 0 || msg.Uid == 0 {
		return fmt.Errorf("invalid message: %s", string(message.Body))
	}
	h := handler.NewDeleteStorageHandler(msg)
	if err := h.Handle(tools.NewCtxWithLogID()); err != nil {
		logrus.Errorf("Delete Storage Error: %s", err)
		return err
	}
	return nil
}

func parseDeleteFileEventMsg(body []byte) *model.DeleteFileEvent {
	msg := &model.DeleteFileEvent{}
	err := json.Unmarshal(body, msg)
	if err != nil {
		logrus.Errorf("unmarshal DeleteFileEvent message Error: %s", err)
	}
	return msg
}
