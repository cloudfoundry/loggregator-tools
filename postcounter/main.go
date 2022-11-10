package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	log.Print("Starting Post Counter...")
	defer log.Print("Closing Post Counter.")
	pt := &postTracker{}

	countDuration, err := time.ParseDuration(os.Getenv("DURATION"))
	if err != nil {
		countDuration = time.Minute
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			pt.countPost()
		} else {
			count := pt.getCounts(countDuration)
			_, _ = w.Write([]byte(fmt.Sprint(count)))
		}
	})

	fmt.Println(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}

type postTracker struct {
	sync.Mutex
	posts []time.Time
}

func (p *postTracker) countPost() {
	p.Lock()
	defer p.Unlock()

	p.posts = append(p.posts, time.Now())
}

func (p *postTracker) getCounts(inLast time.Duration) int {
	p.Lock()
	defer p.Unlock()

	threshold := time.Now().Add(-inLast)
	for i, t := range p.posts {
		if t.After(threshold) {
			count := len(p.posts) - i
			p.posts = p.posts[i:]

			return count
		}
	}

	return 0
}
