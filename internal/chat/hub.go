package chat

type (
	Hub interface {
		Run()

		Add(c *Client)
		Remove(c *Client)
		GetBroadcastingChannel() chan []byte
	}

	hub struct {
		add       chan *Client
		remove    chan *Client
		clients   map[*Client]bool
		broadcast chan []byte
	}
)

func NewHub() Hub {
	return hub{
		add:       make(chan *Client),
		remove:    make(chan *Client),
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

func (h hub) Run() {
	for {
		select {
		case c := <-h.add:
			h.clients[c] = true

		case c := <-h.remove:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close((*c).Buffer())
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

func (h hub) GetBroadcastingChannel() chan []byte {
	return h.broadcast
}
