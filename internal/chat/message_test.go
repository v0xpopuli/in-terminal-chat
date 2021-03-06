package chat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type MessageTestSuite struct {
	suite.Suite

	name          string
	unixTimestamp int64
}

func (s *MessageTestSuite) SetupSuite() {
	s.name = "Sherlock Holmes"
	s.unixTimestamp = time.Date(2022, time.June, 25, 11, 50, 0, 0, time.UTC).Unix()
}

func TestMessageTestSuite(t *testing.T) {
	suite.Run(t, new(MessageTestSuite))
}

func (s *MessageTestSuite) TestBuildMessage() {
	text := "You know my methods, Watson."

	expected := Message{Owner: s.name, Text: text, UnixTimestamp: s.unixTimestamp}
	actual := BuildMessage(s.name, text, s.unixTimestamp)

	s.Equal(expected, actual)
}

func (s *MessageTestSuite) TestBuildJoinMessage() {
	expected := Message{Owner: s.name, Text: joinChatMessage, UnixTimestamp: s.unixTimestamp}
	actual := BuildJoinMessage(s.name, s.unixTimestamp)

	s.Equal(expected, actual)
}

func (s *MessageTestSuite) TestBuildDisconnectMessage() {
	expected := Message{Owner: s.name, Text: disconnectFromChatMessage, UnixTimestamp: s.unixTimestamp}
	actual := BuildDisconnectMessage(s.name, s.unixTimestamp)

	s.Equal(expected, actual)
}

func (s *MessageTestSuite) TestBuildNameExistsMessage() {
	expected := Message{Owner: s.name, Text: nameExistsMessage, UnixTimestamp: s.unixTimestamp}
	actual := BuildNameExistsMessage(s.name, s.unixTimestamp)

	s.Equal(expected, actual)
}
