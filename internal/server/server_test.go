package server

import (
	"in-terminal-chat/internal/chat"
	"strconv"
	"testing"
	"time"

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
	logrus.SetLevel(logrus.DebugLevel)

	hub := chat.NewHub()
	go hub.Run()

	availablePort, err := freeport.GetFreePort()
	s.NoError(err)

	s.address = "localhost:" + strconv.Itoa(availablePort)
	s.fullULR = "ws://" + s.address + "/start?name="

	go New(s.address, hub).Run()
}

func (s *ServerTestSuite) TestServer() {
	egor, wenjie := "Egor", "Wenjie"

	connOne, _, err := websocket.DefaultDialer.Dial(s.fullULR+egor, nil)
	s.NoError(err)

	time.Sleep(2 * time.Second)
	_, _, err = websocket.DefaultDialer.Dial(s.fullULR+wenjie, nil)
	s.NoError(err)

	actualMessages := make([]chat.Message, 0)
	go func() {
		for {
			var m chat.Message
			if err := connOne.ReadJSON(&m); err != nil {
				break
			}
			actualMessages = append(actualMessages, m)
		}
	}()
	time.Sleep(2 * time.Second)

	s.Len(actualMessages, 2)
}
