package api

import (
	"encoding/json"
	"net/http"
	"os"

	"yourproject/engine"
)

func FastRead(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open("data/primary.ndjson")
	if err != nil {
		http.Error(w, "no data", 500)
		return
	}
	defer file.Close()

	start := len(engine.Index) - 20
	if start < 0 {
		start = 0
	}

	result := []string{}

	for i := start; i < len(engine.Index); i++ {

		entry := engine.Index[i]

		buffer := make([]byte, entry.Length)

		file.Seek(entry.Offset, 0)
		file.Read(buffer)

		result = append(result, string(buffer))
	}

	json.NewEncoder(w).Encode(result)
}
