package ws

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"yourproject/auth"
	"yourproject/engine"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	Conn *websocket.Conn
	User string
	Send chan string
}

var Clients = map[*Client]bool{}
var Broadcast = make(chan string)

func Handler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")

	user, ok := auth.Validate(token)
	if !ok {
		http.Error(w, "unauthorized", 401)
		return
	}

	conn, _ := upgrader.Upgrade(w, r, nil)

	c := &Client{
		Conn: conn,
		User: user,
		Send: make(chan string),
	}

	Clients[c] = true

	go readLoop(c)
	go writeLoop(c)
}

func readLoop(c *Client) {
	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		final := fmt.Sprintf("[%s]: %s", c.User, string(msg))

		engine.Write(final)
		Broadcast <- final
	}
}

func writeLoop(c *Client) {
	for msg := range c.Send {
		c.Conn.WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
