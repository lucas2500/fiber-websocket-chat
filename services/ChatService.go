package services

import (
	"log"
	"sync"

	"github.com/gofiber/contrib/websocket"
)

type Client struct {
	User   string
	RoomId string
	Conn   *websocket.Conn
}

type Chat struct {
	// Clients    map[*Client]bool
	Rooms      map[*Client]string
	Broadcast  chan Broadcast
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.Mutex
}

type Broadcast struct {
	RoomId  string
	Conn    *websocket.Conn
	Message []byte
}

var chat = setChatAttributes()

func NewChat(c *websocket.Conn) {

	user := c.Query("name")
	RoomId := c.Query("RoomId")

	if RoomId == "" {
		RoomId = "public"
	}

	client := &Client{User: user, RoomId: RoomId, Conn: c}
	chat.Register <- client

	defer func() {
		chat.Unregister <- client
	}()

	for {
		_, msg, err := c.ReadMessage()

		NewMsg := user + ": " + string(msg)

		if err != nil {

			log.Printf("Unable to read message: %s\n", err)

			break
		}

		chat.Broadcast <- Broadcast{
			RoomId:  RoomId,
			Conn:    c,
			Message: []byte(NewMsg),
		}
	}
}

func RunChat() {

	for {

		select {
		case client := <-chat.Register:

			chat.Mu.Lock()

			chat.Rooms[client] = client.RoomId

			// chat.Clients[client] = true
			chat.Mu.Unlock()

			log.Printf("%s is connected on room %s\n", client.User, client.RoomId)

		case client := <-chat.Unregister:

			chat.Mu.Lock()

			if _, exists := chat.Rooms[client]; exists {
				delete(chat.Rooms, client)
				client.Conn.Close()
			}

			chat.Mu.Unlock()

			log.Printf("%s disconnected\n", client.User)

		case message := <-chat.Broadcast:

			chat.Mu.Lock()

			for client := range chat.Rooms {

				if client.RoomId == message.RoomId && client.Conn != message.Conn {

					if err := client.Conn.WriteMessage(websocket.TextMessage, message.Message); err != nil {

						log.Printf("Unable to send message: %s\n", err)

						client.Conn.Close()
						delete(chat.Rooms, client)
					}
				}
			}

			chat.Mu.Unlock()
		}
	}
}

func setChatAttributes() Chat {

	return Chat{
		// Clients:    make(map[*Client]bool),
		Rooms:      make(map[*Client]string),
		Broadcast:  make(chan Broadcast),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}
