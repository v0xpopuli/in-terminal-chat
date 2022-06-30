package server

import (
	"in-terminal-chat/internal/chat"
	"net/http"
	"time"

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

func start(hub chat.Hub, upgrader websocket.Upgrader) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		wsConn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logrus.WithError(err).Error("Failed to upgrade connection to the WebSocket protocol")
			return
		}

		name := r.URL.Query().Get("name")

		notifyExit := make(chan struct{})
		broadcaster := hub.GetBroadcastingChannel()

		conn := chat.NewConnection(wsConn)
		client := chat.NewClient(name, conn, notifyExit, broadcaster)
		if hub.NameExists(name) {
			client.CloseConnectionWithMessage(chat.BuildNameExistsMessage(name, time.Now().Unix()))
			logrus.WithField("client", client).Debug("Name exists in current session, close connection")
			return
		}
		hub.Add(&client)

		go client.Listen()
		go client.Publish()

		hub.NotifyJoin(name)
		<-notifyExit
		hub.Remove(&client)

		hub.NotifyDisconnect(name)
	}
}
