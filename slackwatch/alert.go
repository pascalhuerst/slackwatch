package slackwatch

func (s slackwatch) alert(m Message) {
	if !*s.armed {
		return
	}

	for _, action := range s.config.Actions {
		go action.Execute(m)
	}
}
