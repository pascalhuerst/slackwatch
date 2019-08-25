// Package slackwatch preforms configured actions when DMed on Slack.
// Out of the box, it expects a JSON formatted config file named .slackwatch
// in your home directory.
//
//   {
//     "SlackToken": "xoxp-123-543",
//     "Actions": [
//       { "Command": "/usr/bin/afplay", "Args": "klaxon.wav" },
//       { "URL": "https://hassio.local/api/services/homeassistant/turn_on?api_password=letmein",
//         "Body": "{\"entity_id\":\"switch.bat_signal\"}"
//       }
//     ]
//   }
package slackwatch

import (
	"github.com/nlopes/slack"
	log "github.com/sirupsen/logrus"
)

// Slackwatch struct holds state. You should call New(config) rather than creating it yourself.
type Slackwatch struct {
	userLookup         map[string]string
	conversationLookup map[string]string
	me                 *slack.UserDetails
	api                *slack.Client
	rtm                *slack.RTM
	armed              bool
	outputAll          bool
	config             *Config
}

// New creates a slackwatch instance and returns a pointer to it.
func New(config Config) *Slackwatch {
	s := Slackwatch{
		api:                slack.New(config.SlackToken),
		userLookup:         make(map[string]string),
		conversationLookup: make(map[string]string),
		armed:              true,
		outputAll:          false,
		config:             &config,
	}
	s.rtm = s.api.NewRTM()
	return &s
}

// Run is a blocking call that makes the connection to Slack and handles incoming events.
func (s *Slackwatch) Run() {
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.Info("Connected")
			s.me = ev.Info.User

		case *slack.DisconnectedEvent:
			log.Error("Disconnected")

		case *slack.MessageEvent:
			if ev.Text != "" {
				m := newMessage(ev.Timestamp, ev.Channel, ev.User, ev.Text, s)
				s.messageReceived(m)
			}

		case *slack.ChannelJoinedEvent:
			name := ev.Channel.Name
			if name == "" {
				name = "DM"
			}
			s.alert(Message{Channel: ev.Channel.Name, Text: "Channel Joined"})
			log.Info("* Joined to new channel", name)

		case *slack.IncomingEventError:
			log.Errorf("Incoming Event Error: %v", ev)

		case *slack.ConnectionErrorEvent:
			log.Errorf("Connection Error: %v", ev)

		case *slack.RTMError:
			log.Errorf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Fatal("Invalid credentials")
			return

		// some types we don't care about
		case *slack.PresenceChangeEvent:
		case *slack.ManualPresenceChangeEvent:
		case *slack.LatencyReport:
		case *slack.HelloEvent:
		case *slack.ConnectingEvent:
		case *slack.UserTypingEvent:
		case *slack.EmojiChangedEvent:
		case *slack.ReactionAddedEvent:
		case *slack.ReactionRemovedEvent:
		case *slack.FilePublicEvent:
		case *slack.FileSharedEvent:
		case *slack.FileChangeEvent:
		case *slack.FileDeletedEvent:
		case *slack.FileUnsharedEvent:
		case *slack.UserChangeEvent:
		case *slack.MemberJoinedChannelEvent:
		case *slack.MemberLeftChannelEvent:
		case *slack.DNDUpdatedEvent:
		case *slack.GroupMarkedEvent:
		case *slack.ChannelMarkedEvent:
		case *slack.ChannelCreatedEvent:
		case *slack.ChannelArchiveEvent:
		case *slack.IMMarkedEvent:
		case *slack.PinAddedEvent:
		case *slack.PinRemovedEvent:
		case *slack.BotAddedEvent:
		case *slack.AckMessage:
		case *slack.PrefChangeEvent:

		default:
			log.Printf("Unknown Event '%T': %v", ev, ev)
		}
	}
}

func (s *Slackwatch) messageReceived(m Message) {
	if m.IsFromMe() && m.Channel == "DM" {
		if s.processCommand(m) {
			return
		}
	}

	if m.IsInteresting() {
		log.Print(m.String())
		s.alert(m)
	} else {
		if s.outputAll {
			string := m.String()
			if len(string) > 60 {
				string = string[:60]
			}
			log.Print(string)
		}
	}
}
