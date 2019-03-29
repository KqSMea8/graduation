package handler

import (
	"context"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
)

type DeleteStorageHandler struct {
	msg *model.DeleteFileEvent
}

func NewDeleteStorageHandler(msg *model.DeleteFileEvent) *DeleteStorageHandler {
	h := &DeleteStorageHandler{
		msg: msg,
	}
	return h
}

func (h *DeleteStorageHandler) Handle(ctx context.Context) error {
	if err := storage.Delete(h.msg.Fid); err != nil {
		logrus.Errorf("Delete Fid: %d Error: %s", h.msg.Fid, err)
		return err
	}
	return nil
}