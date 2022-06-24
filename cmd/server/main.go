package main

import (
	"flag"
	"in-terminal-chat/internal/chat"
	"in-terminal-chat/internal/server"

	"github.com/sirupsen/logrus"
)

func main() {
	address := flag.String("address", "localhost:8080", "http service address")
	debug := flag.Bool("debug", true, "debug mode")
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	hub := chat.NewHub()
	go hub.Run()

	logrus.Info("Server up and running...")
	server.New(*address, hub).Run()
}
