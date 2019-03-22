package handler

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/dal/mq"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
)

type PostHandler struct {
	*CommonHandler
	File       multipart.File
	FileHeader *multipart.FileHeader
	FileMeta   *model.File
	Bytes      []byte
}

func NewPostHandler(r *http.Request) *PostHandler {
	h := &PostHandler{
		CommonHandler: NewCommonHandler(r),
	}
	return h
}

// TODO 添加 generate response
func (h *PostHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	var err error
	if err = h.parseParams(ctx, r); err != nil {
		return err
	}
	h.BuildFileMeta()
	if err = h.SaveFile(ctx); err != nil {
		return err
	}
	return nil
}

func (h *PostHandler) parseParams(ctx context.Context, r *http.Request) error {
	var err error
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	h.File, h.FileHeader, err = r.FormFile(constdef.Param_File)
	if err != nil {
		logrus.Errorf("Get Upload FileBytes Error: %s", err.Error())
		return err
	}
	h.Bytes = make([]byte, int(h.FileHeader.Size))
	if _, err = h.File.Read(h.Bytes); err != nil {
		logrus.Errorf("Read File Error: %s", err)
		return err
	}
	return nil
}

// 有可能有网络错误，唯一键冲突等
// 因为整个过程涉及到不少网络操作，所以需要使用事务，免得数据库中插入了无用记录
// TODO 以下操作需要改为并发，而且需要考虑处理失败的场景
func (h *PostHandler) SaveFile(ctx context.Context) (err error) {
	logrus.Debugf("SaveFile fid: %d", h.FileMeta.Fid)
	db := mysql.FileMySQL.Conn.Begin()
	defer func() {
		if err == nil {
			db.Commit()
		} else {
			db.Rollback()
		}
	}()
	// 保存文件元信息
	if err = mysql.FileMySQL.SaveFileMeta(db, h.FileMeta); err != nil {
		return
	}
	// 保存文件内容
	err = storage.Write(h.FileMeta.Fid, h.Bytes)
	if err != nil {
		logrus.Errorf("Persistent File Error: %s", err)
		return
	}
	err = h.PublishPostFileEvent()
	return
}

// 发送一条消息到消息队列
func (h *PostHandler) PublishPostFileEvent() error {
	msg := &model.PostFileEvent{
		Fid: h.FileMeta.Fid,
		Uid: h.FileMeta.Uid,
		Timestamp: time.Now().Unix(),
	}
	b, _ := json.Marshal(msg)
	return mq.Publish(constdef.PostFileEventTopic, b)
}

func (h *PostHandler) BuildFileMeta() {
	h.FileMeta = &model.File{
		Uid:  h.UserId,
		Fid:  tools.GenID().Int64(),
		Name: h.FileHeader.Filename,
		Size: h.FileHeader.Size,
		Md5:  h.CalculateMd5(),
	}
}

func (h *PostHandler) CalculateMd5() string {
	return fmt.Sprintf("%x", md5.Sum(h.Bytes))
}
