package slackwatch

import (
  "fmt"


  // "github.com/davecgh/go-spew/spew"
)

func (s slackwatch) getCachedUser(id string) string {
  if id == "" {
    return ""
  }

  nick, ok := s.userLookup[id]
  if ok {
    return nick
  }

  fmt.Println("Looking up ", id)
  user, err := s.api.GetUserInfo(id)
  if err != nil {
    return ""
  }

  s.userLookup[id] = user.Name
  return user.Name
}

func (s slackwatch) getCachedConversation(id string) string {
  if id == "" {
    return ""
  }

  conv, ok := s.conversationLookup[id]
  if ok {
    return conv
  }

  fmt.Println("Looking up ", id)
  info, err := s.api.GetConversationInfo(id, false)
  if err != nil {
    fmt.Printf("Error Looking up Conversation: %v\n", err)
    return ""
  }

  name := info.Name
  if !info.IsChannel && name == "" {
    name = "DM"
  }

  s.conversationLookup[id] = name
  return name
}
