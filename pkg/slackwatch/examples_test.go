package slackwatch_test

import (
	"os"

	"github.com/mikegrb/slackwatch/pkg/slackwatch"
)

func Example() {
	config := slackwatch.Config{
		SlackToken: "xoxp-123",
		Actions: []slackwatch.Action{
			slackwatch.CommandAction{Command: "/usr/bin/afplay", Args: "klaxon.wav"},
			slackwatch.URLAction{URL: "https://example.com/trigger_action"},
			slackwatch.URLAction{URL: "https://example.com/other_endpoint", Body: "Slackwatch Triggered"},
		},
	}

	s := slackwatch.New(config)
	s.Run()
}

func ExampleNewConfigFromFile() {
	config := slackwatch.NewConfigFromFile(os.Getenv("HOME") + "/.slackwatch")
	s := slackwatch.New(config)
	s.Run()
}
