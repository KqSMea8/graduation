package handler

import (
	"context"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/read_api/loader"
	"github.com/g10guang/graduation/tools"
	"net/http"
	"time"
)

type GetHandler struct {
	*CommonHandler
}

func NewGetHandler() *GetHandler {
	return &GetHandler{}
}

func (h *GetHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, 400)
		return err
	}

	// 1、获取图片元信息
	// 2、获取图片二进制内容
	jobmgr := tools.NewJobMgr(time.Second)
	jobmgr.AddJob(loader.NewFileMetaLoader(h.Fid))
	jobmgr.AddJob(loader.NewFileContentLoader(h.Fid, storage))
	if err = jobmgr.Start(ctx); err != nil {
		h.genResponse(out, 500)
		return
	}
	if result := jobmgr.GetResult(loader.LoaderName_FileMeta); result.Result != nil {
		switch v := result.Result.(type) {
		case model.File:
			h.FileMeta = v
		}
	}

	if result := jobmgr.GetResult(loader.LoaderName_FileContent); result.Result != nil {
		switch v := result.Result.(type) {
		case []byte:
			h.Bytes = v
		}
	}

	// 返回正常
	h.genResponse(out, 200)
	return
}

func (h *GetHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	return
}

func (h *GetHandler) genResponse(out http.ResponseWriter, statusCode int) {
	out.WriteHeader(statusCode)
	if statusCode == 200 {
		out.Header().Set("filename", h.FileMeta.Name)
		out.Write(h.Bytes)
	}
}
