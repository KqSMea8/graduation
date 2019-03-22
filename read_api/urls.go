package main

import (
	"github.com/g10guang/graduation/read_api/handler"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"net/http"
)

func get(out http.ResponseWriter, r *http.Request) {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewGetHandler()
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("post Error: %s", err.Error())
	}

}

func head(out http.ResponseWriter, r *http.Request) {

}
