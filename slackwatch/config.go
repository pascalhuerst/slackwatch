package slackwatch

import (
  "encoding/json"
  "io/ioutil"

)

type SlackwatchConfig struct {
  SlackToken string
  Actions []*Action
}

// https://stackoverflow.com/questions/28817992/how-to-set-bool-pointer-to-true-in-struct-literal
func NewConfigFromFile(path string) SlackwatchConfig {
  dat, err := ioutil.ReadFile(path)
  if err != nil {
    panic(err)
  }

  var c SlackwatchConfig
  err = json.Unmarshal(dat, &c)
  if err != nil {
    panic(err)
  }

  return c
}
