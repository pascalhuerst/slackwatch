package slackwatch

func (s slackwatch) alert() {
	if !*s.armed {
		return
	}

	for _, action := range s.config.Actions {
		go action.execute()
	}
}
