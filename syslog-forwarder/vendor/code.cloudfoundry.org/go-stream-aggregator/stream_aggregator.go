package streamaggregator

import (
	"context"
	"io/ioutil"
	"log"
	"sync"

	"code.cloudfoundry.org/go-pubsub"
	"code.cloudfoundry.org/go-stream-aggregator/internal/traverser"
)

// StreamAggregator takes a dynamic list of producers and writes to interested
// consumers. When a consumer is removed, it will be told to not retry any
// connect logic via the given context being cancelled. When a producer is
// added, it will be instructed to write to all the subscribed consumers. If
// producers are available when a consumer is added, each producer will start
// writing to the consumer.
type StreamAggregator struct {
	ps  *pubsub.PubSub
	log *log.Logger

	mu         sync.Mutex
	globalList map[string]Producer
}

// New constructs a new StreamAggregator with the given options.
func New(opts ...StreamAggregatorOption) *StreamAggregator {
	a := &StreamAggregator{
		log:        log.New(ioutil.Discard, "", 0),
		globalList: make(map[string]Producer),
		ps:         pubsub.New(pubsub.WithNoMutex()),
	}

	for _, o := range opts {
		o.configure(a)
	}

	a.globalListSubscribe()

	return a
}

// StreamAggregatorOption is used to configure the StreamAggregator
type StreamAggregatorOption interface {
	configure(*StreamAggregator)
}

type streamAggregatorOptionFunc func(*StreamAggregator)

func (f streamAggregatorOptionFunc) configure(a *StreamAggregator) {
	f(a)
}

// WithLogger configures a logger for the StreamAggregator.
func WithLogger(l *log.Logger) StreamAggregatorOption {
	return streamAggregatorOptionFunc(func(a *StreamAggregator) {
		a.log = l
	})
}

// Producer is used to produce data for subscribed consumers.
type Producer interface {
	// Produce is invoked for each consumer. The given context will be cancelled
	// when the producer should stop trying to write data. Until the context
	// is cancelled, it is expected to write available data to the given
	// channel. The Producer must not close the channel. If the producer
	// encounters an error, it is expected to do its own retry logic.
	Produce(ctx context.Context, request interface{}, c chan<- interface{})
}

// ProducerFunc is an adapter to allow ordinary functions to be a Producer.
type ProducerFunc func(ctx context.Context, request interface{}, c chan<- interface{})

// Producer implements Producer.
func (f ProducerFunc) Produce(ctx context.Context, request interface{}, c chan<- interface{}) {
	f(ctx, request, c)
}

// AddProducer adds a producer to the StreamAggregator.
func (a *StreamAggregator) AddProducer(key string, p Producer) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ps.Publish(traverser.Notification{
		Added:    true,
		Key:      key,
		Producer: p,
	}, traverser.TraverserTraverse)
}

// RemoveProducer removes a producer from the StreamAggregator.
func (a *StreamAggregator) RemoveProducer(key string) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.ps.Publish(traverser.Notification{
		Added: false,
		Key:   key,
	}, traverser.TraverserTraverse)
}

// Consume starts consuming data from all the given Producers. As producers
// are added, each will start writing data to the consumer. The returned
// channel will not closed once the context is cancelled and all the producers
// exit.
func (a *StreamAggregator) Consume(ctx context.Context, request interface{}, opts ...ConsumeOption) <-chan interface{} {
	a.mu.Lock()
	defer a.mu.Unlock()

	var conf consumeConfig
	for _, o := range opts {
		o.configure(&conf)
	}

	c := make(chan interface{}, conf.length)
	var wg sync.WaitGroup

	var unsubscribes []func()

	for key, p := range a.globalList {
		unsubscribe := a.startProducer(ctx, key, request, c, p, &wg)
		unsubscribes = append(unsubscribes, unsubscribe)
	}

	unsubscribe := a.whenProducersAreChanged(true, "", func(n traverser.Notification) {
		u := a.startProducer(ctx, n.Key, request, c, n.Producer, &wg)
		unsubscribes = append(unsubscribes, u)
	})

	unsubscribes = append(unsubscribes, unsubscribe)

	go func() {
		<-ctx.Done()

		a.mu.Lock()
		for _, un := range unsubscribes {
			un()
		}
		a.mu.Unlock()

		// We can safely wait for the WaitGroup to drain because any new producers
		// added after the lock was released will not add to the WaitGroup now
		// that the context is cancelled.
		wg.Wait()
		close(c)
	}()

	return c
}

func (a *StreamAggregator) startProducer(
	ctx context.Context,
	key string,
	r interface{},
	c chan<- interface{},
	p Producer,
	wg *sync.WaitGroup,
) func() {
	// Ensure we aren't adding a producer when the context just got cancelled
	select {
	case <-ctx.Done():
		return func() {} // NOP
	default:
	}

	wg.Add(1)
	producerCtx, cancel := context.WithCancel(ctx)

	unsubscribe := a.whenProducersAreChanged(false, key, func(n traverser.Notification) {
		cancel()
	})

	go func() {
		defer wg.Done()
		p.Produce(producerCtx, r, c)
	}()

	return unsubscribe
}

// WithConsumeChannelLength sets the channel buffer length for the resulting
// channel.
func WithConsumeChannelLength(length int) ConsumeOption {
	return consumeOptionFunc(func(c *consumeConfig) {
		c.length = length
	})
}

type consumeConfig struct {
	length int
}

// ConsumeOption is used to configure a consumer. Defaults to 0.
type ConsumeOption interface {
	configure(*consumeConfig)
}

type consumeOptionFunc func(*consumeConfig)

func (f consumeOptionFunc) configure(c *consumeConfig) {
	f(c)
}

func (a *StreamAggregator) globalListSubscribe() {
	a.whenProducersAreChanged(true, "", func(n traverser.Notification) {
		a.globalList[n.Key] = n.Producer
	})

	a.whenProducersAreChanged(false, "", func(n traverser.Notification) {
		delete(a.globalList, n.Key)
	})
}

func (a *StreamAggregator) whenProducersAreChanged(
	added bool,
	key string,
	f func(n traverser.Notification),
) func() {
	filter := &traverser.NotificationFilter{
		Added: &added,
	}

	if key != "" {
		filter.Key = &key
	}

	path := traverser.TraverserCreatePath(filter)

	return a.ps.Subscribe(func(data interface{}) {
		n := data.(traverser.Notification)
		f(n)
	}, pubsub.WithPath(path))
}
