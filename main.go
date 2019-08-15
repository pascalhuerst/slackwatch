package main

import (
	"log"
	"os"
	"time"

	"slackwatch/slackwatch"
)

func main() {
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
