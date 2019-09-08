package slackwatch

import (
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/nlopes/slack"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestSlackwatch(t *testing.T) {
	// Setup
	hook := test.NewGlobal()
	filepath := "/tmp/slackwatchtest-" + strconv.Itoa(os.Getpid())
	config := Config{
		SlackToken: "xoxp-123",
		Actions: []Action{
			&CommandAction{Command: "/usr/bin/touch", Args: filepath},
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

	// Verbose/Quiet
	if sw.verbose {
		t.Error("verbose began as true")
	}

	sw.messageReceived(newMessage("123", "D231", "U123", "!verbose", sw))
	if !sw.verbose {
		t.Error("verbose failed to turn on")
	}
	if hook.LastEntry().Message != "Verbose Set" {
		t.Error("verbose to generate verbose message in log")
	}

	sw.messageReceived(newMessage("123", "D231", "U123", "!quiet", sw))
	if sw.verbose {
		t.Error("verbose failed to turn off")
	}
	if hook.LastEntry().Message != "Quiet Set" {
		t.Error("failed to generate quiet message in log")
	}

	// Armed/Disarmed
	if !sw.armed {
		t.Error("armed began as false")
	}

	sw.messageReceived(newMessage("123", "D231", "U123", "!disarm", sw))
	if sw.armed {
		t.Error("armed failed to turn off")
	}
	if hook.LastEntry().Message != "Disarmed" {
		t.Error("failed to generate disarmed message in log")
	}

	sw.messageReceived(newMessage("123", "D231", "U123", "!arm", sw))
	if !sw.armed {
		t.Error("armed failed to turn on")
	}
	if hook.LastEntry().Message != "Armed" {
		t.Error("failed to generate armed message in log")
	}

	// todo: we need to mock s.rtm.NewOutgoingMessage so we can verify these
	sw.messageReceived(newMessage("123", "D231", "U123", "!status", sw))
	sw.messageReceived(newMessage("123", "D231", "U123", "!help", sw))
	sw.messageReceived(newMessage("123", "D231", "U123", "!invalid", sw))

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

	sw.messageReceived(newMessage("123", "C123", "U321", "hello", sw))
	if hook.LastEntry().Message != "  [general] <user2> hello" {
		t.Error("failed to generate proper uninteresting message log")
	}
	sw.messageReceived(newMessage("123", "C123", "U321", "hello <@U123>", sw))
	if hook.LastEntry().Message != "* [general] <user2> hello @user1" {
		t.Error("failed to generate proper interesting message log")
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
