package gmetric

import (
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/viant/gmetric/counter"
	"github.com/viant/gmetric/stat"
)

// Service represents operation metrics
type Service struct {
	operations []Operation
	counters   []Counter

	lock *sync.Mutex
}

// OperationCounter register operation counters
func (s *Service) OperationCounter(location, name, description string, unit, loopbackUnit time.Duration, RecentBuckets int) *Operation {
	counter := NewOperation(location, name, description, RecentBuckets, loopbackUnit, unit, nil)

	s.lock.Lock()
	defer s.lock.Unlock()

	s.operations = append(s.operations, counter)
	return &counter
}

// MultiOperationCounter register multi value operation counters
func (s *Service) MultiOperationCounter(location, name, description string, unit, loopbackUnit time.Duration, loopbackSize int, provider counter.Provider) *Operation {
	counter := NewOperation(location, name, description, loopbackSize, loopbackUnit, unit, provider)

	s.lock.Lock()
	defer s.lock.Unlock()

	s.operations = append(s.operations, counter)
	return &s.operations[len(s.operations)-1]
}

// Counter register counters
func (s *Service) Counter(location, name, description string) *Counter {
	counter := NewCounter(location, name, description)

	s.lock.Lock()
	defer s.lock.Unlock()

	s.counters = append(s.counters, counter)
	return &s.counters[len(s.counters)-1]
}

// LookupOperation returns operation counters
func (s *Service) LookupOperation(name string) *Operation {
	for _, candidate := range s.operations {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}

var errMetric = errors.New("metric error")

// LookupOperationRecentMetric returns operation metric current bucket value
func (s *Service) LookupOperationRecentMetric(operationName, metric string) int64 {
	operation := s.LookupOperation(operationName)
	if operation == nil {
		return 0
	}
	recentIndex := operation.Index(time.Now())
	counterMetrics := operation.Recent[recentIndex]
	isPct := len(metric) > len(stat.CounterPctKey) && strings.HasSuffix(metric, stat.CounterPctKey)
	if isPct {
		if index := strings.LastIndex(metric, "."); index != -1 {
			metric = metric[:index]
		}
	}
	valueIndex := s.getMetricValueIndex(metric, operation)
	if isPct {
		metric = stat.CounterPctKey
	}
	return s.getCounterValue(metric, counterMetrics, valueIndex)
}

// LookupOperationRecentMetrics returns operation metrics current bucket values
func (s *Service) LookupOperationRecentMetrics(operationName string) counter.Operation {
	operation := s.LookupOperation(operationName)
	if operation == nil {
		return counter.Operation{}
	}
	recentIndex := operation.Index(time.Now())
	counterMetrics := operation.Recent[recentIndex]
	return *counterMetrics
}

// LookupOperationCumulativeMetric returns operation metric cumulative value
func (s *Service) LookupOperationCumulativeMetric(operationName, metric string) int64 {
	operation := s.LookupOperation(operationName)
	if operation == nil {
		return 0
	}
	isPct := len(metric) > len(stat.CounterPctKey) && strings.HasSuffix(metric, stat.CounterPctKey)
	if isPct {
		if index := strings.LastIndex(metric, "."); index != -1 {
			metric = metric[:index]
		}
	}
	valueIndex := s.getMetricValueIndex(metric, operation)
	if isPct {
		metric = stat.CounterPctKey
	}
	return s.getCounterValue(metric, operation.Operation.Operation, valueIndex)
}

func (s *Service) getCounterValue(metric string, operation *counter.Operation, valueIndex int) int64 {
	counterMetrics := operation.Counters
	switch metric {
	case stat.CounterValueKey:
		return operation.CountValue()
	case stat.CounterMineKey:
		return atomic.LoadInt64(&operation.Min)
	case stat.CounterMaxKey:
		return atomic.LoadInt64(&operation.Max)
	case stat.CounterAvgKey:
		return int64(atomic.LoadInt32(&operation.Avg))
	case stat.CounterTimeTakenKey:
		return atomic.LoadInt64(&operation.TimeTaken)
	case stat.CounterPctKey:
		if valueIndex >= 0 && valueIndex < len(counterMetrics) {
			return int64(atomic.LoadInt32(&counterMetrics[valueIndex].Pct))
		}
		return 0
	default:
		if valueIndex >= 0 && valueIndex < len(counterMetrics) {
			return counterMetrics[valueIndex].CountValue()
		}
		return 0
	}
}

func (s *Service) getMetricValueIndex(metric string, operation *Operation) int {
	var metricValue interface{} = metric
	if metric == stat.ErrorKey {
		metricValue = errMetric
	}
	valueIndex := operation.Provider.Map(metricValue)
	return valueIndex
}

// LookupCounter returns counters
func (s *Service) LookupCounter(name string) *Counter {
	for _, candidate := range s.counters {
		if candidate.Name == name {
			return &candidate
		}
	}
	return nil
}

// OperationCounters returns operation counters
func (s *Service) OperationCounters() []Operation {
	return s.operations
}

// FilteredOperationCounters returns operation counters
func (s *Service) FilteredOperationCounters(URI string) func() []Operation {
	var filtered []Operation
	for _, candidate := range s.operations {
		if candidate.Location != "" && strings.Contains(candidate.Location, strings.Trim(URI, "/")) {
			filtered = append(filtered, candidate)
		}
	}
	return func() []Operation {
		return filtered
	}
}

// Counters returns counters
func (s *Service) Counters() []Counter {
	return s.counters
}

// New creates a new metric Service
func New() *Service {
	return &Service{
		operations: make([]Operation, 0),
		counters:   make([]Counter, 0),
		lock:       new(sync.Mutex),
	}
}
