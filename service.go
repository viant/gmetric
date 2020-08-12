package gmetric

import (
	"errors"
	"github.com/viant/gmetric/counter"
	"github.com/viant/gmetric/stat"
	"time"
)

//Service represents operation metrics
type Service struct {
	operations []Operation
	counters   []Counter
}

//OperationCounter register operation counters
func (s *Service) OperationCounter(location, name, description string, unit, loopbackUnit time.Duration, RecentBuckets int) *Operation {
	counter := NewOperation(location, name, description, RecentBuckets, loopbackUnit, unit, nil)
	s.operations = append(s.operations, counter)
	return &counter
}

//MultiOperationCounter register multi value operation counters
func (s *Service) MultiOperationCounter(location, name, description string, unit, loopbackUnit time.Duration, loopbackSize int, provider counter.Provider) *Operation {
	counter := NewOperation(location, name, description, loopbackSize, loopbackUnit, unit, provider)
	s.operations = append(s.operations, counter)
	return &s.operations[len(s.operations)-1]
}

//Counter register counters
func (s *Service) Counter(location, name, description string) *Counter {
	counter := NewCounter(location, name, description)
	s.counters = append(s.counters, counter)
	return &s.counters[len(s.counters)-1]
}

//LookupOperation returns operation counters
func (s *Service) LookupOperation(name string) *Operation {
	for _, candidate := range s.operations {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}

var errMetric = errors.New("metric error")

//LookupOperationRecentMetric returns operation metric current bucket value
func (s *Service) LookupOperationRecentMetric(operationName, metric string) int64 {
	operation := s.LookupOperation(operationName)
	if operation == nil {
		return 0
	}
	recentIndex := operation.Index(time.Now())
	counterMetrics := operation.Recent[recentIndex].Counters
	if metric == stat.CounterValueKey {
		return operation.Recent[recentIndex].CountValue()
	}
	valueIndex := s.getMetricValueIndex(metric, operation)
	if valueIndex >= 0 && valueIndex < len(counterMetrics) {
		return counterMetrics[valueIndex].CountValue()
	}
	return 0
}

//LookupOperationCumulativeMetric returns operation metric cumulative value
func (s *Service) LookupOperationCumulativeMetric(operationName, metric string) int64 {
	operation := s.LookupOperation(operationName)
	if operation == nil {
		return 0
	}
	counterMetrics := operation.Counters
	if metric == stat.CounterValueKey {
		return operation.CountValue()
	}
	valueIndex := s.getMetricValueIndex(metric, operation)
	if valueIndex >= 0 && valueIndex < len(counterMetrics) {
		return counterMetrics[valueIndex].CountValue()
	}
	return 0
}

func (s *Service) getMetricValueIndex(metric string, operation *Operation) int {
	var metricValue interface{} = metric
	if metric == stat.ErrorKey {
		metricValue = errMetric
	}
	valueIndex := operation.Provider.Map(metricValue)
	return valueIndex
}

//LookupCounter returns counters
func (s *Service) LookupCounter(name string) *Counter {
	for _, candidate := range s.counters {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}

//OperationCounters returns operation counters
func (s *Service) OperationCounters() []Operation {
	return s.operations
}

//Counters returns counters
func (s *Service) Counters() []Counter {
	return s.counters
}

//New creates a new metric Service
func New() *Service {
	return &Service{
		operations: make([]Operation, 0),
		counters:   make([]Counter, 0),
	}
}
