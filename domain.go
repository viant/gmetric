package gmetric

func (m *OperationMetric) Clone() *OperationMetric {
	result := OperationMetric{
		Name:m.Name,
		Unit:m.Unit,
		Description:m.Description,
		ErrorCount:m.ErrorCount,
		Count:m.Count,
		RecentValues:make([]int64, len(m.RecentValues)),
		Averages:make([]int64, len(m.Averages)),
		AvgValue:m.AvgValue,
		MinValue:m.MinValue,
		MaxValue:m.MaxValue,
	}
	for i, v := range m.RecentValues {
		result.RecentValues[i] = v
	}
	for i, v := range m.Averages {
		result.Averages[i] = v
	}
	return &result
}

func (m *KeyedOperationMetric) Clone() *KeyedOperationMetric {
	result := &KeyedOperationMetric{
		Metrics:make(map[string]*OperationMetric),
	}
	for k, v := range m.Metrics {
		result.Metrics[k] = v.Clone()
	}

	return result
}

func (m *OperationMetricPackage) Clone() *OperationMetricPackage {
	result := &OperationMetricPackage{
		Name:m.Name,
		Metrics:make( map[string]*OperationMetric),
		KeyedMetrics:make(map[string]*KeyedOperationMetric),
	}

	for k, v:=range m.Metrics {
		result.Metrics[k] = v.Clone()
	}
	for k, v:=range m.KeyedMetrics {
		result.KeyedMetrics[k] = v.Clone()
	}
	return result
}
