package slackwatch

import (
	"fmt"
	"strings"
)

func (s slackwatch) processCommand(m Message) bool {
	if !strings.HasPrefix(m.Text, "!") {
		return false
	}

	switch m.Text {
	case "!arm":
		*s.armed = true
		s.sendStatus(m.ChannelID)
	case "!disarm":
		*s.armed = false
		s.sendStatus(m.ChannelID)
	case "!verbose":
		*s.outputAll = true
		s.sendStatus(m.ChannelID)
	case "!quiet":
		*s.outputAll = false
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

func (s slackwatch) sendStatus(channelID string) {
	status := fmt.Sprintf("armed: %v, outputAll: %v", *s.armed, *s.outputAll)
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(status, channelID))
}

func (s slackwatch) sendHelp(channelID string) {
	text := "!arm !disarm !verbose !quiet !status !help"
	s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channelID))
}
