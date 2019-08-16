package slackwatch

import (
	"os"
	"testing"

	"github.com/nlopes/slack"
)

func TestIsInteresting(t *testing.T) {
	sw := slackwatch{me: &slack.UserDetails{ID: "U123"}}
	tables := []struct {
		d string
		m Message
		r bool
	}{
		{"channel message", Message{Channel: "general", Text: "hi guys", sw: &sw}, false},
		{"private message", Message{Channel: "DM", Text: "hi", sw: &sw}, true},
		{"channel mention", Message{Channel: "general", Text: "hi <@U123>", sw: &sw}, true},
		{"mpdm message", Message{Channel: "mpdm-usr1-usr2-usr3-1", Text: "hi", sw: &sw}, true},
	}

	for _, table := range tables {
		result := table.m.IsInteresting()
		if result != table.r {
			t.Errorf("isInteresting failed test '%s', got: %v, want: %v", table.d, result, table.r)
		}
	}
}

func TestIsFromMe(t *testing.T) {
	sw := slackwatch{me: &slack.UserDetails{ID: "U123"}}
	tables := []struct {
		d string
		m Message
		r bool
	}{
		{"message from someone else", Message{UserID: "U321", sw: &sw}, false},
		{"message from me", Message{UserID: "U123", sw: &sw}, true},
	}
	for _, table := range tables {
		result := table.m.IsFromMe()
		if result != table.r {
			t.Errorf("isFromMe failed test '%s', got: %v, want: %v", table.d, result, table.r)
		}
	}
}

func TestTS(t *testing.T) {
	os.Setenv("TZ", "Etc/UTC")

	tables := []struct {
		input  string
		output string
	}{
		{"1565957851", "12:17:31"},
		{"1565913600", "00:00:00"},
	}
	for _, table := range tables {
		result := ts(table.input)
		if result != table.output {
			t.Errorf("ts failed test, got: %v, want: %v", result, table.output)
		}
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	ts("one")
}

func TestMessage(t *testing.T) {
	sw := slackwatch{
		userLookup: map[string]string{
			"U123": "user1",
			"U321": "user2",
		},
		conversationLookup: map[string]string{
			"C123": "general",
			"C321": "dev",
		},
		me: &slack.UserDetails{ID: "U123"},
	}

	tables := []struct {
		d string
		m Message
		r string
	}{
		{"message 1", newMessage("0", "C123", "U123", "hi gals", &sw), "  [general] <user1> hi gals"},
		{"message 2", newMessage("0", "C123", "U321", "hi guys", &sw), "  [general] <user2> hi guys"},
		{"message 3", newMessage("0", "C123", "U321", "hi <@U123>", &sw), "* [general] <user2> hi @user1"},
		{"message 4", newMessage("0", "", "", "hi", &sw), "  [] <> hi"},
	}

	for _, table := range tables {
		result := table.m.String()
		if result != table.r {
			t.Errorf("Failed test '%s', got: '%v', want: '%v'", table.d, result, table.r)
		}
	}

}
