package app

// Queue is a thread safe queue.
type Queue struct {
	buffer chan []byte
}

// NewQueue returns a new Queue configured with the given size.
func NewQueue(size int) *Queue {
	return &Queue{
		buffer: make(chan []byte, size),
	}
}

// Push adds a slice of bytes to the queue. If the queue is full the given
// slice will be dropped.
func (q *Queue) Push(eb []byte) {
	select {
	case q.buffer <- eb:
	default:
		// Drop message
	}
}

// Pop pulls a slice of bytes off the queue and returns the slice. If the
// queue is empty the returned slice will be nil and the bool will be false.
func (q *Queue) Pop() ([]byte, bool) {
	select {
	case b := <-q.buffer:
		return b, true
	default:
		return nil, false
	}
}
