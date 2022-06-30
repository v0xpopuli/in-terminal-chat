package chat

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type (
	Client interface {
		Publish()
		Listen()

		Name() string
		Buffer() chan Message
		CloseConnectionWithMessage(Message)
	}

	client struct {
		name                string
		conn                Connection
		buffer, broadcaster chan Message
		notifyExit          chan struct{}
	}
)

func NewClient(name string, conn Connection, notifyExit chan struct{}, broadcaster chan Message) Client {
	return client{
		name:        name,
		conn:        conn,
		notifyExit:  notifyExit,
		broadcaster: broadcaster,
		buffer:      make(chan Message),
	}
}

func (c client) Publish() {
	for {
		message, err := c.conn.ReadJSONMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.WithField("name", c.name).WithError(err).Error("Websocket unexpectedly closed")
			}
			break
		}
		c.broadcaster <- message
	}
	c.conn.Close()
	c.notifyExit <- struct{}{}
}

func (c client) Listen() {
	for {
		select {
		case message, ok := <-c.buffer:
			if !ok {
				c.conn.WriteCloseMessage()
				logrus.WithField("name", c.name).Warn("Attempt to read from closed channel")
				return
			}

			if err := c.conn.WriteJSONMessage(message); err != nil {
				logrus.WithField("name", c.name).WithError(err).Error("Error occurred during attempt to send message")
				return
			}
		}
	}
	c.conn.Close()
}

func (c client) Name() string {
	return c.name
}

func (c client) Buffer() chan Message {
	return c.buffer
}

func (c client) CloseConnectionWithMessage(m Message) {
	c.conn.WriteJSONMessage(m)
	c.conn.WriteCloseMessage()
}
