package chat

const (
	messageJoined       = "*joined the chat*"
	messageDisconnected = "*disconnected from the chat*"
)

type Message struct {
	Owner, Text   string
	UnixTimestamp int64
}

func BuildMessage(name, text string, unixTimestamp int64) Message {
	return Message{Owner: name, Text: text, UnixTimestamp: unixTimestamp}
}

func BuildJoinMessage(name string, unixTimestamp int64) Message {
	return BuildMessage(name, messageJoined, unixTimestamp)
}

func BuildDisconnectMessage(name string, unixTimestamp int64) Message {
	return BuildMessage(name, messageDisconnected, unixTimestamp)
}
