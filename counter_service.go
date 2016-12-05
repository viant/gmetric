package gmetric

import (
	"fmt"
	"strings"
	"sync"
)

type counterService struct {
	mutex *sync.RWMutex
	Map   map[string]*SynchronizedOperationMetricPackage
}

func (s *counterService) getPackage(name string) *SynchronizedOperationMetricPackage {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	if result, found := s.Map[name]; found {
		return result
	}
	return nil
}

func (s *counterService) registerPackage(operationMetricPackage *SynchronizedOperationMetricPackage) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.Map[operationMetricPackage.Name] = operationMetricPackage
}

func (s *counterService) getOrCreatePackage(packageName string) *SynchronizedOperationMetricPackage {
	operationMetricPackage := s.getPackage(packageName)
	if operationMetricPackage == nil {
		operationMetricPackage = NewSynchronizedOperationMetricPackage(packageName)
		s.registerPackage(operationMetricPackage)
	}
	return operationMetricPackage
}

func (s *counterService) RegisterCounter(packageName, name, unit, description string, size int, valueProvider ValueProvider) *OperationCounter {
	operationMetricPackage := s.getOrCreatePackage(packageName)
	counter := NewOperationCounter(name, unit, description, size, valueProvider)
	operationMetricPackage.PutCounter(counter)
	return counter
}

func (s *counterService) RegisterKeyCounter(packageName, name, unit, description string, size int, keyProvider KeyProvider, valueProvider ValueProvider) *KeyedOperationCounter {
	operationMetricPackage := s.getOrCreatePackage(packageName)
	counter := NewKeyedOperationCounter(name, unit, description, size, keyProvider, valueProvider)
	operationMetricPackage.PutKeyeCounter(counter)
	return counter
}

//NewCounterService create a new instance of CounterService
func NewCounterService() CounterService {
	return &counterService{
		mutex: &sync.RWMutex{},
		Map:   make(map[string]*SynchronizedOperationMetricPackage),
	}
}

//QueryExpression represents  QueryExpression
type QueryExpression struct {
	packageExpr, metricExpr string
}

//BuildQueryExpression build an expression for passed in text
func BuildQueryExpression(expression string) *QueryExpression {
	theLastDotPosition := strings.LastIndex(expression, "/")
	if theLastDotPosition != -1 {
		packageExpr := expression[0:theLastDotPosition]
		metricExpr := expression[theLastDotPosition+1:]
		return &QueryExpression{
			packageExpr: packageExpr,
			metricExpr:  metricExpr,
		}
	} else if expression == "*" {
		return &QueryExpression{
			packageExpr: "*",
			metricExpr:  "*",
		}
	}
	return nil
}

func (s *counterService) Query(queryExpression string) (map[string]*OperationMetricPackage, error) {
	expression := BuildQueryExpression(queryExpression)
	if expression == nil {
		return nil, fmt.Errorf("Invlid expressoin %v", queryExpression)
	}
	var result = make(map[string]*OperationMetricPackage)
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	for k, v := range s.Map {
		if expression.packageExpr == k || expression.packageExpr == "*" {
			if expression.metricExpr == "*" {
				result[k] = v.OperationMetricPackage
			} else {
				result[k] = NewSynchronizedOperationMetricPackage(v.Name).OperationMetricPackage
				//if expression.metricExpr == v.Name
				for metricName, metric := range v.Metrics {
					if metricName == expression.metricExpr {
						result[k].Metrics[metricName] = metric
					}
				}
				for metricName, metric := range v.KeyedMetrics {
					if metricName == expression.metricExpr {
						result[k].KeyedMetrics[metricName] = metric
					}
				}
			}
		}

	}
	return result, nil
}
