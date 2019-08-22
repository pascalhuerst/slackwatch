package slackwatch

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

// An Action 's execute method is called when an interesting message is received.
type Action interface {
	Execute(Message)
}

// DefaultAction supports two types of common actions, Command, and URL. For
// URLs, if Body is provided an HTTP Post is performed, otherwise an HTTP Get.
type DefaultAction struct {
	Command string
	Args    string
	URL     string
	Body    string
}

// Execute runs the specified command or preforms an HTTP request.
func (a DefaultAction) Execute(m Message) {
	if a.Command != "" {
		a.runCommand()
	}
	if a.URL != "" {
		if a.Body != "" {
			a.postURL()
		} else {
			a.getURL()
		}
	}
}

func (a DefaultAction) runCommand() {
	cmd := exec.Command(a.Command, a.Args)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running", a.Command, err)
	}
}

func (a DefaultAction) postURL() {
	res, err := http.Post(a.URL, "application/octet-stream", strings.NewReader(a.Body))
	if err != nil {
		fmt.Println("Error POSTing", a.URL, err)
	}
	res.Body.Close()
}

func (a DefaultAction) getURL() {
	res, err := http.Get(a.URL)
	if err != nil {
		fmt.Println("Error GETing", a.URL, err)
	}
	res.Body.Close()
}
