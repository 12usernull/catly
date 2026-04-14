package api

import (
	"bufio"
	"net/http"
	"os"
)

func ReadMessages(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data/primary.ndjson")
	if err != nil {
		http.Error(w, "no data", 500)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var result []string

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	for _, m := range result {
		w.Write([]byte(m + "\n"))
	}
}
