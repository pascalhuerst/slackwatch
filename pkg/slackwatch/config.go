package slackwatch

import (
	"encoding/json"
	"io/ioutil"
)

// Config is passed to the New constructor.
type Config struct {
	SlackToken string
	Actions    []Action
}

type jsonConfig struct {
	SlackToken string
	Actions    []defaultAction
}

type defaultAction struct {
	CommandAction
	URLAction
}

// NewConfigFromFile parses a json formatted config file and returns a Config.
func NewConfigFromFile(path string) Config {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var j jsonConfig
	err = json.Unmarshal(dat, &j)
	if err != nil {
		panic(err)
	}

	var action []Action
	for _, a := range j.Actions {
		if a.Command != "" {
			action = append(action, CommandAction{Command: a.Command, Args: a.Args})
		}
		if a.URL != "" {
			action = append(action, URLAction{URL: a.URL, Body: a.Body, ContentType: a.ContentType})
		}
	}

	return Config{
		SlackToken: j.SlackToken,
		Actions:    action,
	}
}
