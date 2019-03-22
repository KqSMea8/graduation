package handler

import (
	"context"
	"encoding/json"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mq"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type DeleteHandler struct {
	*CommonHandler
	Fid int64
}

func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{
		CommonHandler: NewCommonHandler(),
	}
}

func (h *DeleteHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, 200)
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
	if h.Fid, err = strconv.ParseInt(r.FormValue(constdef.Param_Fid), 10, 64); err != nil {
		logrus.Errorf("parse fid Error: %s", err)
		return err
	}
	return
}

// 事务 + 并发
func (h *DeleteHandler) delete_(ctx context.Context) (err error) {
	db := mysql.FileMySQL.Conn.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()

	// TODO 删除 redis 缓存
	if err = mysql.FileMySQL.Delete(db, h.Fid); err != nil {
		return
	}

	go storage.Delete(h.Fid)

	err = h.PublishDeleteEvent(ctx)
	return
}

func (h *DeleteHandler) PublishDeleteEvent(ctx context.Context) (err error) {
	msg := &model.DeleteFileEvent{
		Uid: h.UserId,
		Fid: h.Fid,
		Timestamp: time.Now().Unix(),
	}
	b, _ := json.Marshal(msg)
	return mq.Publish(constdef.DeleteFileEventTopic, b)
}

func (h *DeleteHandler) genResponse(out http.ResponseWriter, statusCode int) {
	out.WriteHeader(statusCode)
}