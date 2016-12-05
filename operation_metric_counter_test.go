package gmetric_test

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/viant/gmetric"
	"testing"
	"time"
)

func TestOperationMetricCounter_Add(t *testing.T) {
	counter := gmetric.NewOperationCounter("name", "ns", "test latency", 4, nil)
	err := testMethod1(10*time.Microsecond, counter, false)
	assert.Nil(t, err)
	assert.Equal(t, uint64(1), counter.Count)
	assert.Equal(t, uint64(0), counter.ErrorCount)

	assert.True(t, int(counter.RecentValues[0]) >= 10000)

	assert.Equal(t, int64(0), counter.Averages[0])
	for i := 0; i < 15; i++ {
		err = testMethod1(10*time.Microsecond, counter, i%4 == 0)
	}
	assert.Nil(t, err)
	assert.Equal(t, uint64(4), counter.ErrorCount)
	assert.True(t, int(counter.Averages[0]) >= 10000 && int(counter.Averages[0]) < 1000000)
}

func testMethod1(sleepTime time.Duration, counter *gmetric.OperationCounter, returnError bool) (err error) {
	defer func(startTime time.Time) {
		counter.AddLatency(startTime, err)
	}(time.Now())

	time.Sleep(sleepTime)
	if returnError {
		err = errors.New("test")
	}
	return err
}

func TestOperationMetricCounter_AddWithSource(t *testing.T) {
	counter := gmetric.NewOperationCounter("name", "bytes", "test payload", 2, func(source interface{}) int {
		if payload, casted := source.(*testPayload); casted {
			return (len(payload.bytes))
		}
		return 0
	})

	err := testMethod2(newTestPayload("test1"), counter, false)
	assert.Nil(t, err)

	err = testMethod2(newTestPayload("2"), counter, true)
	assert.NotNil(t, err)
	assert.Equal(t, uint64(2), counter.Count)
	assert.Equal(t, uint64(1), counter.ErrorCount)
	assert.Equal(t, int64(5), counter.RecentValues[0])
	assert.Equal(t, int64(1), counter.RecentValues[1])

	err = testMethod2(newTestPayload("12"), counter, true)
	assert.NotNil(t, err)

	assert.Equal(t, int64(3), counter.Averages[0])

	assert.Equal(t, int64(5), counter.MaxValue)
	assert.Equal(t, int64(1), counter.MinValue)

	err = testMethod2(newTestPayload("12"), counter, false)
	assert.Nil(t, err)

	err = testMethod2(newTestPayload("100000000"), counter, false)
	assert.Nil(t, err)
	assert.Equal(t, int64(2), counter.Averages[1])

}

func testMethod2(payload *testPayload, counter *gmetric.OperationCounter, returnError bool) (err error) {
	defer func(startTime time.Time) {
		counter.AddFromSource(payload, err)
	}(time.Now())

	if returnError {
		return errors.New("test error")
	}
	return nil
}

type testPayload struct {
	bytes []byte
}

func newTestPayload(fragment string) *testPayload {
	return &testPayload{
		bytes: []byte(fragment),
	}
}
