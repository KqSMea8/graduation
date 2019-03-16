package handler

import (
	"context"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
)

type PostHandler struct {
	*CommonHandler
	File       multipart.File
	FileHeader *multipart.FileHeader
}

func NewPostHandler(r *http.Request) *PostHandler {
	var err error
	h := &PostHandler{
		CommonHandler: NewCommonHandler(r),
	}
	h.File, h.FileHeader, err = r.FormFile(constdef.Param_File)
	if err != nil {
		logrus.Errorf("Get Upload File Error: %s", err.Error())
		return nil
	}
	return h
}

func (h *PostHandler) Handle(ctx context.Context, out http.ResponseWriter, r *http.Request) error {
	fid := tools.GenID()
	return nil
}

