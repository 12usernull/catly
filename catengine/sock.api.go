package api

import (
	"encoding/json"
	"net/http"
	"time"

	"yourproject/engine"
)

type SendReq struct {
	User string `json:"user"`
	Text string `json:"text"`
}

func WriteMessage(w http.ResponseWriter, r *http.Request) {
	var req SendReq
	json.NewDecoder(r.Body).Decode(&req)

	msg := "[" + req.User + "]: " + req.Text

	engine.Write(msg)

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().String(),
	})
}
