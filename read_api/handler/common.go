package handler

import (
	"context"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/dal/mysql"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/store"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type CommonHandler struct {
	Fid int64
	FileMeta model.File
	Bytes []byte
}

func NewCommonHandler() *CommonHandler {
	return &CommonHandler{}
}

func (h *CommonHandler) parseParams(ctx context.Context, r *http.Request) (err error) {
	h.Fid, err = strconv.ParseInt(r.FormValue(constdef.Param_Fid), 10, 64)
	if err != nil {
		logrus.Errorf("parse Fid Error: %s", err)
	}
	return err
}

func (h *CommonHandler) loadFileMeta(ctx context.Context) (err error) {
	h.FileMeta, err = mysql.FileMySQL.Get(h.Fid)
	return
}

func (h *CommonHandler) loadFileContent(ctx context.Context) (err error) {
	h.Bytes, err = storage.Read(h.Fid)
	return
}

var storage store.Storage

func init() {
	storage = store.NewLocalStorage()
}
