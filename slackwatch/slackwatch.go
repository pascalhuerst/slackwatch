// Package slackwatch preforms configfured actions when DMed on Slack.
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
	"log"

	"github.com/nlopes/slack"
)

type slackwatch struct {
	userLookup         map[string]string
	conversationLookup map[string]string
	me                 *slack.UserDetails
	api                *slack.Client
	rtm                *slack.RTM
	armed              *bool
	outputAll          *bool
	config             *Config
}

// New creates a slackwatch instance
func New(config Config) slackwatch {
	// https://stackoverflow.com/questions/28817992/how-to-set-bool-pointer-to-true-in-struct-literal
	s := slackwatch{
		api:                slack.New(config.SlackToken),
		userLookup:         make(map[string]string),
		conversationLookup: make(map[string]string),
		armed:              &[]bool{true}[0],
		outputAll:          &[]bool{false}[0],
		config:             &config,
	}
	return s
}

func (s slackwatch) Run() {
	s.rtm = s.api.NewRTM()
	go s.rtm.ManageConnection()

	for msg := range s.rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.ConnectedEvent:
			log.Print("Connected")
			s.me = ev.Info.User

		case *slack.DisconnectedEvent:
			log.Print("Disconnected")

		case *slack.MessageEvent:
			s.messageReceived(ev)

		case *slack.ChannelJoinedEvent:
			name := ev.Channel.Name
			if name == "" {
				name = "DM"
			}
			s.alert(Message{Channel: ev.Channel.Name, Text: "Channel Joined"})
			log.Print("* Joined to new channel", name)

		case *slack.IncomingEventError:
			log.Printf("Incoming Event Error: %v", ev)

		case *slack.ConnectionErrorEvent:
			log.Printf("Connection Error: %v", ev)

		case *slack.RTMError:
			log.Printf("Error: %s\n", ev.Error())

		case *slack.InvalidAuthEvent:
			log.Printf("Invalid credentials")
			return

		default:
		}
	}
}

func (s slackwatch) messageReceived(msg *slack.MessageEvent) {
	if msg.Text == "" {
		return
	}

	m := newMessage(msg.Timestamp, msg.Channel, msg.User, msg.Text, &s)

	if m.IsFromMe() && m.Channel == "DM" {
		if s.processCommand(m) {
			return
		}
	}

	if m.IsInteresting() {
		log.Print(m.String())
		s.alert(m)
	} else {
		if *s.outputAll {
			string := m.String()
			if len(string) > 60 {
				string = string[:60]
			}
			log.Print(string)
		}
	}
}
