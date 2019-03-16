package handler

import (
	"context"
	"errors"
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type Handler interface {
	Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error
}

type CommonHandler struct {
	UserId int64
}

func NewCommonHandler(r *http.Request) *CommonHandler {
	h := &CommonHandler{}
	h.UserId, _ = strconv.ParseInt(r.FormValue(constdef.Param_Uid), 10, 64)
	return h
}

func (h *CommonHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	return nil
}

func (h *CommonHandler) CheckParams(ctx context.Context) error {
	if h.UserId == 0 {
		logrus.Errorf("uid == 0")
		return errors.New("uid == 0")
	}
	return nil
}
