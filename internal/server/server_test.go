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

	hub chat.Hub

	address string
	fullULR string
}

func TestServerTestSuite(t *testing.T) {
	suite.Run(t, new(ServerTestSuite))
}

func (s *ServerTestSuite) SetupSuite() {
	logrus.SetLevel(logrus.DebugLevel)

	s.hub = chat.NewHub()
	go s.hub.Run()

	availablePort, err := freeport.GetFreePort()
	s.NoError(err)

	s.address = "localhost:" + strconv.Itoa(availablePort)
	s.fullULR = "ws://" + s.address + "/start?name="

	go New(s.address, s.hub).Run()
}

func (s *ServerTestSuite) TestServer() {
	var m chat.Message
	egor, wenjie := "Egor", "Wenjie"

	conn, _, err := websocket.DefaultDialer.Dial(s.fullULR+egor, nil)
	s.NoError(err)

	time.Sleep(2 * time.Second)
	conn.ReadJSON(&m)
	s.Equal(egor, m.Owner)

	_, _, err = websocket.DefaultDialer.Dial(s.fullULR+wenjie, nil)
	s.NoError(err)

	time.Sleep(2 * time.Second)
	conn.ReadJSON(&m)
	s.Equal(wenjie, m.Owner)
}
