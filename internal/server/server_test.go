package server

// import (
// 	"in-terminal-chat/internal/chat"
// 	"strconv"
// 	"testing"
// 	"time"
//
// 	"github.com/gorilla/websocket"
// 	"github.com/phayes/freeport"
// 	"github.com/stretchr/testify/suite"
// )
//
// type ServerTestSuite struct {
// 	suite.Suite
//
// 	address string
// 	fullULR string
//
// 	ticker *time.Ticker
// }
//
// func TestServerTestSuite(t *testing.T) {
// 	suite.Run(t, new(ServerTestSuite))
// }
//
// func (s *ServerTestSuite) SetupSuite() {
// 	hub := chat.NewHub()
// 	go hub.Run()
//
// 	availablePort, err := freeport.GetFreePort()
// 	s.NoError(err)
//
// 	s.address = "localhost:" + strconv.Itoa(availablePort)
// 	s.fullULR = "ws://" + s.address + "/start?name="
//
// 	s.ticker = time.NewTicker(3 * time.Second)
//
// 	go New(s.address, hub).Run()
// }
//
// func (s *ServerTestSuite) TearDownSuite() {
// 	s.ticker.Stop()
// }
//
// func (s *ServerTestSuite) TestFlow() {
// 	egor, wenjie := "Egor", "Wenjie"
// 	hiThere, byeBye := "Hi there!", "Bye bye!"
//
// 	expectedMessages := []chat.Message{
// 		{Owner: egor, Text: "*join the chat*"},
// 		{Owner: wenjie, Text: "*join the chat*"},
// 		{Owner: egor, Text: hiThere},
// 		{Owner: wenjie, Text: byeBye},
// 		{Owner: wenjie, Text: "*disconnect from the chat*"},
// 	}
//
// 	connOne, _, err := websocket.DefaultDialer.Dial(s.fullULR+egor, nil)
// 	s.NoError(err)
//
// 	connTwo, _, err := websocket.DefaultDialer.Dial(s.fullULR+wenjie, nil)
// 	s.NoError(err)
//
// 	actualMessages := make([]chat.Message, 0)
// 	go func() {
// 		for {
// 			var m chat.Message
// 			if err := connOne.ReadJSON(&m); err != nil {
// 				break
// 			}
// 			actualMessages = append(actualMessages, m)
// 		}
// 	}()
//
// 	s.NoError(connOne.WriteJSON(chat.Message{Owner: egor, Text: hiThere}))
//
// 	<-s.ticker.C
// 	s.NoError(connTwo.WriteJSON(chat.Message{Owner: wenjie, Text: byeBye}))
// 	s.NoError(connTwo.WriteMessage(websocket.CloseMessage, []byte{}))
//
// 	<-s.ticker.C
// 	s.Len(actualMessages, 5)
// 	s.Equal(expectedMessages[0].Owner, actualMessages[0].Owner)
// 	s.Equal(expectedMessages[0].Text, actualMessages[0].Text)
// 	s.Equal(expectedMessages[1].Owner, actualMessages[1].Owner)
// 	s.Equal(expectedMessages[1].Text, actualMessages[1].Text)
// 	s.Equal(expectedMessages[2].Owner, actualMessages[2].Owner)
// 	s.Equal(expectedMessages[2].Text, actualMessages[2].Text)
// 	s.Equal(expectedMessages[3].Owner, actualMessages[3].Owner)
// 	s.Equal(expectedMessages[3].Text, actualMessages[3].Text)
// 	s.Equal(expectedMessages[4].Owner, actualMessages[4].Owner)
// 	s.Equal(expectedMessages[4].Text, actualMessages[4].Text)
// }
