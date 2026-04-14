package main

import (
	"fmt"
	"net/http"

	"yourproject/api"
	"yourproject/engine"
	"yourproject/ws"
)

func main() {

	// 🧠 storage engine
	engine.Init()
	engine.LoadIndex()

	// 💬 websocket chat
	http.HandleFunc("/ws", ws.Handler)

	// 📖 read API
	http.HandleFunc("/read", api.FastRead)

	// ✍️ write API
	http.HandleFunc("/write", func(w http.ResponseWriter, r *http.Request) {
		msg := r.URL.Query().Get("m")
		engine.Write(msg)
		w.Write([]byte("ok"))
	})

	// 🔥 broadcast loop (chat realtime)
	go ws.StartBroadcast()

	fmt.Println("🚀 FULL MODULAR SERVER :8080")
	fmt.Println("   /ws    -> websocket chat")
	fmt.Println("   /read  -> fast read API")
	fmt.Println("   /write -> write message API")

	http.ListenAndServe(":8080", nil)
}
