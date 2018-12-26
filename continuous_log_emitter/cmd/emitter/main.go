package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	logger *log.Logger = log.New(os.Stderr, "", log.LstdFlags)
)

func main() {
	sen := "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium, totam rem aperiam, ea	que ipsa quae ab illo inventore veritatis et quasi architecto beatae vitae dicta sunt explicabo."

	e := os.Getenv("EMIT_INTERVAL")
	emitInterval, err := time.ParseDuration(e)

	if err != nil {
		panic(err)
	}

	logMarker := fmt.Sprintf("[test-message]: %s", sen)

	writeLogs(emitInterval, logMarker)
}

func writeLogs(ei time.Duration, logMarker string) {
	ticker := time.NewTicker(ei)

	fmt.Printf("%s Message\n", logMarker)
	for range ticker.C {
		fmt.Printf("%s Message\n", logMarker)
	}
}
