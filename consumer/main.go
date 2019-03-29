package main

import (
	"github.com/nsqio/go-nsq"
	"image/jpeg"
)

func main() {
	config := nsq.NewConfig()
	// TODO 修改消息队列的 topic
	q, err := nsq.NewConsumer("topic", "channel", config)
	if err != nil {
		panic(err)
	}
	q.AddHandler(nsq.HandlerFunc(jpegCompress))
	q.AddHandler(nsq.HandlerFunc(pngCompress))
}

