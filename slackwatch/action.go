package slackwatch

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

type Action struct {
	Command string
	Args    string
	URL     string
	Body    string
}

func (a Action) execute() {
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

func (a Action) runCommand() {
	cmd := exec.Command(a.Command, a.Args)
	err := cmd.Run()
	if err != nil {
		fmt.Println("Error running", a.Command, err)
	}
}

func (a Action) postURL() {
	res, err := http.Post(a.URL, "application/octet-stream", strings.NewReader(a.Body))
	if err != nil {
		fmt.Println("Error POSTing", a.URL, err)
	}
	res.Body.Close()
}

func (a Action) getURL() {
	res, err := http.Get(a.URL)
	if err != nil {
		fmt.Println("Error GETing", a.URL, err)
	}
	res.Body.Close()
}
