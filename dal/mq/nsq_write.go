package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// when a file posted, send a message to nsq
var nsqProducer *nsq.Producer

func init() {
	var err error
	const addr = "10.8.118.15:10004"
	config := nsq.NewConfig()
	nsqProducer, err = nsq.NewProducer(addr, config)
	if err != nil {
		panic(err)
	}
}

func StopNsqProducer() {
	nsqProducer.Stop()
}

func PublishNsq(topic string, content []byte) error {
	err := nsqProducer.Publish(topic, content)
	if err != nil {
		logrus.Errorf("PublishNsq Error: %s topic: %s content: %s", err, topic, string(content))
	}
	return err
}
