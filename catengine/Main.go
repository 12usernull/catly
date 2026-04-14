package main

import (
	"fmt"
	"net/http"

	"yourproject/engine"
	"yourproject/ws"
)

func main() {
	engine.Init()

	http.HandleFunc("/ws", ws.Handler)

	go ws.StartBroadcast()

	fmt.Println("🚀 Modüler server :8080")
	http.ListenAndServe(":8080", nil)
}
