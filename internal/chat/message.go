package chat

import (
	"time"
)

const (
	joinChatMessage           = "*join the chat*"
	disconnectFromChatMessage = "*disconnect from the chat*"
)

type Message struct {
	Owner, Text   string
	UnixTimestamp int64
}

func BuildMessage(name, text string) Message {
	return Message{Owner: name, Text: text, UnixTimestamp: time.Now().Unix()}
}

func BuildJoinMessage(name string) Message {
	return BuildMessage(name, joinChatMessage)
}

func BuildDisconnectMessage(name string) Message {
	return BuildMessage(name, disconnectFromChatMessage)
}
