package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TwiN/go-color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uilive"
	"github.com/sirupsen/logrus"
)

func main() {
	address := flag.String("address", "ws://localhost:8080", "http service address")
	name := flag.String("name", "", "client name")
	flag.Parse()

	if strings.TrimSpace(*name) == "" {
		logrus.Error("Name field can't be empty")
		os.Exit(1)
	}

	url := *address + "/start?name=" + *name
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {

		logrus.WithError(err).Error("Failed to establish connection with server")
		os.Exit(1)
	}

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	go func() {
		for {
			_, m, err := conn.ReadMessage()
			if err != nil {
				logrus.WithError(err).Error("Failed to read message")
			}
			if strings.Contains(string(m), *name) {
				fmt.Fprintln(writer, color.InCyan(string(m)))
				writer.Flush()
			} else {
				fmt.Println(color.InYellow(string(m)))
			}
		}
	}()

	fmt.Fprint(writer, color.InGreen("You are successfully connected!\n\n"))
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}

		conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
		if err != nil {
			logrus.WithError(err).Error("Failed to send message")
			break
		}
	}
}
