package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	emitInterval := os.Getenv("EMIT_INTERVAL")

	d, err := time.ParseDuration(emitInterval)
	if err != nil {
		d = 3 * time.Millisecond
	}

	ticker := time.NewTicker(d)
	for range ticker.C {
		fmt.Println("LogEmitter: emitting log")
	}
}
