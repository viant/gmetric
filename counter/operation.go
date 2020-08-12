package counter

import (
	"sync/atomic"
	"time"
)

//OnDone on done callback
type OnDone func(end time.Time, values ...interface{}) int64

//Operation represents basic metrics
type Operation struct {
	*MultiCounter
	TimeTaken int64
	Max       int64
	Min       int64
	Avg       int32
	timeUnit  time.Duration
}

//Begin begin count an event
func (b *Operation) Begin(started time.Time) OnDone {
	if b.timeUnit == 0 {
		b.timeUnit = 1
	}
	count := b.Increment()
	return func(end time.Time, values ...interface{}) int64 {
		values = NormalizeValue(values)
		elapsed := int64(end.Sub(started) / b.timeUnit)
		timeTaken := atomic.AddInt64(&b.TimeTaken, elapsed)
		avgTime := time.Duration(timeTaken / count)
		atomic.StoreInt32(&b.Avg, int32(avgTime))
		if elapsed > atomic.LoadInt64(&b.Max) {
			atomic.StoreInt64(&b.Max, elapsed)
		}
		min := atomic.LoadInt64(&b.Min)
		if elapsed < min || min == 0 {
			atomic.StoreInt64(&b.Min, elapsed)
		}

		if len(values) > 0 && b.provider != nil {
			for i := range values {
				b.incrementValueBy(values[i], count, 1)
			}
		}
		return count
	}
}

//NewOperation creates operation metrics
func NewOperation(timeUnit time.Duration, provider Provider) *Operation {
	op := &Operation{
		MultiCounter: &MultiCounter{
			Counter:  &Counter{},
			provider: provider,
		},
	}
	if provider != nil {
		op.Counters = make([]*Value, len(provider.Keys()))
		for i, val := range provider.Keys() {
			op.Counters[i] = &Value{
				Value: val,
			}
		}
	}
	op.timeUnit = timeUnit

	return op
}
