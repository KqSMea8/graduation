package main

import (
	"github.com/g10guang/graduation/dal/mq"
	"github.com/g10guang/graduation/tools"
	"github.com/g10guang/graduation/write_api/handler"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	defer mq.StopNsqProducer()
	var err error
	initHttpHandler()
	if err = http.ListenAndServe("0.0.0.0:10001", nil); err != nil {
		logrus.Panicf("http.ListenAndServe Error: %s", err)
	}
	logrus.Infof("Main goroutine exit")
}

func initHttpHandler() {
	logrus.Info("Init HttpHandler")
	http.HandleFunc("/post", post)
	http.HandleFunc("/update", delete_)
}

// Restful interface

func post(out http.ResponseWriter, r *http.Request) () {
	r.ParseForm()
	//logrus.Infof("PostForm: %+v", r.PostForm)
	out.WriteHeader(200)
	out.Write([]byte("hello world"))
	return
	ctx := tools.NewCtxWithLogID()
	h := handler.NewPostHandler()
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("post Error: %s", err.Error())
	}
}

func delete_(out http.ResponseWriter, r *http.Request) {
	ctx := tools.NewCtxWithLogID()
	h := handler.NewDeleteHandler()
	err := h.Handle(ctx, out, r)
	if err != nil {
		logrus.Errorf("delete Error: %s", err.Error())
	}
}
