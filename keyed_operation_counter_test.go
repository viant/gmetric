package gmetric_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/viant/gmetric"
	"testing"
	"time"
)

func TestKeyedOperationMetricCounter_AddLatency(t *testing.T) {
	counter := gmetric.NewKeyedOperationCounter("test testKeyedMethod1", "ns", "testKeyedMethod1 Latency", 10, nil, nil)
	for _, key := range []string{"k1", "k2"} {
		err := testKeyedMethod1(10*time.Millisecond, key, counter, false)
		assert.Nil(t, err)
		assert.Equal(t, uint64(1), counter.Metrics[key].Count)
		assert.True(t, counter.Metrics[key].RecentValues[0] >= int64(10))

		err = testKeyedMethod1(10*time.Millisecond, key, counter, true)
		assert.NotNil(t, err)
		assert.Equal(t, uint64(1), counter.Metrics[key].ErrorCount)
		assert.Equal(t, uint64(2), counter.Metrics[key].Count)
	}
}

func testKeyedMethod1(sleepTime time.Duration, key string, counter *gmetric.KeyedOperationCounter, returnError bool) (err error) {
	defer func(startTime time.Time) {
		counter.AddLatency(key, startTime, err)
	}(time.Now())
	time.Sleep(sleepTime)
	if returnError {
		err = errors.New("test")
	}
	return err
}

func TestKeyedOperationMetricCounter_AddFromSource(t *testing.T) {
	counter := gmetric.NewKeyedOperationCounter("test testKeyedMethod1", "ns", "testKeyedMethod1 Latency", 10, func(source interface{}) string {
		if payload, casted := source.(*testKeyedPayload); casted {
			return payload.key
		}
		return ""
	}, func(source interface{}) int {
		if payload, casted := source.(*testKeyedPayload); casted {
			return (len(payload.bytes))
		}
		return 0
	})

	for _, key := range []string{"k1", "k2"} {
		err := testKeyedMethod2(newTestKeyPayload(key, "abcd"), counter, false)
		assert.Nil(t, err)

		assert.Equal(t, uint64(1), counter.Metrics[key].Count)
		assert.Equal(t, int64(4), counter.Metrics[key].RecentValues[0])
	}
	for _, key := range []string{"k3", "k4"} {
		err := testKeyedMethod3(newTestKeyPayload(key, "abcd"), counter, false)
		assert.Nil(t, err)

		assert.Equal(t, uint64(1), counter.Metrics[key].Count)
		assert.Equal(t, int64(4), counter.Metrics[key].RecentValues[0])
	}
}

func testKeyedMethod2(payload *testKeyedPayload, counter *gmetric.KeyedOperationCounter, returnError bool) (err error) {
	defer func(startTime time.Time) {
		counter.AddFromSource(payload, payload, err)
	}(time.Now())

	if returnError {
		return errors.New("test error")
	}
	return nil
}

func testKeyedMethod3(payload *testKeyedPayload, counter *gmetric.KeyedOperationCounter, returnError bool) (err error) {
	defer func(startTime time.Time) {
		counter.Add(payload.key, (len(payload.bytes)), err)
	}(time.Now())

	if returnError {
		return errors.New("test error")
	}
	return nil
}

type testKeyedPayload struct {
	key   string
	bytes []byte
}

func newTestKeyPayload(key, fragment string) *testKeyedPayload {
	return &testKeyedPayload{
		key:   key,
		bytes: []byte(fragment),
	}
}
