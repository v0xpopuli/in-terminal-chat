package server

import (
	"in-terminal-chat/internal/chat"
	"strconv"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/phayes/freeport"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ServerTestSuite struct {
	suite.Suite

	address string
	fullULR string
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	logrus.Info("Setup server suite!")
	hub := chat.NewHub()
	go hub.Run()

	availablePort, err := freeport.GetFreePort()
	s.NoError(err)

	s.address = "localhost:" + strconv.Itoa(availablePort)
	s.fullULR = "ws://" + s.address + "/start?name="

	go New(s.address, hub).Run()
}

func (s *ServerTestSuite) TestServer() {
	logrus.Info("Start server test!!!")
	logrus.SetLevel(logrus.DebugLevel)

	egor, wenjie := "Egor", "Wenjie"
	joinTheChatMessage := "*join the chat*"

	connOne, _, err := websocket.DefaultDialer.Dial(s.fullULR+egor, nil)
	s.NoError(err)

	_, _, err = websocket.DefaultDialer.Dial(s.fullULR+wenjie, nil)
	s.NoError(err)

	release := make(chan struct{})
	actualMessages := make([]chat.Message, 0)
	go func() {
		for {
			var m chat.Message
			if err := connOne.ReadJSON(&m); err != nil {
				release <- struct{}{}
				break
			}
			actualMessages = append(actualMessages, m)
			if len(actualMessages) == 2 {
				release <- struct{}{}
			}
		}
	}()
	logrus.Info("Stuck before release...")
	<-release

	s.Equal(egor, actualMessages[0].Owner)
	s.Equal(joinTheChatMessage, actualMessages[0].Text)
	s.NotEmpty(actualMessages[0].UnixTimestamp)

	s.Equal(wenjie, actualMessages[1].Owner)
	s.Equal(joinTheChatMessage, actualMessages[1].Text)
	s.NotEmpty(actualMessages[1].UnixTimestamp)

}
