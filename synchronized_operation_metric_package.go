package gmetric

import "sync"

//SynchronizedOperationMetricPackage represents a sync SynchronizedOperationMetricPackage
type SynchronizedOperationMetricPackage struct {
	mutex *sync.RWMutex
	*OperationMetricPackage
	MetricCounters      map[string]*OperationCounter
	KeyedMetricCounters map[string]*KeyedOperationCounter
}

//GetMetric retusna metric for passed in key
func (t *SynchronizedOperationMetricPackage) GetMetric(key string) *OperationMetric {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if result, found := t.Metrics[key]; found {
		return result
	}
	return nil
}

//GetCounter returns an OperationCounter for passed in key
func (t *SynchronizedOperationMetricPackage) GetCounter(key string) *OperationCounter {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if result, found := t.MetricCounters[key]; found {
		return result
	}
	return nil
}

// PutCounter adds provided OperationCounter.
func (t *SynchronizedOperationMetricPackage) PutCounter(counter *OperationCounter) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.MetricCounters[counter.Name] = counter
	t.Metrics[counter.Name] = counter.OperationMetric
}

//GetKeyedMetric returns KeyedOperationMetric for passed in key
func (t *SynchronizedOperationMetricPackage) GetKeyedMetric(key string) *KeyedOperationMetric {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if result, found := t.KeyedMetrics[key]; found {
		return result
	}
	return nil
}

//GetKeyedCounter returns an GetKeyedCounter for passed in key.
func (t *SynchronizedOperationMetricPackage) GetKeyedCounter(key string) *KeyedOperationCounter {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	if result, found := t.KeyedMetricCounters[key]; found {
		return result
	}
	return nil
}

//PutKeyeCounter adds a PutKeyeCounter
func (t *SynchronizedOperationMetricPackage) PutKeyeCounter(counter *KeyedOperationCounter) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.KeyedMetricCounters[counter.Name] = counter
	t.KeyedMetrics[counter.Name] = counter.KeyedOperationMetric

}

//NewSynchronizedOperationMetricPackage create a new NewSynchronizedOperationMetricPackage.
func NewSynchronizedOperationMetricPackage(name string) *SynchronizedOperationMetricPackage {
	operationMetricPackage := &OperationMetricPackage{
		Name:         name,
		Metrics:      make(map[string]*OperationMetric),
		KeyedMetrics: make(map[string]*KeyedOperationMetric),
	}
	return &SynchronizedOperationMetricPackage{
		mutex: &sync.RWMutex{},
		OperationMetricPackage: operationMetricPackage,
		MetricCounters:         make(map[string]*OperationCounter),
		KeyedMetricCounters:    make(map[string]*KeyedOperationCounter),
	}
}
