package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	logcache "code.cloudfoundry.org/go-log-cache"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

type LatencyTestResult struct {
	Latency          float64   `json:"latency"`
	QueryTimes       []float64 `json:"queryTime"`
	AverageQueryTime float64   `json:"averageQueryTime"`
}

func measureLatency(ctx context.Context, reader logcache.Reader, name string) ([]byte, error) {
	var queryTimes []time.Duration
	for j := 0; j < 10; j++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		expectedLog := fmt.Sprintf("Test log - %d", rand.Int63())
		logStartTime := time.Now()
		fmt.Println(expectedLog)

		for i := 0; i < 100; i++ {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
			}

			// Look at 110 pages. This equates to 11000 possible envelopes. The
			// reliability test emits 10000 envelopes and could be clutter within the
			// cache.
			queryDurations, success := consumePages(
				110,
				logStartTime.Add(-7500*time.Millisecond),
				logStartTime.Add(7500*time.Millisecond),
				reader,
				name,
				func(envelopes []*loggregator_v2.Envelope) (bool, int64) {
					return lookForMsg(expectedLog, envelopes)
				},
			)
			queryTimes = append(queryTimes, queryDurations...)

			if success {
				var queryTimesMs []float64
				var totalQueryTimes time.Duration
				for _, qt := range queryTimes {
					totalQueryTimes += qt
					queryTimesMs = append(queryTimesMs, float64(qt)/float64(time.Millisecond))
				}

				avgQT := int64(totalQueryTimes) / int64(len(queryTimes))

				// Success
				testResults := LatencyTestResult{
					Latency:          float64(time.Since(logStartTime)) / float64(time.Millisecond),
					QueryTimes:       queryTimesMs,
					AverageQueryTime: float64(avgQT) / float64(time.Millisecond),
				}

				resultData, err := json.Marshal(testResults)
				if err != nil {
					return nil, err
				}

				return resultData, nil
			}
		}
	}

	// Fail...
	return nil, errors.New("Never found the log")
}

func consumePages(pages int, start, end time.Time, reader logcache.Reader, name string, f func([]*loggregator_v2.Envelope) (bool, int64)) ([]time.Duration, bool) {
	var queryDurations []time.Duration

	for i := 0; i < pages; i++ {
		queryDuration, last, count, err := consumeTimeRange(start, end, reader, name, f)
		if err != nil {
			log.Println(err)
			return nil, false
		}

		queryDurations = append(queryDurations, queryDuration)
		if !last.IsZero() && count > 0 {
			start = last.Add(time.Nanosecond)
			continue
		}

		return queryDurations, true
	}

	return queryDurations, false
}

func consumeTimeRange(start, end time.Time, reader logcache.Reader, name string, f func([]*loggregator_v2.Envelope) (bool, int64)) (time.Duration, time.Time, int, error) {
	queryStart := time.Now()
	data, err := reader(
		context.Background(),
		name,
		start,
		logcache.WithEndTime(end),
	)

	if err != nil {
		return 0, time.Time{}, 0, err
	}
	queryDuration := time.Since(queryStart)

	found, last := f(data)
	if !found {
		return queryDuration, time.Unix(0, last), len(data), nil
	}

	return queryDuration, time.Time{}, len(data), nil
}

func lookForMsg(msg string, envelopes []*loggregator_v2.Envelope) (bool, int64) {
	var last int64
	for _, e := range envelopes {
		last = e.Timestamp
		if string(e.GetLog().GetPayload()) == msg {
			return true, last
		}
	}

	return false, last
}
