package main

import (
	"github.com/g10guang/graduation/read_api/handler"
	"github.com/g10guang/graduation/tools"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	var err error
	initHttpHandler()
	if err = http.ListenAndServe("0.0.0.0:10002", nil); err != nil {
		logrus.Panicf("http.ListenAndServe Error: %s", err)
	}
	logrus.Infof("Main goroutine exit")
}

func initHttpHandler() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/head", head)
}

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
