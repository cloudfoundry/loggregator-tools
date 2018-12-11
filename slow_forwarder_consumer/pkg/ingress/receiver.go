package ingress

import (
	"log"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"golang.org/x/net/context"
)

type Receiver struct {
	blockChan chan struct{}
}

func NewReceiver() *Receiver {
	return &Receiver{
		blockChan: make(chan struct{}),
	}
}

func (s *Receiver) Sender(sender loggregator_v2.Ingress_SenderServer) error {
	for {
		_, err := sender.Recv()
		if err != nil {
			log.Printf("Failed to receive data: %s", err)
			return err
		}
		<-s.blockChan
	}

	return nil
}

func (s *Receiver) BatchSender(sender loggregator_v2.Ingress_BatchSenderServer) error {
	for {
		_, err := sender.Recv()
		if err != nil {
			log.Printf("Failed to receive data: %s", err)
			return err
		}
		<-s.blockChan
	}

	return nil
}

func (s *Receiver) Send(_ context.Context, b *loggregator_v2.EnvelopeBatch) (*loggregator_v2.SendResponse, error) {
	for range b.Batch {
		<-s.blockChan
	}

	return &loggregator_v2.SendResponse{}, nil
}
