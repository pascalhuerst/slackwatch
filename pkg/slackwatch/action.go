package slackwatch

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

// An Action 's execute method is called when an interesting message is
// received if armed.
type Action interface {
	Execute(Message)
}

// URLAction specifies an HTTP request to make on Alert. If Body is provided,
// an HTTP post is made, otherwise, an HTTP get.
type URLAction struct {
	URL         string
	Body        string
	ContentType string
	AccessToken string
}

// Execute is called to make the HTTP request
func (u URLAction) Execute(m Message) {

	client := &http.Client{}
	req, err := http.NewRequest("POST", u.URL, strings.NewReader(u.Body))
	if err != nil {
		logrus.Error("Error requesting %s: %v", u.URL, err)
		return
	}

	req.Header.Add("Content-Type", u.ContentType)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", u.AccessToken))

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Error requesting %s: %v", u.URL, err)
	}
	defer resp.Body.Close()

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
		logrus.Errorf("Error running %s: %v", c.Command, err)
	}
}
