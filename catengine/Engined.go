package engine

import (
	"encoding/json"
	"os"
)

type IndexEntry struct {
	Segment int   `json:"segment"`
	Offset  int64 `json:"offset"`
	Length  int   `json:"length"`
}

var (
	buffer       []string
	Index        []IndexEntry
	primary      *os.File
	replica      *os.File
	offset int64 = 0
)

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

	for _, msg := range buffer {

		line := msg + "\n"

		n, _ := primary.WriteString(line)
		replica.WriteString(line)

		Index = append(Index, IndexEntry{
			Segment: 1,
			Offset:  offset,
			Length:  n,
		})

		offset += int64(n)
	}

	primary.Sync()
	replica.Sync()

	buffer = buffer[:0]

	SaveIndex()
}

func SaveIndex() {
	file, _ := os.Create("data/index.json")
	defer file.Close()

	json.NewEncoder(file).Encode(Index)
}

func LoadIndex() {
	file, err := os.Open("data/index.json")
	if err != nil {
		return
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&Index)
}
