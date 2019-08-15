package slackwatch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type message struct {
	timestamp string
	channelID string
	channel   string
	userID    string
	user      string
	text      string
	sw        *slackwatch
}

func newMessage(timestamp string, channel string, user string, text string, sw *slackwatch) message {
	m := message{
		timestamp: ts(timestamp),
		channelID: channel,
		channel:   sw.getCachedConversation(channel),
		userID:    user,
		user:      sw.getCachedUser(user),
		text:      text,
		sw:        sw,
	}
	return m
}

func ts(ts string) string {
	i, err := strconv.ParseFloat(ts, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(int64(i), 0)
	return tm.Format("15:04:05")
}

func (m message) asString() string {
	interesting := " "
	if m.isInteresting() {
		interesting = "*"
	}

	return fmt.Sprintf("%s %s [%s] <%s> %s",
		interesting,
		m.timestamp,
		m.channel,
		m.user,
		m.detokenizedText())
}

func (m message) isFromMe() bool {
	if m.userID == m.sw.me.ID {
		return true
	}
	return false
}

func (m message) isInteresting() bool {
	if m.channel == "DM" {
		return true
	}
	if strings.HasPrefix(m.channel, "mpdm-") {
		return true
	}

	if strings.Contains(m.text, m.sw.me.ID) {
		return true
	}

	return false
}

func (m message) detokenizedText() string {
	t := m.text
	re := regexp.MustCompile("<@(\\w+)>")
	matches := re.FindAllStringSubmatch(t, -1)

	if matches == nil {
		return t
	}

	for _, match := range matches {
		t = strings.Replace(t, "<@"+match[1]+">", "@"+m.sw.getCachedUser(match[1]), -1)
	}
	return t
}
