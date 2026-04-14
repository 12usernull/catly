package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	BATCH_SIZE = 10
	MAX_SIZE   = 256 * 1024 // 256KB
)

var (
	buffer    []string
	primary   *os.File
	replica   *os.File
	segmentID = 1
	lastFlush = time.Now()
)

func openFiles() {
	var err error

	primary, err = os.OpenFile(fmt.Sprintf("../data/primary_%04d.ndjson", segmentID),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	replica, err = os.OpenFile(fmt.Sprintf("../data/replica_%04d.ndjson", segmentID),
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	fmt.Println("📁 Segment açıldı:", segmentID)
}

func rotate() {
	primary.Close()
	replica.Close()

	segmentID++
	openFiles()
}

func checkRotate() {
	info, _ := primary.Stat()

	if info.Size() >= MAX_SIZE {
		flush()
		rotate()
	}
}

func flush() {
	if len(buffer) == 0 {
		return
	}

	for _, msg := range buffer {
		primary.WriteString(msg + "\n")
		replica.WriteString(msg + "\n")
	}

	primary.Sync()
	replica.Sync()

	fmt.Println("⚡ FLUSH:", len(buffer), "mesaj")

	buffer = buffer[:0]
	lastFlush = time.Now()

	checkRotate()
}

func add(msg string) {
	buffer = append(buffer, msg)

	if len(buffer) >= BATCH_SIZE {
		flush()
	}
}

func timeoutFlush() {
	for {
		time.Sleep(200 * time.Millisecond)

		if time.Since(lastFlush) > time.Second {
			flush()
		}
	}
}

func main() {
	buffer = make([]string, 0, BATCH_SIZE)

	openFiles()

	go timeoutFlush()

	fmt.Println("🚀 Engine hazır (mesaj gir):")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		add(scanner.Text())
	}
}
