package chat

import (
	"bytes"
	"log"

	"github.com/gorilla/websocket"
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

type (
	Data struct {
		Name    string `json:"name"`
		Message string `json:"message"`
	}

	Client interface {
		Publish()
		Listen()

		Buffer() chan []byte
	}

	client struct {
		name                string
		conn                *websocket.Conn
		buffer, broadcaster chan []byte
		notifyExit          chan struct{}
	}
)

func NewClient(name string, conn *websocket.Conn, notifyExit chan struct{}, broadcaster chan []byte) Client {
	return client{name: name, conn: conn, notifyExit: notifyExit, broadcaster: broadcaster, buffer: make(chan []byte)}
}

func (c client) Publish() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.broadcaster <- c.prependName(message)
	}
	c.conn.Close()
	c.notifyExit <- struct{}{}
}

func (c client) Listen() {
	for {
		select {
		case message, ok := <-c.buffer:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			if err := w.Close(); err != nil {
				return
			}
		}
	}
	c.conn.Close()
}

func (c client) Buffer() chan []byte {
	return c.buffer
}

func (c client) prependName(message []byte) []byte {
	return []byte(c.name + ": " + string(message))
}
