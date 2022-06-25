package chat

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	expectedName := "Agent Smith"
	expectedConn := &websocket.Conn{}
	expectedNotifyExit := make(chan struct{})
	expectedBroadcaster := make(chan Message)

	actual := NewClient(expectedName, expectedConn, expectedNotifyExit, expectedBroadcaster).(client)

	assert.Equal(t, expectedName, actual.name)
	assert.Equal(t, expectedConn, actual.conn)
	assert.Equal(t, expectedNotifyExit, actual.notifyExit)
	assert.Equal(t, expectedBroadcaster, actual.broadcaster)
}
