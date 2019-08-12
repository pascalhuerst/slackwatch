package slackwatch

func (s slackwatch) alert(m message) {
  if !*s.armed {
    return
  }

  go s.playAudio()
}
