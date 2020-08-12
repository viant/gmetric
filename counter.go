package gmetric

import "github.com/viant/gmetric/counter"

//Counter represents a counter
type Counter struct {
	Identity
	counter.Counter
}

//Reset resets counters
func (c *Counter) Reset() {
	c.Counter = counter.Counter{}
}

//NewCounter creates a counters
func NewCounter(location, name, description string) Counter {
	return Counter{
		Identity: Identity{
			Location:    location,
			Name:        name,
			Description: description,
		},
		Counter: counter.Counter{},
	}
}
