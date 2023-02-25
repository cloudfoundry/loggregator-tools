package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/loggregator-release/src/plumbing"
	"google.golang.org/grpc"
)

func main() {
	caFile := flag.String("ca", "", "CA certificate")
	certFile := flag.String("cert", "", "TLS certificate")
	keyFile := flag.String("key", "", "TLS private key")
	port := flag.Int("port", 8081, "Port of server")
	flag.Parse()

	serverCreds, err := plumbing.NewServerCredentials(
		*certFile,
		*keyFile,
		*caFile,
	)

	if err != nil {
		log.Fatalf("failed to configure server TLS: %s", err)
	}

	rx := newEchoReceiver()
	addr := fmt.Sprintf("127.0.0.1:%d", *port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Printf("grpc bound to: %s", lis.Addr())

	grpcServer := grpc.NewServer(grpc.Creds(serverCreds))
	loggregator_v2.RegisterIngressServer(grpcServer, rx)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type echoReceiver struct{}

func newEchoReceiver() *echoReceiver {
	return &echoReceiver{}
}

func (r *echoReceiver) Sender(sender loggregator_v2.Ingress_SenderServer) error {
	for {
		e, err := sender.Recv()
		if err != nil {
			log.Printf("Failed to receive data: %s", err)
			return err
		}

		fmt.Printf("%+v \n", e)
	}
}

func (r *echoReceiver) BatchSender(sender loggregator_v2.Ingress_BatchSenderServer) error {
	for {
		envs, err := sender.Recv()
		if err != nil {
			log.Printf("Failed to receive data: %s", err)
			return err
		}

		for _, e := range envs.Batch {
			fmt.Printf("%+v \n", e)
		}
	}
}

func (r *echoReceiver) Send(_ context.Context, b *loggregator_v2.EnvelopeBatch) (*loggregator_v2.SendResponse, error) {
	for _, e := range b.Batch {
		fmt.Printf("%+v \n", e)
	}

	return &loggregator_v2.SendResponse{}, nil
}
