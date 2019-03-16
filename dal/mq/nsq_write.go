package mq

import (
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
)

// when a file posted, send a message to nsq
var nsqProducer *nsq.Producer

func init() {
	var err error
	const addr = "127.0.0.1:4301"
	config := nsq.NewConfig()
	nsqProducer, err = nsq.NewProducer(addr, config)
	if err != nil {
		panic(err)
	}
}

func StopNsqProducer() {
	nsqProducer.Stop()
}

func Publish(topic string, content []byte) error {
	logrus.Errorf("Send Topic %s content %s", string(content))
	err := nsqProducer.Publish(topic, content)
	logrus.Errorf("Send Topic %s Error: %s", topic, err)
	return err
}
