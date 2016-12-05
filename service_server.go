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
			metric.Unit = unit
			divider = int64(time.Millisecond)

		}
	case "s":
		if metric.Unit == "ns" {
			metric.Unit = unit
			divider = int64(time.Millisecond * 1000)

		}
	case "kb":
		if metric.Unit == "bytes" {
			metric.Unit = unit
			divider = int64(1000)

		}
	case "mb":
		if metric.Unit == "bytes" {
			metric.Unit = unit
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

func applyUnits(pkg *OperationMetricPackage, units map[string]string) {
	for _, v := range pkg.Metrics {
		if unit, found := units[v.Name]; found {
			applyUnit(v, unit)
		}
	}
	for _, v := range pkg.KeyedMetrics {
		for _, metric := range v.Metrics {
			if unit, found := units[metric.Name]; found {
				applyUnit(metric, unit)
			}
		}
	}
}

func (s *serviceServer) Query(context context.Context, request *QueryRequest) (response *QueryResponse, err error) {
	response = &QueryResponse{}
	metrics, err := s.service.Query(request.Query)
	if err != nil {
		response.Error = err.Error()
		return response, err
	}

	if request.Summary || len(request.Units) > 0 {
		for k, v := range metrics {
			metrics[k] = v.Clone()
		}
	}
	if request.Summary {
		for _, v := range metrics {
			applySummary(v)

		}
	}
	if len(request.Units) > 0 {
		for _, v := range metrics {
			applyUnits(v, request.Units)

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
