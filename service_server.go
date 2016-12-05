package gmetric

import (
	"golang.org/x/net/context"
)

type serviceServer struct {
	service CounterService
}



func applySummary(pkg *OperationMetricPackage) {
	for k, v := range pkg.Metrics {
		metric := *v
		metric.Averages = nil
		metric.RecentValues = nil
		pkg.Metrics[k] = &metric
	}

	for metricKey, v := range pkg.KeyedMetrics {
		keyMetric := *v
		for k, metricPointer := range  keyMetric.Metrics {
			metric := *metricPointer
			metric.Averages = nil
			metric.RecentValues = nil
			pkg.KeyedMetrics[metricKey].Metrics[k]= &metric
		}
		pkg.KeyedMetrics[metricKey] = &keyMetric
	}
}


func (s *serviceServer) Query(context context.Context, request *QueryRequest) (response *QueryResponse, err error) {
	response = &QueryResponse{}
	metrics, err := s.service.Query(request.Query)
	if err != nil {
		response.Error = err.Error()
		return response, err
	}
	if request.Summary {
		for k, v := range metrics {
			pkg := *v
			applySummary(&pkg)
			metrics[k] = &pkg
		}
	}
	response.Metrics = metrics
	return response, err
}

//NewServiceServer creates a new service server for passed in service.
func mewServiceServer(service CounterService) ServiceServer {
	return &serviceServer{
		service: service,
	}
}
