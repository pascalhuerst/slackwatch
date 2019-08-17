package slackwatch_test

import (
	"os"
	"slackwatch/slackwatch"
)

func Example() {
	config := slackwatch.NewConfigFromFile(os.Getenv("HOME") + "/.slackwatch")
	s := slackwatch.New(config)
	s.Run()
}

func ExampleDefaultAction() {
	_ = slackwatch.DefaultAction{Command: "/usr/bin/afplay", Args: "klaxon.wav"}
	_ = slackwatch.DefaultAction{URL: "https://example.com/trigger_action"}
	_ = slackwatch.DefaultAction{URL: "https://example.com/other_endpoint", Body: "Slackwatch Triggered"}
}

func ExampleNewConfigFromFile() {
	config := slackwatch.NewConfigFromFile(os.Getenv("HOME") + "/.slackwatch")
	s := slackwatch.New(config)
	s.Run()
}
