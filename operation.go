package gmetric

import (
	"github.com/viant/gmetric/counter"
	"github.com/viant/gmetric/window"
	"time"
)

//Operation represents named counters metrics
type Operation struct {
	Identity
	window.Operation
	reset func() window.Operation
}

//Reset resets counters
func (c *Operation) Reset() {
	c.Operation = c.reset()
}

//NewOperation creates a new counters
func NewOperation(location, name, description string, size int, recentUnit, unit time.Duration, provider counter.Provider) Operation {
	newCounter := func() window.Operation { return window.NewOperation(size, recentUnit, unit, provider) }
	return Operation{
		Identity: Identity{
			Location:    location,
			Description: description,
			Name:        name,
		},
		Operation: newCounter(),
		reset:     newCounter,
	}
}
