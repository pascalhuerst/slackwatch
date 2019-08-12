package slackwatch

import (
  "fmt"
  "strings"

)

func (s slackwatch) processCommand(m message) bool {
  if !strings.HasPrefix(m.text, "!") {
    return false
  }

  switch m.text {
  case "!arm":
    *s.armed = true
    s.sendStatus(m.channelID)
  case "!disarm":
    *s.armed = false
    s.sendStatus(m.channelID)
  case "!verbose":
    *s.outputAll = true
    s.sendStatus(m.channelID)
  case "!quiet":
    *s.outputAll = false
    s.sendStatus(m.channelID)
  case "!status":
    s.sendStatus(m.channelID)
  case "!help":
    s.sendHelp(m.channelID)
  default:
    s.rtm.SendMessage(s.rtm.NewOutgoingMessage("Invalid Command", m.channelID))
  }

  return true
}

func (s slackwatch) sendStatus(channelID string) {
  status := fmt.Sprintf("armed: %v, outputAll: %v", *s.armed, *s.outputAll )
  s.rtm.SendMessage(s.rtm.NewOutgoingMessage(status, channelID))
}

func (s slackwatch) sendHelp(channelID string) {
  text := "!arm !disarm !verbose !quiet !status !help"
  s.rtm.SendMessage(s.rtm.NewOutgoingMessage(text, channelID))
}
