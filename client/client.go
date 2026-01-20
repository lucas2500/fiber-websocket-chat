package main

import (
	"bufio"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/fasthttp/websocket"
)

type Client struct {
	User    string
	Conn    *websocket.Conn
	MsgChan chan string
}

func main() {

	var user string

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Type your name: ")

	if scanner.Scan() {
		user = scanner.Text()
	}

	client := newClient(user)

	go client.sendMessage()
	go client.receiveMessage()

	fmt.Println("Type your message: ")

	for scanner.Scan() {

		message := scanner.Text()

		if message == "quit" {
			break
		}

		client.MsgChan <- message
	}

	client.Conn.Close()

	os.Exit(0)
}

func newClient(user string) *Client {

	u := url.URL{
		Scheme:   "ws",
		Host:     "127.0.0.1:3000",
		Path:     "/ws/NewChat",
		RawQuery: fmt.Sprintf("name=%s", user),
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.Fatal("dial:", err)
	}

	return &Client{
		User:    user,
		Conn:    conn,
		MsgChan: make(chan string),
	}
}

func (c *Client) sendMessage() {

	for msg := range c.MsgChan {

		err := c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))

		if err != nil {

			log.Println("Unable to write message:", err)

			return
		}
	}
}

func (c *Client) receiveMessage() {

	for {

		_, message, err := c.Conn.ReadMessage()

		if err != nil {

			log.Println("Unable to read message:", err)

			return
		}

		fmt.Println(string(message))
	}
}
