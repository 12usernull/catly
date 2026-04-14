package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn *websocket.Conn
	send chan string
}

var clients = make(map[*Client]bool)
var broadcast = make(chan string)

func handleMessages() {
	for {
		msg := <-broadcast

		for client := range clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(clients, client)
			}
		}
	}
}

func handleClientMessages(client *Client) {
	defer func() {
		client.conn.Close()
		delete(clients, client)
	}()

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			break
		}

		finalMsg := fmt.Sprintf("%s | %s", time.Now().Format("15:04:05"), string(msg))

		// broadcast
		broadcast <- finalMsg
	}
}

func writeToClient(client *Client) {
	for msg := range client.send {
		client.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{
		conn: conn,
		send: make(chan string),
	}

	clients[client] = true

	go writeToClient(client)
	go handleClientMessages(client)
}

func main() {
	http.HandleFunc("/ws", wsHandler)

	go handleMessages()

	fmt.Println("🚀 WebSocket server :8080/ws")

	http.ListenAndServe(":8080", nil)
}
