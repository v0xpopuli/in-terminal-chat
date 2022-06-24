package chat

import "github.com/sirupsen/logrus"

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
			logrus.Debugf("Connection added: %+v", *c)

		case c := <-h.remove:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close((*c).Buffer())
				logrus.Debugf("Connection removed: %+v", *c)
			}

		case message := <-h.broadcast:
			for c := range h.clients {
				select {
				case (*c).Buffer() <- message:
				default:
					close((*c).Buffer())
					delete(h.clients, c)
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
	h.broadcast <- BuildJoinMessage(name)
}

func (h hub) NotifyDisconnect(name string) {
	h.broadcast <- BuildDisconnectMessage(name)
}

func (h hub) GetBroadcastingChannel() chan Message {
	return h.broadcast
}
