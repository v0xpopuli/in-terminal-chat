package server

import (
	"in-terminal-chat/internal/chat"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type (
	Server interface {
		Run()
	}

	server struct {
		address string
		router  *mux.Router
		hub     chat.Hub
	}
)

func New(address string, hub chat.Hub) Server {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}

	router := mux.NewRouter()
	router.HandleFunc("/start", start(hub, upgrader))

	return server{address: address, hub: hub, router: router}
}

func (s server) Run() {
	logrus.Fatal(http.ListenAndServe(s.address, s.router))
}

func start(h chat.Hub, upgrader websocket.Upgrader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.WithError(err).Error("Failed to upgrade connection to the WebSocket protocol")
			return
		}

		name := r.URL.Query().Get("name")

		notifyExit := make(chan struct{})
		broadcaster := h.GetBroadcastingChannel()

		c := chat.NewClient(name, conn, notifyExit, broadcaster)
		h.Add(&c)

		go c.Listen()
		go c.Publish()
		<-notifyExit

		h.Remove(&c)
	}
}
