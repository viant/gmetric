package counter

import "sync/atomic"

//Counter represents a counter
type Counter struct {
	Count int64
}

//CountValue returns count
func (c *Counter) CountValue() int64 {
	return atomic.LoadInt64(&c.Count)
}

//IncrementBy increments counter
func (c *Counter) IncrementBy(i int64) int64 {
	return atomic.AddInt64(&c.Count, i)
}

//Increment increments counter
func (c *Counter) Increment() int64 {
	return c.IncrementBy(1)
}

//Decrement decrement counter
func (c *Counter) Decrement() int64 {
	return c.IncrementBy(-1)
}
