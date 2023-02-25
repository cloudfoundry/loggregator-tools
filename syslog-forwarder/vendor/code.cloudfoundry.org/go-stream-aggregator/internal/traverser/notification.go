package traverser

import "context"

//go:generate pubsub-gen --struct-name=code.cloudfoundry.org/go-stream-aggregator/internal/traverser.Notification --package=traverser --traverser=Traverser --output=$GOPATH/src/code.cloudfoundry.org/go-stream-aggregator/internal/traverser/notification_traverser.gen.go --blacklist-fields=Notification.Producer

// Producer is a type copy of streamaggregator.Producer
type Producer interface {
	// Produce implements streamaggregator.Producer
	Produce(ctx context.Context, request interface{}, c chan<- interface{})
}

// Notification is used in the PubSub to alert upon producers coming
// and going.
type Notification struct {
	Added    bool
	Key      string
	Producer Producer
}
