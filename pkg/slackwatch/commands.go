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

	switch m.Text {
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
	case "!status":
		s.sendStatus(m.ChannelID)
	case "!help":
		s.sendHelp(m.ChannelID)
	default:
		s.rtm.SendMessage(s.rtm.NewOutgoingMessage("Invalid Command", m.ChannelID))
	}

	return true
}

func (s *Slackwatch) sendStatus(channelID string) {
	status := fmt.Sprintf("armed: %v, outputAll: %v", s.armed, s.outputAll)
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(status, channelID))
}

func (s *Slackwatch) sendHelp(channelID string) {
	text := "!arm !disarm !verbose !quiet !status !help"
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channelID))
}
