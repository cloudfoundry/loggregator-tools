package client_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync/atomic"

	sharedapi "code.cloudfoundry.org/loggregator-tools/reliability/api"
	"code.cloudfoundry.org/loggregator-tools/reliability/worker/internal/client"

	"github.com/gorilla/websocket"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("WorkerClient", func() {
	It("receives tests to run", func() {
		server := newFakeWSServer()
		runner := &spyRunner{}

		client := client.NewWorkerClient(server.wsAddr(), true, runner)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			err := client.Run(ctx)
			Expect(err).ToNot(HaveOccurred())
		}()
		Eventually(server.connections).Should(Equal(int64(1)))

		server.tests <- sharedapi.Test{}
		server.tests <- sharedapi.Test{}
		Eventually(runner.Count).Should(Equal(int64(2)))
	})
})

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type fakeWSServer struct {
	listener net.Listener
	tests    chan sharedapi.Test

	_connections int64
}

func newFakeWSServer() *fakeWSServer {
	server := &fakeWSServer{
		tests: make(chan sharedapi.Test, 100),
	}
	http.Handle("/", server)

	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}
	server.listener = lis

	go func() {
		log.Println(http.Serve(lis, nil))
	}()

	return server
}

func (f *fakeWSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	defer conn.Close()

	atomic.AddInt64(&f._connections, 1)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		defer cancel()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()

	for {
		select {
		case test := <-f.tests:
			err := conn.WriteJSON(&test)
			if err != nil {
				panic(err)
			}
		case <-ctx.Done():
			return
		}
	}

}

func (f *fakeWSServer) wsAddr() string {
	return "ws://" + f.listener.Addr().String()
}

func (f *fakeWSServer) connections() int64 {
	return atomic.LoadInt64(&f._connections)
}

func (f *fakeWSServer) stop() {
	err := f.listener.Close()
	if err != nil {
		panic(err)
	}
}

type spyRunner struct {
	client.Runner
	runCallCount int64
}

func (s *spyRunner) Run(*sharedapi.Test) {
	atomic.AddInt64(&s.runCallCount, 1)
}

func (s *spyRunner) Count() int64 {
	return atomic.LoadInt64(&s.runCallCount)
}
