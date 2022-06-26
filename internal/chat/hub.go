package chat

import (
	"time"

	"github.com/sirupsen/logrus"
)

type (
	Hub interface {
		Run()
		Add(c *Client)
		Remove(c *Client)
		GetBroadcastingChannel() chan Message
		NotifyJoin(name string)
		NotifyDisconnect(name string)
	}

	hub struct {
		add       chan *Client
		remove    chan *Client
		clients   map[*Client]bool
		broadcast chan Message
	}
)

func NewHub() Hub {
	return hub{
		add:       make(chan *Client),
		remove:    make(chan *Client),
		clients:   make(map[*Client]bool),
		broadcast: make(chan Message),
	}
}

func (h hub) Run() {
	for {
		select {
		case c := <-h.add:
			h.clients[c] = true
			logrus.WithField("client", *c).Debug("Connection added")

		case c := <-h.remove:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close((*c).Buffer())
				logrus.WithField("client", *c).Debug("Connection removed")
			}

		case message := <-h.broadcast:
			logrus.WithField("message", message).Debug("Message received")
			for c := range h.clients {
				select {
				case (*c).Buffer() <- message:
					logrus.WithField("client", *c).Debug("Message sent to client")
				default:
				}
			}
		}
	}
}

func (h hub) Add(c *Client) {
	h.add <- c
}

func (h hub) Remove(c *Client) {
	h.remove <- c
}

func (h hub) NotifyJoin(name string) {
	h.broadcast <- BuildJoinMessage(name, time.Now().Unix())
}

func (h hub) NotifyDisconnect(name string) {
	h.broadcast <- BuildDisconnectMessage(name, time.Now().Unix())
}

func (h hub) GetBroadcastingChannel() chan Message {
	return h.broadcast
}
