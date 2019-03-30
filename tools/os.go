package tools

import (
	"github.com/g10guang/graduation/constdef"
	"github.com/sirupsen/logrus"
	"os"
)

func IsProductEnv() bool {
	_, exists := os.LookupEnv(constdef.ENV_ProductEnv)
	return exists
}

func InitLog() {
	if !IsProductEnv() {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
	}
}
