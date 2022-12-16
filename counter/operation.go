package counter

import (
	"sync"
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

//Begin begins count an event
func (o *Operation) Begin(started time.Time) OnDone {
	return o.BeginWithInc(1, started)
}

//BeginWithInc begins count an event
func (o *Operation) BeginWithInc(inc int64, started time.Time) OnDone {
	if inc == 0 {
		inc = 1
	}
	if o.timeUnit == 0 {
		o.timeUnit = 1
	}
	count := o.IncrementBy(inc)
	return func(end time.Time, values ...interface{}) int64 {
		values = NormalizeValue(values)
		elapsed := int64(end.Sub(started) / o.timeUnit)
		timeTaken := atomic.AddInt64(&o.TimeTaken, elapsed)
		avgTime := time.Duration(timeTaken / count)
		atomic.StoreInt32(&o.Avg, int32(avgTime))
		if elapsed > atomic.LoadInt64(&o.Max) {
			atomic.StoreInt64(&o.Max, elapsed)
		}
		min := atomic.LoadInt64(&o.Min)
		if elapsed < min || min == 0 {
			atomic.StoreInt64(&o.Min, elapsed)
		}

		if len(values) > 0 && o.provider != nil {
			for i := range values {
				o.incrementValueBy(values[i], count, 1)
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
			locker:   &sync.Mutex{},
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
