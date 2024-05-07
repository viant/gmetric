package counter

import (
	"fmt"
	"github.com/viant/toolbox"
	"sync"
	"sync/atomic"
)

// MultiCounter represents multi value counter
type MultiCounter struct {
	*Counter
	Counters []*Value
	provider Provider
	locker   *sync.Mutex
}

// IncrementValue increments counter
func (c *MultiCounter) IncrementValue(value interface{}) int64 {
	return c.IncrementValueBy(value, 1)
}

// DecrementValue decrements counter by 1
func (c *MultiCounter) DecrementValue(value interface{}) int64 {
	return c.IncrementValueBy(value, -1)
}

// IncrementValueBy increments counter
func (c *MultiCounter) IncrementValueBy(value interface{}, i int64) int64 {
	return c.incrementValueBy(value, c.Counter.CountValue(), i)
}

// IncrementValueBy increments counter
func (c *MultiCounter) incrementValueBy(value interface{}, count, i int64) int64 {
	if len(c.Counters) == 0 {
		return 0
	}
	idx := c.provider.Map(value)
	if idx < 0 {
		return 0
	}
	if idx >= len(c.Counters) {
		idx = idx % len(c.Counters)
	}
	if _, ok := value.(string); !ok {
		c.locker.Lock()
		if custom, ok := value.(CustomCounter); ok {
			c.Counters[idx].Custom.Aggregate(custom)
		} else {
			stringer, ok := value.(fmt.Stringer)
			if ok {
				c.Counters[idx].Value = stringer.String()
			} else {
				c.Counters[idx].Value = toolbox.AsString(value)
			}
		}
		c.locker.Unlock()
	}
	valueCount := c.Counters[idx].IncrementBy(i)
	if count > 0 {
		atomic.SwapInt32(&c.Counters[idx].Pct, int32(100*valueCount/count))
	} else {
		atomic.SwapInt32(&c.Counters[idx].Pct, int32(0))
	}
	return count
}
