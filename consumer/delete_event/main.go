package main

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer clean()

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(constdef.DeleteFileEventTopic, "delete", config)
	consumer.ChangeMaxInFlight(200)
	consumer.AddHandler(nsq.HandlerFunc(delete_))

	if err = consumer.ConnectToNSQLookupds([]string{constdef.NsqLookupdIp}); err != nil {
		logrus.Panicf("ConnectToNSQLookupds Error: %s", err)
		panic(err)
	}

	shutdown := make(chan os.Signal, 2)
	signal.Notify(shutdown, syscall.SIGINT)

	for {
		select {
		case <-consumer.StopChan:
			goto exit
		case <-shutdown:
			consumer.Stop()
			goto exit
		}
	}

exit:
	logrus.Infof("consumer exit")
}

func clean() {

}
