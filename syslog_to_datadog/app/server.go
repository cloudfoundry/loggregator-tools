package app

import (
	"log"
	"net"
	"net/http"
	"time"

	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/internal/processor"
	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/internal/web"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// Server handles initilizing an HTTP server as well as a rfc 5425 message
// processor.
type Server struct {
	listener       net.Listener
	queue          *Queue
	datadogBaseURL string
	datadogKey     string
}

// ServerOption is a func that can be passed into NewServer to configure
// optional settings on the Server.
type ServerOption func(*Server)

// WithDatadogBaseURL sets the base url for the datadog API.
func WithDatadogBaseURL(baseURL string) ServerOption {
	return func(s *Server) {
		s.datadogBaseURL = baseURL
	}
}

// NewServer will start a listener on the given addr and configure a server
// for processing rfc5424 messages.
func NewServer(addr, datadogKey string, opts ...ServerOption) *Server {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to open listener (%s): %s", addr, err)
	}
	log.Printf("listening on: %s", addr)

	queue := NewQueue(10000)
	server := &Server{
		listener:       lis,
		queue:          queue,
		datadogBaseURL: "https://app.datadoghq.com",
		datadogKey:     datadogKey,
	}
	for _, o := range opts {
		o(server)
	}
	return server
}

// Run will start the message processor and serve the HTTP server.
func (s *Server) Run() {
	go processor.New(
		s.queue.Pop,
		httpClient,
		s.datadogBaseURL,
		s.datadogKey,
	).Run()
	httpServer := &http.Server{Handler: web.NewHandler(s.queue.Push)}
	log.Fatalf("http server shutting down: %s", httpServer.Serve(s.listener))
}

// Addr returns the address the listener is bound to.
func (s *Server) Addr() string {
	return s.listener.Addr().String()
}
