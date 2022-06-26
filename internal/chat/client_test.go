package chat

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ClientTestSuite struct {
	suite.Suite
	controller *gomock.Controller

	name        string
	conn        *MockConnection
	notifyExit  chan struct{}
	broadcaster chan Message
	client      Client
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}

func (s *ClientTestSuite) SetupSuite() {
	s.controller = gomock.NewController(s.T())
	s.name = "Ukupnik"
}

func (s *ClientTestSuite) SetupTest() {
	s.conn = NewMockConnection(s.controller)
	s.notifyExit = make(chan struct{})
	s.broadcaster = make(chan Message)
	s.client = NewClient(s.name, s.conn, s.notifyExit, s.broadcaster)
}

func (s *ClientTestSuite) TearDownSuite() {
	s.controller.Finish()
}

func (s *ClientTestSuite) TestPublish() {
	expected := Message{Owner: "Kirkorov", Text: "Hi!", UnixTimestamp: 0}

	s.conn.EXPECT().
		ReadJSONMessage().
		DoAndReturn(func() (Message, error) {
			return BuildMessage("Kirkorov", "Hi!", 0), nil
		}).
		AnyTimes()

	go s.client.Publish()
	s.Equal(expected, <-s.broadcaster)
}

func (s *ClientTestSuite) TestPublish_Error() {
	expected := struct{}{}

	s.conn.EXPECT().
		ReadJSONMessage().
		DoAndReturn(func() (Message, error) {
			return Message{}, errors.New("something went wrong")
		}).
		AnyTimes()

	s.conn.EXPECT().
		Close().
		Return(nil)

	go s.client.Publish()
	s.Equal(expected, <-s.notifyExit)
}

func (s *ClientTestSuite) TestListen() {
	message := BuildMessage("Galkin", "Hello!", 0)
	s.conn.EXPECT().
		WriteJSONMessage(message).
		Return(nil).
		AnyTimes()

	s.conn.EXPECT().
		WriteCloseMessage().
		Return(nil).
		AnyTimes()

	s.conn.EXPECT().
		Close().
		Return(nil).
		AnyTimes()

	go s.client.Listen()

	s.client.Buffer() <- message
	close(s.client.Buffer())
}

func (s *ClientTestSuite) TestListen_Error() {
	message := BuildMessage("Galkin", "Hello!", 0)
	s.conn.EXPECT().
		WriteJSONMessage(message).
		Return(errors.New("something went wrong")).
		AnyTimes()

	s.conn.EXPECT().
		Close().
		Return(nil).
		AnyTimes()

	go s.client.Listen()

	s.client.Buffer() <- message
}
