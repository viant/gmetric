package gmetric

import (
	"sync"
	"time"
)

//KeyedOperationCounter represents a operation metrics by key/type
type KeyedOperationCounter struct {
	*KeyedOperationMetric
	KeyProvider   KeyProvider
	ValueProvider ValueProvider
	TrackerMap    map[string]*OperationCounter
	Name          string
	Unit          string
	Description   string
	Size          int
	mutex         *sync.RWMutex
}

func (t *KeyedOperationCounter) get(key string) *OperationCounter {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if result, found := t.TrackerMap[key]; found {
		return result
	}
	return nil
}

func (t *KeyedOperationCounter) getOrCreate(key string) *OperationCounter {
	result := t.get(key)
	if result != nil {
		return result
	}
	t.mutex.Lock()
	defer t.mutex.Unlock()
	result = NewOperationCounter(t.Name, t.Unit, t.Description, t.Size, t.ValueProvider)
	t.TrackerMap[key] = result
	t.Metrics[key] = result.OperationMetric
	return result
}

//AddFromSource adds to counter value from source evaluated   by value provider (specified at a counter creation type)
func (t *KeyedOperationCounter) AddFromSource(keySource, valueSource interface{}, err error) {
	key := t.KeyProvider(keySource)
	operationPerformance := t.getOrCreate(key)
	operationPerformance.AddFromSource(valueSource, err)
}

//Add adds to counter passed in value for supplied key
func (t *KeyedOperationCounter) Add(key string, value int, err error) {
	operationPerformance := t.getOrCreate(key)
	operationPerformance.Add(value, err)
}

//AddLatency adds to counter latency computed as difference between now and passed startTime.
func (t *KeyedOperationCounter) AddLatency(key string, startTime time.Time, err error) {
	operationPerformance := t.getOrCreate(key)
	operationPerformance.AddLatency(startTime, err)
}

//NewKeyedOperationCounter create a new key operation counter
func NewKeyedOperationCounter(name, unit, description string, size int, keyProvider KeyProvider, valueProvider ValueProvider) *KeyedOperationCounter {
	keyedOperationMetric := &KeyedOperationMetric{
		Metrics: make(map[string]*OperationMetric),
	}
	return &KeyedOperationCounter{
		KeyedOperationMetric: keyedOperationMetric,
		Name:                 name,
		Description:          description,
		Size:                 size,
		Unit:                 unit,
		KeyProvider:          keyProvider,
		ValueProvider:        valueProvider,
		TrackerMap:           make(map[string]*OperationCounter),
		mutex:                &sync.RWMutex{},
	}
}
