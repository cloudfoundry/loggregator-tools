package main

import (
	"flag"
	"fmt"
	"log"

	"code.cloudfoundry.org/loggregator-tools/slow_forwarder_consumer/pkg/ingress"
	"code.cloudfoundry.org/loggregator-tools/slow_forwarder_consumer/pkg/plumbing"

	"google.golang.org/grpc"
)

var (
	certFile = flag.String("cert", "", "cert to use to connect to rlp")
	keyFile  = flag.String("key", "", "key to use to connect to rlp")
	caFile   = flag.String("ca", "", "ca cert to use to connect to rlp")
	port     = flag.String("port", "", "port to listen for envelopes")
)

func main() {
	flag.Parse()

	serverCreds, err := plumbing.NewServerCredentials(
		*certFile,
		*keyFile,
		*caFile,
	)
	if err != nil {
		log.Fatalf("failed to configure server TLS: %s", err)
	}

	rx := ingress.NewReceiver()
	srv := ingress.NewServer(
		fmt.Sprintf("127.0.0.1:%s", *port),
		rx,
		grpc.Creds(serverCreds),
	)
	srv.Start()
}
