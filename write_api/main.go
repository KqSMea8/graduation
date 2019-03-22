package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
)


func main() {
	var err error
	initHttpHandler()
	if err = http.ListenAndServe("0.0.0.0:9807", nil); err != nil {
		logrus.Panicf("http.ListenAndServe Error: %s", err)
	}
	logrus.Infof("Main goroutine exit")
}

func initHttpHandler() {
	logrus.Info("Init HttpHandler")
	http.HandleFunc("/post", post)
	http.HandleFunc("/update", delete_)
}

