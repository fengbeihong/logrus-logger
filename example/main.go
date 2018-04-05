package main

import (
	"github.com/sirupsen/logrus"
	"mixotc.com/mixotc/logrus-logger"
)

func main() {
	log.InitDefaultMyLogger(logrus.DebugLevel, "test.log")

	log.Debug("test debug")
	log.Info("test info")
	log.Error("test error")
}
