package gmetric_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/gmetric"
	"testing"
)

func TestService_RegisterKeyMetric(t *testing.T) {
	service := gmetric.NewCounterService()
	counter1 := service.RegisterKeyCounter("com/viant/gmetric", "testService1", "ns", "Test1", 10, nil, nil)
	assert.NotNil(t, counter1)

	counter1.Add("k1", 10, nil)
	assert.EqualValues(t, 10, counter1.Metrics["k1"].RecentValues[0])

	counter1.Add("k1", 20, nil)
	assert.EqualValues(t, 20, counter1.Metrics["k1"].RecentValues[1])

	counter2 := service.RegisterKeyCounter("com/viant/gmetric", "testService2", "ns", "Test2", 10, nil, nil)
	assert.NotNil(t, counter2)

}

func TestService_RegisterMetric(t *testing.T) {
	service := gmetric.NewCounterService()
	counter1 := service.RegisterCounter("com/viant/gmetric", "testService1", "ns", "Test1", 10, nil)
	assert.NotNil(t, counter1)

	counter1.Add(10, nil)
	assert.EqualValues(t, 10, counter1.RecentValues[0])

	counter1.Add(20, nil)
	assert.EqualValues(t, 20, counter1.RecentValues[1])

	counter2 := service.RegisterCounter("com/viant/gmetric", "testService2", "ns", "Test2", 10, nil)
	assert.NotNil(t, counter2)

}

func TestService_Query(t *testing.T) {
	service := gmetric.NewCounterService()

	{
		counter := service.RegisterKeyCounter("com/viant/gmetric", "Metric1", "ns", "Test1", 10, nil, nil)
		counter.Add("k1", 10, nil)

	}

	{
		counter := service.RegisterCounter("com/viant/gmetric", "Metric2", "ns", "Test1", 10, nil)
		counter.Add(20, nil)
	}

	{
		//query with wildcard
		packages, err := service.Query("com/viant/gmetric/*")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(packages))

		assert.Equal(t, 1, len(packages["com/viant/gmetric"].Metrics))
		assert.Equal(t, 1, len(packages["com/viant/gmetric"].KeyedMetrics))
	}

	{
		//query with metric name
		packages, err := service.Query("com/viant/gmetric/Metric2")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(packages))

		assert.Equal(t, 1, len(packages["com/viant/gmetric"].Metrics))
		assert.Equal(t, 0, len(packages["com/viant/gmetric"].KeyedMetrics))
	}

	{
		//query with metric name
		packages, err := service.Query("com/viant/gmetric/Metric1")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(packages))

		assert.Equal(t, 0, len(packages["com/viant/gmetric"].Metrics))
		assert.Equal(t, 1, len(packages["com/viant/gmetric"].KeyedMetrics))
	}

	{
		//query with metric name
		packages, err := service.Query("*")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(packages))

	}

	{
		//invlaid query with metric name
		_, err := service.Query("*Metric1")
		assert.NotNil(t, err)

	}
}
