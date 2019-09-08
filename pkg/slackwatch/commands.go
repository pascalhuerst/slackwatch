package slackwatch

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

func (s *Slackwatch) processCommand(m Message) bool {
	if !strings.HasPrefix(m.Text, "!") {
		return false
	}

	command := strings.Fields(m.Text)[0]
	switch command {
	case "!arm":
		s.armed = true
		logrus.Info("Armed")
		s.sendStatus(m.ChannelID)
	case "!disarm":
		s.armed = false
		logrus.Info("Disarmed")
		s.sendStatus(m.ChannelID)
	case "!verbose":
		s.outputAll = true
		logrus.Info("Verbose Set")
		s.sendStatus(m.ChannelID)
	case "!quiet":
		s.outputAll = false
		logrus.Info("Quiet Set")
		s.sendStatus(m.ChannelID)
	case "!chanls":
		s.sendChannels(m.ChannelID)
	case "!chanadd":
		s.addChannel(m)
		s.sendChannels(m.ChannelID)
	case "!chanrm":
		s.rmChannel(m)
		s.sendChannels(m.ChannelID)
	case "!status":
		s.sendStatus(m.ChannelID)
	case "!help":
		s.sendHelp(m.ChannelID)
	default:
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("Invalid Command", m.ChannelID))
	}

	return true
}

func (s *Slackwatch) addChannel(m Message) {
	channel := strings.TrimPrefix(m.DetokenizedText(), "!chanadd ")
	s.interestingChan = append(s.interestingChan, channel)
}

func (s *Slackwatch) rmChannel(m Message) {
	channel := strings.TrimPrefix(m.DetokenizedText(), "!chanrm ")
	newChans := make([]string, 0)
	for _, c := range s.interestingChan {
		if c != channel {
			newChans = append(newChans, c)
		}
	}
	s.interestingChan = newChans
}

func (s *Slackwatch) sendChannels(channelID string) {
	message := strings.Join(s.interestingChan, ", ")
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage("Current watched channels: "+message, channelID))
}

func (s *Slackwatch) sendStatus(channelID string) {
	status := fmt.Sprintf("armed: %v, outputAll: %v", s.armed, s.outputAll)
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(status, channelID))
}

func (s *Slackwatch) sendHelp(channelID string) {
	text := "!arm !disarm !verbose !quiet !status !help\n" +
		"For watched channels: !chanls !chanadd <channel w/o #>, !chanrm <channel w/o #>"
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channelID))
}
