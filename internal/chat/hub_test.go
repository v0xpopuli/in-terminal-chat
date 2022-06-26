package chat

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type HubTestSuite struct {
	suite.Suite

	ticker *time.Ticker

	hub hub
}

func TestHubSuiteTestSuite(t *testing.T) {
	suite.Run(t, new(HubTestSuite))
}

func (s *HubTestSuite) SetupSuite() {
	s.hub = NewHub().(hub)
	go s.hub.Run()
}

func (s *HubTestSuite) SetupTest() {
	s.ticker = time.NewTicker(3 * time.Second)
}

func (s *HubTestSuite) TearDownSuite() {
	s.ticker.Stop()
}

func (s *HubTestSuite) TestAdd() {
	c := NewClient("Tony Stark", nil, nil, nil)
	s.hub.Add(&c)

	<-s.ticker.C
	s.True(s.hub.clients[&c])
}

func (s *HubTestSuite) TestRemove() {
	c := NewClient("Thor Odinson", nil, nil, nil)
	s.hub.clients[&c] = true

	s.True(s.hub.clients[&c])

	s.hub.Remove(&c)

	<-s.ticker.C
	_, actual := s.hub.clients[&c]
	s.False(actual)
}

func (s *HubTestSuite) TestBroadcast() {
	expected := Message{Owner: "Captain America", Text: "Avengers assemble"}

	c := NewClient("Hawkeye", nil, nil, nil)
	s.hub.Add(&c)

	go func() {
		s.Equal(expected, <-c.Buffer())
	}()

	<-s.ticker.C
	message := BuildMessage("Captain America", "Avengers assemble", 0)
	s.hub.GetBroadcastingChannel() <- message
}

func (s *HubTestSuite) TestNotifyJoin() {
	expected := Message{Owner: "Captain America", Text: "*join the chat*"}

	c := NewClient("Hulk", nil, nil, nil)
	s.hub.Add(&c)

	go func() {
		actual := <-c.Buffer()
		s.Equal(expected.Owner, actual.Owner)
		s.Equal(expected.Text, actual.Text)
	}()

	<-s.ticker.C
	s.hub.NotifyJoin("Captain America")
}

func (s *HubTestSuite) TestNotifyDisconnect() {
	expected := Message{Owner: "Captain America", Text: "*disconnect from the chat*"}

	c := NewClient("Hulk", nil, nil, nil)
	s.hub.Add(&c)

	go func() {
		actual := <-c.Buffer()
		s.Equal(expected.Owner, actual.Owner)
		s.Equal(expected.Text, actual.Text)
	}()

	<-s.ticker.C
	s.hub.NotifyDisconnect("Captain America")
}

func (s *HubTestSuite) TestGetBroadcastingChannel() {
	s.NotNil(s.hub.GetBroadcastingChannel())
}
