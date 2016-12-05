package gmetric

import (
	"sync/atomic"
	"time"
)

const MaxInt = int64(^uint(0) >> 1)
const MinInt = -(MaxInt - 1)

//OperationCounter represents a application operation counter
type OperationCounter struct {
	*OperationMetric
	valueProvider ValueProvider
	averageIndex  int64
}

//AddLatency adds to counter bifferance between now and start time in ns.
func (t *OperationCounter) AddLatency(startTime time.Time, err error) {
	operationLength := len(t.RecentValues)
	if operationLength > 0 {
		t.Add(int(time.Now().UnixNano() - startTime.UnixNano()), err)
	} else {
		t.Add(0, err)
	}
}

//AddFromSource adds a value from source to a counter.
func (t *OperationCounter) AddFromSource(valueSource interface{}, err error) {
	value := t.valueProvider(valueSource)
	t.Add(value, err)
}

func (t *OperationCounter) computeRecentMetrics(limit int) (avg, min, max float64) {
	var cumulative int64 = 0
	min = float64(MaxInt)
	max = float64(MinInt)
	for i := 0; i < limit; i++ {
		recentValue := atomic.LoadInt64(&t.RecentValues[i])
		cumulative = cumulative + recentValue
		recentValueFloat := float64(recentValue)
		if recentValue != 0 && recentValueFloat < min {
			min = recentValueFloat
		}
		if recentValueFloat > max {
			max = recentValueFloat
		}
	}
	return float64(cumulative) / float64(limit), float64(min), float64(max)
}

//Add add a value to counter.
func (t *OperationCounter) Add(value int, err error) {
	if err != nil {
		atomic.AddUint64(&t.ErrorCount, 1)
	}
	operationLength := len(t.RecentValues)
	timeTakenAvgLength := len(t.Averages)
	var count = atomic.LoadUint64(&t.Count)
	if !atomic.CompareAndSwapUint64(&t.Count, count, count + 1) {
		t.Add(value, err)
		return
	}
	if operationLength > 0 {
		if int(count) > 0 {

			var limit = operationLength
			if (int(count) % operationLength) > 0 {
				limit = int(count) % operationLength
			}
			avg, min, max := t.computeRecentMetrics(limit)

			if (int(count) % operationLength) == 0 {
				avgIndex := int(t.averageIndex) % timeTakenAvgLength
				atomic.StoreInt64(&t.Averages[avgIndex], int64(avg))
				atomic.AddInt64(&t.averageIndex, 1)
			}
			t.AvgValue = avg
			if max == float64(MinInt) {
				max = 0
			}
			t.MaxValue = max
			if min == float64(MaxInt) {
				min = 0
			}
			t.MinValue = min
		}
		var index = int(count) % operationLength
		atomic.StoreInt64(&t.RecentValues[index % operationLength], int64(value))

	}
}

//NewOperationCounter create a new operation counter.
func NewOperationCounter(name, unit, description string, size int, valueProvider ValueProvider) *OperationCounter {
	result := &OperationCounter{
		valueProvider: valueProvider,
		OperationMetric: &OperationMetric{
			Name:         name,
			Unit:         unit,
			Description:  description,
			RecentValues: make([]int64, size),
			Averages:     make([]int64, size),
		}}
	return result
}
