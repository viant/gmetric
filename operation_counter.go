package gmetric

import (
	"sync/atomic"
	"time"
)

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
		t.Add(int(time.Now().UnixNano()-startTime.UnixNano()), err)
	} else {
		t.Add(0, err)
	}
}

//AddFromSource adds a value from source to a counter.
func (t *OperationCounter) AddFromSource(valueSource interface{}, err error) {
	value := t.valueProvider(valueSource)
	t.Add(value, err)
}

//Add add a value to counter.
func (t *OperationCounter) Add(value int, err error) {
	if err != nil {
		atomic.AddInt64(&t.ErrorCount, 1)
	}
	operationLength := len(t.RecentValues)
	timeTakenAvgLength := len(t.Averages)
	var count = atomic.LoadInt64(&t.Count)
	if !atomic.CompareAndSwapInt64(&t.Count, count, count+1) {
		t.Add(value, err)
		return
	}
	if operationLength > 0 {
		if int(count) > 0 && (int(count)%operationLength) == 0 {
			var cumulative int64
			for i := 0; i < operationLength; i++ {
				cumulative = cumulative + atomic.LoadInt64(&t.RecentValues[i])
			}
			avgIndex := int(t.averageIndex) % timeTakenAvgLength
			atomic.StoreInt64(&t.Averages[avgIndex], cumulative/int64(operationLength))
			atomic.AddInt64(&t.averageIndex, 1)
		}
		var index = int(count) % operationLength
		atomic.StoreInt64(&t.RecentValues[index%operationLength], int64(value))

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
