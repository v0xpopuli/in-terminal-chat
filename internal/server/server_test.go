package server

import (
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotEmpty(t, New("", nil))
}

func Test_start(t *testing.T) {
	assert.NotEmpty(t, start(nil, websocket.Upgrader{}))
}
