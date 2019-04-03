package handler

import (
	"context"
	"encoding/json"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mq"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type DeleteHandler struct {
	*CommonHandler
	Fids []int64
}

func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{
		CommonHandler: NewCommonHandler(),
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, 400)
		return err
	}
	if err = h.delete_(ctx); err != nil {
		h.genResponse(out, 500)
		return err
	}
	h.genResponse(out, 200)
	return nil
}

func (h *DeleteHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	fids := r.PostForm[constdef.Param_Fid]
	h.Fids = make([]int64, len(fids))
	for i, fid := range fids {
		h.Fids[i], err = strconv.ParseInt(fid, 10, 64)
		if err != nil {
			logrus.Errorf("parse fid Error: %s", err)
		}
	}
	if len(h.Fids) == 0 {
		logrus.Errorf("empty delete fids")
		return errors.New("empty delete fids")
	}
	logrus.Infof("fids: %+v", h.Fids)
	return
}

// 事务 + 并发
func (h *DeleteHandler) delete_(ctx context.Context) (err error) {
	if err = mysql.FileMySQL.MDelete(nil, h.Fids); err == gorm.ErrRecordNotFound {
		logrus.Errorf("Uid: %d delete fid: %v not exist", h.UserId, h.Fids)
		return err
	}
	// 发送消息队列异步化
	go h.PublishDeleteEvent(ctx)
	return
}

func (h *DeleteHandler) PublishDeleteEvent(ctx context.Context) (err error) {
	for _, fid := range h.Fids {
		msg := &model.DeleteFileEvent{
			Uid:       h.UserId,
			Fid:       fid,
			Timestamp: time.Now().Unix(),
		}
		b, _ := json.Marshal(msg)
		if err = mq.PublishNsq(constdef.DeleteFileEventTopic, b); err != nil {
			logrus.Errorf("Publish Delete Event Error: %s", err)
		}
	}
	return err
}

func (h *DeleteHandler) genResponse(out http.ResponseWriter, statusCode int) {
	out.WriteHeader(statusCode)
}
