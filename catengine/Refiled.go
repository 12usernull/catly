package engine

import (
	"os"
)

var buffer []string
var primary *os.File
var replica *os.File

const BatchSize = 10

func Init() {
	os.MkdirAll("data", 0755)

	p, _ := os.OpenFile("data/primary.ndjson", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	r, _ := os.OpenFile("data/replica.ndjson", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)

	primary = p
	replica = r
}

func Write(msg string) {
	buffer = append(buffer, msg)

	if len(buffer) >= BatchSize {
		Flush()
	}
}

func Flush() {
	if len(buffer) == 0 {
		return
	}

	for _, m := range buffer {
		primary.WriteString(m + "\n")
		replica.WriteString(m + "\n")
	}

	primary.Sync()
	replica.Sync()

	buffer = buffer[:0]
}
