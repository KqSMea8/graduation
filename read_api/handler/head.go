package handler

import (
	"context"
	"encoding/json"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/read_api/loader"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type HeadHandler struct {
	*CommonHandler
	FileMeta *model.File
}

func NewHeadHandler() *HeadHandler {
	return &HeadHandler{}
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
	out.WriteHeader(statusCode)
	if statusCode == 200 {
		b, _ := json.Marshal(h.FileMeta)
		if _, err := out.Write(b); err != nil {
			logrus.Errorf("Http Write File Error: %s", err)
		}
	}
}
