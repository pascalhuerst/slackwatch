package main

import (
  "fmt"
  "time"
  "os"

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
      fmt.Println("Recovered from panic:", r)
    }
  }()

  config := slackwatch.NewConfigFromFile(os.Getenv("HOME") + "/.slackwatch")
  s := slackwatch.New(config)
  s.Run()
}
