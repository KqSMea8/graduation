package main

import (
	"encoding/json"
	"errors"
	"github.com/g10guang/graduation/constdef"
	"github.com/g10guang/graduation/consumer/post_event/handler"
	"github.com/g10guang/graduation/model"
	"github.com/g10guang/graduation/tools"
	"github.com/nsqio/go-nsq"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	defer clean()

	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(constdef.PostFileEventTopic, "compress", config)
	if err != nil {
		panic(err)
	}
	consumer.ChangeMaxInFlight(200)
	consumer.AddHandler(nsq.HandlerFunc(compress))
	if err = consumer.ConnectToNSQLookupds([]string{constdef.NsqLookupdAddr}); err != nil {
		logrus.Panicf("ConnectToNSQLookupds Error: %s", err)
		panic(err)
	}
	//if err = consumer.ConnectToNSQDs([]string{constdef.NsqdAddr}); err != nil {
	//	logrus.Panicf("ConnectToNSQDs Error: %s", err)
	//	panic(err)
	//}

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

// 将图片转化为 jpeg/png 格式
func compress(message *nsq.Message) error {
	logrus.Infof("message: %+v", message)
	msg := parsePostFileEventMsg(message.Body)
	if msg == nil {
		return errors.New("message error")
	}
	h := handler.NewCompressHandler(msg)
	if err := h.Handle(tools.NewCtxWithLogID()); err != nil {
		logrus.Errorf("CompressHandler Error: %s", err)
		return err
	}
	logrus.Infof("CompressHandler Success")
	return nil
}

func parsePostFileEventMsg(body []byte) *model.PostFileEvent {
	logrus.Infof("post_file event message: %s", string(body))
	m := new(model.PostFileEvent)
	if err := json.Unmarshal(body, m); err != nil {
		logrus.Errorf("PostFileEvent message Error: %s", err)
		return nil
	}
	return m
}
