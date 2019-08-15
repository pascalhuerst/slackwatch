package slackwatch

import (
	"encoding/json"
	"io/ioutil"
)

// Config is passed to the New constructor.
type Config struct {
	SlackToken string
	Actions    []*Action
}

// NewConfigFromFile parses a json formatted config file and returns a Config.
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
