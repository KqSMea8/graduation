package handler

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
)

type HeadHandler struct {
	*CommonHandler
}

func NewHeadHandler() *HeadHandler {
	return &HeadHandler{}
}

func (h *HeadHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) (err error) {
	if err = h.parseParams(ctx, r); err != nil {
		h.genResponse(out, 400)
		return
	}
	// TODO 增加 redis 缓存的访问
	if err = h.loadFileMeta(ctx); err != nil {
		h.genResponse(out, 500)
		return
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
