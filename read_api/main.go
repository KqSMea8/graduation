package main

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	var err error
	initHttpHandler()
	if err = http.ListenAndServe("0.0.0.0:9806", nil); err != nil {
		logrus.Panicf("http.ListenAndServe Error: %s", err)
	}
	logrus.Infof("Main goroutine exit")
}

func initHttpHandler() {
	http.HandleFunc("/get", get)
	http.HandleFunc("/head", head)
}
