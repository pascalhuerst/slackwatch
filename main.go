package main

import (
	"os"
	"time"

	"slackwatch/slackwatch"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
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
			log.Print("Recovered from panic:", r)
		}
	}()

	config := slackwatch.NewConfigFromFile(os.Getenv("HOME") + "/.slackwatch")
	s := slackwatch.New(config)
	s.Run()
}
