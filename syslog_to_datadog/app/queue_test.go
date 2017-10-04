package app_test

import (
	"fmt"

	"code.cloudfoundry.org/loggregator-tools/syslog_to_datadog/app"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Queue", func() {
	It("can push and pop items to/from the queue", func() {
		q := app.NewQueue(20)

		for i := 0; i < 20; i++ {
			q.Push([]byte(fmt.Sprintf("%d", i)))
		}

		var items [][]byte
		for i := 0; i < 20; i++ {
			b, _ := q.Pop()
			items = append(items, b)
		}

		Expect(items).To(HaveLen(20))
		Expect(items).ToNot(ContainElement(nil))
	})

	It("drops items when the queue is full", func(done Done) {
		defer close(done)
		q := app.NewQueue(20)

		for i := 0; i < 20; i++ {
			q.Push([]byte(fmt.Sprintf("%d", i)))
		}

		q.Push([]byte{'1'})
		q.Push([]byte{'2'})
	})

	It("does not block if buffer is empty", func(done Done) {
		defer close(done)
		q := app.NewQueue(20)

		item, ok := q.Pop()
		Expect(item).To(BeNil())
		Expect(ok).To(BeFalse())
	})
})
