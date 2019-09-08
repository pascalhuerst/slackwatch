package slackwatch

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Message represents an incoming message
type Message struct {
	Timestamp string
	ChannelID string // of the form C02G9FMGB
	Channel   string // without preceding #
	UserID    string // of the form U3L6DT8P3
	User      string
	Text      string // raw text of the message, you probably want DetokenizedText()
	sw        *Slackwatch
}

func newMessage(timestamp string, channel string, user string, text string, sw *Slackwatch) Message {
	m := Message{
		Timestamp: ts(timestamp),
		ChannelID: channel,
		Channel:   sw.getCachedConversation(channel),
		UserID:    user,
		User:      sw.getCachedUser(user),
		Text:      text,
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

// String returns a string representation of the detokenized text preceded by
// the channel and username, suitable for logging or showing to the user.
// * [general] <user2> hello @user1
// The initial asterisk or space indicates if the message is considered
// interesting, see IsInteresting.
func (m Message) String() string {
	interesting := " "
	if m.IsInteresting() {
		interesting = "*"
	}

	return fmt.Sprintf("%s [%s] <%s> %s",
		interesting,
		m.Channel,
		m.User,
		m.DetokenizedText())
}

// IsFromMe returns true if the user our API token belongs to sent the message.
func (m Message) IsFromMe() bool {
	if m.UserID == m.sw.me.ID {
		return true
	}
	return false
}

// IsInteresting returns true if we were mentioned or this is a DM.
func (m Message) IsInteresting() bool {
	if m.Channel == "DM" {
		return true
	}

	if m.IsWatchedChannel() {
		return true
	}

	if strings.HasPrefix(m.Channel, "mpdm-") {
		return true
	}

	if strings.Contains(m.Text, m.sw.me.ID) {
		return true
	}

	return false
}

// IsWatchedChannel returns true if the channel is on the watch list.
func (m Message) IsWatchedChannel() bool {
	for _, c := range m.sw.watchedChan {
		if c == m.Channel {
			return true
		}
	}
	return false
}

// DetokenizedText replaces user mentions with their name instead of their ID.
func (m Message) DetokenizedText() string {
	t := m.Text
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
