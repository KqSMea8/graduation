package handler

import (
	"context"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/read_api/loader"
	"github.com/g10guang/graduation/tools"
	"net/http"
	"strconv"
	"time"
)

type HeadHandler struct {
	*CommonHandler
	FileMeta *model.File
}

func NewHeadHandler() *HeadHandler {
	return &HeadHandler{
		CommonHandler: NewCommonHandler(),
	}
}

func (h *HeadHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, 400)
		return
	}
	jobmgr := tools.NewJobMgr(time.Second)
	jobmgr.AddJob(loader.NewFileMetaLoader(h.Fid))
	if err = jobmgr.Start(ctx); err != nil {
		h.genResponse(out, 500)
		return
	}
	if result := jobmgr.GetResult(loader.LoaderName_FileMeta); result.Result != nil {
		switch v := result.Result.(type) {
		case model.File:
			h.FileMeta = &v
		}
	}
	h.genResponse(out, 200)
	return
}

func (h *HeadHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	return
}

func (h *HeadHandler) genResponse(out http.ResponseWriter, statusCode int) {
	if statusCode == 200 {
		out.Header().Set("fid", strconv.FormatInt(h.FileMeta.Fid, 10))
		out.Header().Set("uid", strconv.FormatInt(h.FileMeta.Uid, 10))
		out.Header().Set("name", h.FileMeta.Name)
		out.Header().Set("size", strconv.FormatInt(h.FileMeta.Size, 10))
		out.Header().Set("md5", h.FileMeta.Md5)
		out.Header().Set("create_time", strconv.FormatInt(h.FileMeta.CreateTime.Unix(), 10))
		out.Header().Set("update_time", strconv.FormatInt(h.FileMeta.UpdateTime.Unix(), 10))
		out.WriteHeader(200)
	} else {
		out.WriteHeader(statusCode)
	}
}
