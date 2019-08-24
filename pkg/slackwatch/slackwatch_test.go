package slackwatch

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/nlopes/slack"
)

func TestSlackwatch(t *testing.T) {
	// Setup
	filepath := "/tmp/slackwatchtest-" + strconv.Itoa(os.Getpid())
	config := Config{
		SlackToken: "xoxp-123",
		Actions: []*DefaultAction{
			&DefaultAction{Command: "/usr/bin/touch", Args: filepath},
		},
	}

	sw := New(config)
	sw.me = &slack.UserDetails{ID: "U123"}
	sw.userLookup["U123"] = "user1"
	sw.userLookup["U321"] = "user2"
	sw.userLookup["U231"] = "user3"
	sw.conversationLookup["C123"] = "general"
	sw.conversationLookup["C321"] = "dev"
	sw.conversationLookup["D231"] = "DM"
	sw.conversationLookup["G123"] = "mpdm--user1--user2--user3--1"

	if sw.outputAll {
		t.Error("outputAll began as true")
	}
	sw.messageReceived(newMessage("123", "D231", "U123", "!verbose", sw))
	if !sw.outputAll {
		t.Error("outputAll failed to turn on")
	}
	sw.messageReceived(newMessage("123", "D231", "U123", "!quiet", sw))
	if sw.outputAll {
		t.Error("outputAll failed to turn off")
	}

	if !sw.armed {
		t.Error("armed began as false")
	}
	sw.messageReceived(newMessage("123", "D231", "U123", "!disarm", sw))
	if sw.outputAll {
		t.Error("armed failed to turn off")
	}
	sw.messageReceived(newMessage("123", "D231", "U123", "!arm", sw))
	if !sw.armed {
		t.Error("armed failed to turn on")
	}

	// todo: we need to mock s.rtm.NewOutgoingMessage so we can verify these
	sw.messageReceived(newMessage("123", "D231", "U123", "!status", sw))
	sw.messageReceived(newMessage("123", "D231", "U123", "!help", sw))
	sw.messageReceived(newMessage("123", "D231", "U123", "!invalid", sw))
	sw.messageReceived(newMessage("123", "D231", "U123", "!help", sw))

	sw.messageReceived(newMessage("123", "D231", "U123", "hello", sw))
	time.Sleep(1 * time.Second)
	if !fileExists(filepath) {
		t.Error("touch file action failed to execute")
	} else {
		os.Remove(filepath)
	}

	sw.messageReceived(newMessage("123", "D231", "U123", "!disarm", sw))
	sw.messageReceived(newMessage("123", "D231", "U123", "hello", sw))
	time.Sleep(1 * time.Second)
	if fileExists(filepath) {
		t.Error("touch file action executed when disarmed")
		os.Remove(filepath)
	}

	sw.messageReceived(newMessage("123", "D231", "U123", "!verbose", sw))

	// TODO capture log.Print
	sw.messageReceived(newMessage("123", "C123", "U321", "hello", sw))

}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
