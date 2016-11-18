package gmetric

const (
	//ApplicationName  name of this application
	ApplicationName = "MetricService"
	//ApplicationVersion version of this application
	ApplicationVersion = "0.1.0"
)

//CounterService represents a counter seric
type CounterService interface {

	//RegisterCounter register a counter for passed in package, name, unit, description, size and optional value provider
	RegisterCounter(packageName, name, unit, description string, size int, valueProvider ValueProvider) *OperationCounter

	//RegisterKeyCounter register a keyed counter for passed in package, name, unit, description, size and optional value provider
	RegisterKeyCounter(packageName, name, unit, description string, size int, keyProvider KeyProvider, valueProvider ValueProvider) *KeyedOperationCounter

	//Returns details about operational metric counters, it support wildcard expression at counter level:  com/vinat/app1/*  or exact counter name com/vinat/app1/CounterName
	Query(expression string) (map[string]*OperationMetricPackage, error)
}

//KeyProvider represents a key provider
type KeyProvider func(source interface{}) string

//ValueProvider represents a value provider
type ValueProvider func(source interface{}) int
