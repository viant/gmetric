package counter

import (
	"github.com/viant/toolbox"
)

//MultiCounter represents multi value counter
type MultiCounter struct {
	*Counter
	Counters []*Value
	provider Provider
}

//IncrementValue increments counter
func (c *MultiCounter) IncrementValue(value interface{}) int64 {
	return c.IncrementValueBy(value, 1)
}

//IncrementValue increments counter
func (c *MultiCounter) DecrementValue(value interface{}) int64 {
	return c.IncrementValueBy(value, -1)
}

//IncrementValueBy increments counter
func (c *MultiCounter) IncrementValueBy(value interface{}, i int64) int64 {
	return c.incrementValueBy(value, c.Counter.CountValue(), i)
}

//IncrementValueBy increments counter
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
		c.Counters[idx].Value = toolbox.AsString(value)
	}

	valueCount := c.Counters[idx].IncrementBy(i)
	if count > 0 {
		c.Counters[idx].Pct = int32(100 * valueCount / count)
	}
	return count
}
