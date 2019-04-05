package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/read_api/loader"
	"github.com/g10guang/graduation/tools"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"time"
)

type MetaHandler struct {
	*CommonHandler
	FileMetas []*model.File
	Fids      []int64
}

func NewMetaHandler() *MetaHandler {
	return &MetaHandler{
		CommonHandler: NewCommonHandler(),
	}
}

func (h *MetaHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, 400)
		return
	}
	jobmgr := tools.NewJobMgr(time.Second)
	jobmgr.AddJob(loader.NewFileMetaLoader(h.Fids))
	if err = jobmgr.Start(ctx); err != nil {
		h.genResponse(out, 500)
		return
	}
	if result := jobmgr.GetResult(loader.LoaderName_FileMeta); result.Result != nil {
		if v, ok := result.Result.(map[int64]*model.File); ok {
			for _, v := range v {
				h.FileMetas = append(h.FileMetas, v)
			}
		}
	}
	h.genResponse(out, 200)
	return
}

func (h *MetaHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	if err = h.CommonHandler.parseParams(ctx, r); err != nil {
		return err
	}
	// 解析 post 表单中提交的 fid slice
	fids := r.PostForm[constdef.Param_Fid]
	h.Fids = make([]int64, len(fids))
	for i, fid := range fids {
		if h.Fids[i], err = strconv.ParseInt(fid, 10, 64); err != nil {
			logrus.Errorf("strconv.ParseInt fid: %s Error: %s", fid, err)
			return err
		}
	}
	if len(h.Fids) == 0 {
		logrus.Errorf("empty meta fids")
		return errors.New("empty fids")
	}
	if len(h.Fids) > 100 {
		logrus.Errorf("too much fids len=%d", len(h.Fids))
		return fmt.Errorf("too much fids len=%s", len(h.Fids))
	}
	return
}

func (h *MetaHandler) genResponse(out http.ResponseWriter, statusCode int) {
	if statusCode == 200 {
		b, err := json.Marshal(h.FileMetas)
		if err != nil {
			logrus.Errorf("Marshal FileMetas Error: %s", err)
			out.WriteHeader(500)
			return
		}
		if _, err = out.Write(b); err != nil {
			logrus.Errorf("write http response Error: %s", err)
			out.WriteHeader(500)
			return
		}
		out.WriteHeader(200)
	} else {
		out.WriteHeader(statusCode)
	}
}
