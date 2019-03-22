package main

import (
	"github.com/g10guang/graduation/tools"
	"github.com/g10guang/graduation/write_api/handler"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Restful interface

func post(out http.ResponseWriter, r *http.Request) () {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewPostHandler(r)
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("post Error: %s", err.Error())
	}
}

func delete_(out http.ResponseWriter, r *http.Request) {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewDeleteHandler(r)
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("delete Error: %s", err.Error())
	}
}
