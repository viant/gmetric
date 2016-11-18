package gmetric_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/gmetric"
	"testing"
)

func TestNewSynchronizedOperationMetricPackage(t *testing.T) {
	pkg := gmetric.NewSynchronizedOperationMetricPackage("com/viant")

	assert.Nil(t, pkg.GetKeyedMetric("k1"))
	assert.Nil(t, pkg.GetKeyedCounter("k1"))

	pkg.PutKeyeCounter(&gmetric.KeyedOperationCounter{
		KeyedOperationMetric: &gmetric.KeyedOperationMetric{
			Metrics: make(map[string]*gmetric.OperationMetric),
		},
		Name: "k1",
	})

	assert.NotNil(t, pkg.GetKeyedMetric("k1"))
	assert.NotNil(t, pkg.GetKeyedCounter("k1"))

	assert.Nil(t, pkg.GetMetric("k2"))
	assert.Nil(t, pkg.GetCounter("k2"))

	pkg.PutCounter(&gmetric.OperationCounter{
		OperationMetric: &gmetric.OperationMetric{
			Name: "k2",
		},
	})

	assert.NotNil(t, pkg.GetMetric("k2"))
	assert.NotNil(t, pkg.GetCounter("k2"))

}
