package slackwatch

import (
	"fmt"
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

		case *slack.MessageEvent:
			s.messageReceived(ev)

		case *slack.ChannelJoinedEvent:
			name := ev.Channel.Name
			if name == "" {
				name = "DM"
			}
			s.alert()
			log.Print("* Joined to new channel", name)

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

	if m.isFromMe() && m.channel == "DM" {
		if s.processCommand(m) {
			return
		}
	}

	if m.isInteresting() {
		fmt.Println(m.asString())
		s.alert()
	} else {
		if *s.outputAll {
			fmt.Println(m.asString())
		}
	}
}
