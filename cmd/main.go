package main

import (
	"os"
	"time"

	"github.com/mikegrb/slackwatch/pkg/slackwatch"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "15:04:05",
	})

	for {
		runSlackwatch()
		time.Sleep(10 * time.Second)
	}
}

func runSlackwatch() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Errorf("Recovered from panic: %v", r)
		}
	}()

	config := slackwatch.NewConfigFromFile(os.Getenv("HOME") + "/.slackwatch")
	s := slackwatch.New(config)
	s.Run()
}
