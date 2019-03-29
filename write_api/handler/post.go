package handler

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/dal/mq"
	"github.com/g10guang/graduation/write_api/jobs"
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

func NewPostHandler() *PostHandler {
	h := &PostHandler{
		CommonHandler: NewCommonHandler(),
	}
	return h
}

func (h *PostHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	var err error
	if err = h.parseParams(ctx, r); err != nil {
		// 用户参数错误
		h.genResponse(out, 400)
		return err
	}
	h.BuildFileMeta()
	if err = h.SaveFile(ctx); err != nil {
		h.genResponse(out, 500)
		return err
	}
	h.genResponse(out, 200)
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
func (h *PostHandler) SaveFile(ctx context.Context) (err error) {
	logrus.Debugf("SaveFile fid: %d", h.FileMeta.Fid)
	db := mysql.FileMySQL.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
			go storage.Delete(h.FileMeta.Fid)
		} else {
			db.Commit()
		}
	}()

	jobmgr := tools.NewJobMgr(time.Second)
	jobmgr.AddJob(jobs.NewSaveFileMetaJob(h.FileMeta, db))
	jobmgr.AddJob(jobs.NewStoreFileJob(h.FileMeta.Fid, h.Bytes, storage))
	if err = jobmgr.Start(ctx); err != nil {
		logrus.Errorf("batch Job process Error: %s", err)
		return err
	}

	// 最后发送消息不能够并行，不然会增大消费者处理难度
	if	err = h.PublishPostFileEvent();  err != nil {
		logrus.Errorf("Send Nsq Error: %s", err)
		return err
	}

	return nil
}

// 发送一条消息到消息队列
func (h *PostHandler) PublishPostFileEvent() error {
	msg := &model.PostFileEvent{
		Fid: h.FileMeta.Fid,
		Uid: h.FileMeta.Uid,
		Timestamp: time.Now().Unix(),
	}
	b, _ := json.Marshal(msg)
	return mq.PublishNsq(constdef.PostFileEventTopic, b)
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

// 生成响应
func (h *PostHandler) genResponse(out http.ResponseWriter, statusCode int) {
	out.WriteHeader(statusCode)
}
