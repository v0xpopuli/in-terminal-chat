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
		logrus.Error(color.InRed("Name field can't be empty"))
		os.Exit(1)
	}

	conn, _, err := websocket.DefaultDialer.Dial(makeDialURL(address, name), nil)
	if err != nil {
		logrus.WithError(err).Error("Failed to establish connection with server")
		os.Exit(1)
	}

	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	go readMessages(writer, conn, *name)

	writeMessage(writer, conn)
}

func makeDialURL(address *string, name *string) string {
	return *address + "/start?name=" + *name
}

func writeMessage(w *uilive.Writer, c *websocket.Conn) {
	fmt.Fprint(w, color.InGreen("You are successfully connected!\n\n"))
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}

		if err := c.WriteMessage(websocket.TextMessage, []byte(scanner.Text())); err != nil {
			logrus.WithError(err).Error("Failed to send message")
			break
		}
	}
}

func readMessages(w *uilive.Writer, c *websocket.Conn, name string) {
	for {
		_, m, err := c.ReadMessage()
		if err != nil {
			logrus.WithError(err).Error("Failed to read message")
		}
		if strings.Contains(string(m), name) {
			fmt.Fprintln(w, color.InCyan(string(m)))
			w.Flush()
		} else {
			fmt.Println(color.InYellow(string(m)))
		}
	}
}
