package slackwatch

import (
  "fmt"
  "os/exec"
  "net/http"
  "strings"
  "io/ioutil"
)

type Action struct {
  Command string
  Args string
  URL string
  Body string
}

func (a Action) Execute() {
  if a.Command != "" {
    a.RunCommand()
  }
  if a.URL != "" {
    if a.Body != "" {
      a.PostURL()
    } else {
      a.GetURL()
    }
  }
}

func (a Action) RunCommand() {
  cmd := exec.Command(a.Command, a.Args)
  err := cmd.Run()
  if err != nil {
    fmt.Println("Error running", a.Command, err)
  }
}

func (a Action) PostURL() {
  res, err := http.Post(a.URL, "application/octet-stream", strings.NewReader(a.Body))
  if err != nil {
    fmt.Println("Error POSTing", a.URL, err)
  }
  body, err := ioutil.ReadAll(res.Body)
  fmt.Println(body)
  res.Body.Close()
}

func (a Action) GetURL() {
  res, err := http.Get(a.URL)
  if err != nil {
    fmt.Println("Error GETing", a.URL, err)
  }
  res.Body.Close()
}
