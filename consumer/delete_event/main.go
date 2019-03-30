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

func delete_(message *nsq.Message) error {
	msg := parseDeleteFileEventMsg(message.Body)
	if msg.Fid == 0 || msg.Uid == 0 {
		return fmt.Errorf("invalid message: %s", string(message.Body))
	}
	h := handler.NewDeleteStorageHandler(msg)
	if err := h.Handle(tools.NewCtxWithLogID()); err != nil {
		logrus.Errorf("Delete Storage Error: %s", err)
		return err
	}
	return nil
}

func parseDeleteFileEventMsg(body []byte) *model.DeleteFileEvent {
	msg := &model.DeleteFileEvent{}
	err := json.Unmarshal(body, msg)
	if err != nil {
		logrus.Errorf("unmarshal DeleteFileEvent message Error: %s", err)
	}
	return msg
}
