package gmetric

import (
	"golang.org/x/net/context"
	"time"
)

type serviceServer struct {
	service CounterService
}

func applySummary(pkg *OperationMetricPackage) {
	for _, v := range pkg.Metrics {
		v.Averages = nil
		v.RecentValues = nil
	}
	for _, v := range pkg.KeyedMetrics {
		for _, metric := range v.Metrics {
			metric.Averages = nil
			metric.RecentValues = nil
		}
	}
}

func applyUnit(metric *OperationMetric, unit string) {
	var divider int64 = 1;
	switch unit {
	case "ms":
		if metric.Unit == "ns" {
			divider = int64(time.Millisecond)

		}
	case "s":
		if metric.Unit == "ns" {
			divider = int64(time.Millisecond * 1000)

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
		for i, v := range metric.RecentValues {
			metric.RecentValues [i] = v / divider
		}
	}
	if len(metric.Averages) > 0 {
		for i, v := range metric.Averages {
			metric.Averages[i] = v / divider
		}
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
			keyMetric.Metrics[k] = &metric
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

	if request.Summary || request.Unit != "" {
		for k, v := range metrics {
			metrics[k] = v.Clone()
		}
	}
	if request.Summary {
		for _, v := range metrics {
			applySummary(v)

		}
	}
	if request.Unit != "" {
		for _, v := range metrics {
			applyUnits(v, request.Unit)

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
