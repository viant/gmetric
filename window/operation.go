package window

import (
	"fmt"
	"github.com/viant/gmetric/counter"
	"sync/atomic"
	"time"
)

//OnCounterDone represents on on counter callback
type OnCounterDone func(end time.Time, flag bool, err error) int64

//Operation represents operation metrics
type Operation struct {
	*counter.Operation
	Recent             []*counter.Operation
	index              int32
	UnitDuration       time.Duration `json:"-"`
	Unit               string
	RecentUnitDuration time.Duration `json:"-"`
	RecentUnit         string
	Provider           counter.Provider `json:"-"`
}

//Begin start metrics
func (c *Operation) Begin(started time.Time) counter.OnDone {
	totalOnDone := c.Operation.Begin(started)
	index := c.Index(started)
	currentIndex := atomic.LoadInt32(&c.index)
	if currentIndex != index && atomic.CompareAndSwapInt32(&c.index, currentIndex, index) {
		c.Recent[index] = counter.NewOperation(c.UnitDuration, c.Provider)
	}
	recentOnDone := c.Recent[index].Begin(started)
	return func(end time.Time, values ...interface{}) int64 {
		values = counter.NormalizeValue(values)

		count := totalOnDone(end, values...)
		_ = recentOnDone(end, values...)
		return count
	}
}

//Index returns recent bucket index for supplied time
func (c *Operation) Index(atTime time.Time) int32 {
	recentIndex := int(time.Duration(atTime.UnixNano()) / c.RecentUnitDuration)
	index := int32(recentIndex % len(c.Recent))
	return index
}

//NewOperation creates basic rolling window counter
func NewOperation(recentBuckets int, recentUnit time.Duration, unit time.Duration, provider counter.Provider) Operation {
	if recentBuckets < 1 {
		recentBuckets = 2
	}
	if recentUnit == 0 {
		recentUnit = time.Minute
	}
	result := Operation{
		Operation:          counter.NewOperation(unit, provider),
		Recent:             make([]*counter.Operation, recentBuckets),
		index:              0,
		UnitDuration:       unit,
		Unit:               fmt.Sprintf("%s", unit),
		RecentUnitDuration: recentUnit,
		RecentUnit:         fmt.Sprintf("%s", recentUnit),
		Provider:           provider,
	}
	for i := range result.Recent {
		result.Recent[i] = counter.NewOperation(unit, provider)
	}
	return result
}
