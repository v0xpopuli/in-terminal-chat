package chat

const (
	joinChatMessage           = "*join the chat*"
	disconnectFromChatMessage = "*disconnect from the chat*"
)

type Message struct {
	Owner, Text   string
	UnixTimestamp int64
}

func BuildMessage(name, text string, unixTimestamp int64) Message {
	return Message{Owner: name, Text: text, UnixTimestamp: unixTimestamp}
}

func BuildJoinMessage(name string, unixTimestamp int64) Message {
	return BuildMessage(name, joinChatMessage, unixTimestamp)
}

func BuildDisconnectMessage(name string, unixTimestamp int64) Message {
	return BuildMessage(name, disconnectFromChatMessage, unixTimestamp)
}
