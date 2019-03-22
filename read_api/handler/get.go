package handler

import (
	"context"
	"github.com/pkg/errors"
	"net/http"
	"sync"
	"sync/atomic"
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

	// TODO 这里 WaitGroup 可以改为 JobMgr
	// 1、获取图片元信息
	// 2、获取图片二进制内容
	var wg sync.WaitGroup
	wg.Add(2)
	var errNum int32
	go func() {
		defer wg.Done()
		// TODO 获取元数据可以增加 redis 缓存
		if err := h.loadFileContent(ctx); err != nil {
			atomic.AddInt32(&errNum, 1)
		}
	}()

	go func() {
		defer wg.Done()
		if err := h.loadFileMeta(ctx); err != nil {
			atomic.AddInt32(&errNum, 1)
		}
	}()

	wg.Wait()

	// 返回错误
	if errNum > 0 {
		h.genResponse(out, 500)
		return errors.New("coroutine load data error")
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
