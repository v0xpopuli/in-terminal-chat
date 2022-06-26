package chat

//go:generate mockgen --source=conn.go --destination=conn_mock.go --package=chat

import (
	"github.com/gorilla/websocket"
)

type (
	Connection interface {
		ReadJSONMessage() (Message, error)
		WriteJSONMessage(Message) error
		Close() error
		WriteCloseMessage() error
	}

	connection struct {
		conn *websocket.Conn
	}
)

func NewConnection(conn *websocket.Conn) Connection {
	return connection{conn: conn}
}

func (c connection) ReadJSONMessage() (m Message, _ error) {
	return m, c.conn.ReadJSON(&m)
}

func (c connection) WriteJSONMessage(m Message) error {
	return c.conn.WriteJSON(m)
}

func (c connection) WriteCloseMessage() error {
	return c.conn.WriteMessage(websocket.CloseMessage, []byte{})
}

func (c connection) Close() error {
	return c.conn.Close()
}
