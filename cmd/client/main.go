package main

import (
	"bufio"
	"flag"
	"fmt"
	"in-terminal-chat/internal/chat"
	"os"
	"strings"
	"time"

	"github.com/TwiN/go-color"
	"github.com/gorilla/websocket"
	"github.com/gosuri/uilive"
)

const messagePattern = "[%s] %s -> %s"

func main() {
	address := flag.String("address", "ws://localhost:8080", "http service address")
	name := flag.String("name", "", "client name")
	flag.Parse()

	if strings.TrimSpace(*name) == "" {
		fmt.Println(color.InRed("Name arg can't be empty"))
		os.Exit(1)
	}

	conn, _, err := websocket.DefaultDialer.Dial(makeDialURL(address, name), nil)
	if err != nil {
		fmt.Printf("Failed to establish connection with server: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf(color.InGreen("You are successfully connected!\n\n"))
	go readMessages(conn, *name)

	writeMessage(conn, *name)
}

func writeMessage(c *websocket.Conn, name string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		message := scanner.Text()
		if message == "" {
			continue
		}

		if err := c.WriteJSON(chat.BuildMessage(name, message)); err != nil {
			fmt.Printf("Failed to send message: %v\n", err)
			break
		}
	}
}

func readMessages(c *websocket.Conn, name string) {
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	for {
		var message chat.Message
		if err := c.ReadJSON(&message); err != nil {
			fmt.Printf("Failed to read message: %v\n", err)
		}

		if message.Owner == name {
			fmt.Fprintln(writer, color.InCyan(assembleMessage(message)))
		} else {
			fmt.Println(color.InYellow(assembleMessage(message)))
		}
	}
}

func assembleMessage(m chat.Message) string {
	return fmt.Sprintf(messagePattern, time.Unix(m.UnixTimestamp, 0), m.Owner, m.Text)
}

func makeDialURL(address *string, name *string) string {
	return *address + "/start?name=" + *name
}
