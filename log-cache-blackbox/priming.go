package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
)

func prime(ctx context.Context, id string, reader logcache.Reader) error {
	done := make(chan struct{})
	defer close(done)

	start := time.Now()

	primerMessage := []byte(fmt.Sprintf("prime %s %d", id, time.Now().UnixNano()))
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-done:
				return
			default:
				log.Println(string(primerMessage))
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		time.Sleep(500 * time.Millisecond)

		select {
		case <-ctx.Done():
			return errors.New("context canceled, aborting priming")
		default:
		}

		log.Printf("primer id: %s", id)
		envs, err := reader(ctx, id, start)
		if err != nil {
			log.Printf("failed to read from id: %s: %s", id, err)
			continue
		}
		for _, e := range envs {
			if bytes.Contains(e.GetLog().GetPayload(), primerMessage) {
				return nil
			}
		}
	}
}
