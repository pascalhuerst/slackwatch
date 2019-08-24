package slackwatch

func (s *Slackwatch) alert(m Message) {
	if !s.armed {
		return
	}

	for _, action := range s.config.Actions {
		go action.Execute(m)
	}
}
