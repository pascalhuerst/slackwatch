package slackwatch

import (
	"testing"

	"github.com/nlopes/slack"
)

func TestIsInteresting(t *testing.T) {
	sw := slackwatch{me: &slack.UserDetails{ID: "U123"}}
	tables := []struct {
		d string
		m message
		r bool
	}{
		{"channel message", message{channel: "general", text: "hi guys", sw: &sw}, false},
		{"private message", message{channel: "DM", text: "hi", sw: &sw}, true},
		{"channel mention", message{channel: "general", text: "hi <@U123>", sw: &sw}, true},
		{"mpdm message", message{channel: "mpdm-usr1-usr2-usr3-1", text: "hi", sw: &sw}, true},
	}

	for _, table := range tables {
		result := table.m.isInteresting()
		if result != table.r {
			t.Errorf("isInteresting failed test '%s', got: %v, want: %v", table.d, result, table.r)
		}
	}
}
