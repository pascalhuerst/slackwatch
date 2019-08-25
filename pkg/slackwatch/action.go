package slackwatch

import (
	"net/http"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

// An Action 's execute method is called when an interesting message is
// received if armed.
type Action interface {
	Execute(Message)
}

// URLAction specifies an HTTP request to make on Alert
type URLAction struct {
	URL  string
	Body string
}

// Execute is called to make the HTTP request
func (u URLAction) Execute(m Message) {
	var res *http.Response
	var err error
	if u.Body != "" {
		res, err = http.Post(u.URL, "application/octet-stream", strings.NewReader(u.Body))
	} else {
		res, err = http.Get(u.URL)
	}
	res.Body.Close()
	if err != nil {
		log.Errorf("Error requesting %s: %v", u.URL, err)
	}
}

// CommandAction specifies a command to execute on Alert
type CommandAction struct {
	Command string
	Args    string
}

// Execute is called to run the command.
func (c CommandAction) Execute(m Message) {
	cmd := exec.Command(c.Command, c.Args)
	err := cmd.Run()
	if err != nil {
		log.Errorf("Error running %s: %v", c.Command, err)
	}
}
