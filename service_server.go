package gmetric

import (
	"golang.org/x/net/context"
	"time"
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
		for k, metricPointer := range keyMetric.Metrics {
			metric := *metricPointer
			metric.Averages = nil
			metric.RecentValues = nil
			pkg.KeyedMetrics[metricKey].Metrics[k] = &metric
		}
		pkg.KeyedMetrics[metricKey] = &keyMetric
	}
}

func applyUnit(metric *OperationMetric, unit string) {
	var divider int64 = 1;
	switch unit {
	case "ms":
		if metric.Unit == "ns" {
			divider = int64(time.Millisecond)

		}
	case "kbytes":
		if metric.Unit == "bytes" {
			divider = int64(1000)

		}
	case "mbytes":
		if metric.Unit == "bytes" {
			divider = int64(1000000)

		}
	}


	metric.MaxValue = metric.MaxValue / divider
	metric.MinValue = metric.MinValue / divider
	metric.AvgValue = metric.AvgValue / divider
	if len(metric.RecentValues) > 0 {
		var values = make([]int64, len(metric.RecentValues))
		for i, v := range metric.RecentValues {
			values[i] = v / divider
		}
		metric.RecentValues = values
	}
	if len(metric.Averages) > 0 {
		var values = make([]int64, len(metric.Averages))
		for i, v := range metric.Averages {
			values[i] = v / divider
		}
		metric.Averages = values
	}
}

func applyUnits(pkg *OperationMetricPackage, unit string) {
	for k, v := range pkg.Metrics {
		metric := *v
		applyUnit(&metric, unit)
		pkg.Metrics[k] = &metric
	}

	for metricKey, v := range pkg.KeyedMetrics {
		keyMetric := *v
		for k, metricPointer := range keyMetric.Metrics {
			metric := *metricPointer
			applyUnit(&metric, unit)
			metric.Averages = nil
			metric.RecentValues = nil
			pkg.KeyedMetrics[metricKey].Metrics[k] = &metric
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
	if request.Unit != "" {
		for k, v := range metrics {
			pkg := *v
			applyUnits(&pkg, request.Unit)
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
