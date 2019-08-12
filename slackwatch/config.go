package slackwatch

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	SlackToken string
	Actions    []*Action
}

func NewConfigFromFile(path string) Config {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var c Config
	err = json.Unmarshal(dat, &c)
	if err != nil {
		panic(err)
	}

	return c
}
