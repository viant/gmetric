package window

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/viant/gmetric/counter"
	"github.com/viant/gmetric/counter/base"
	"github.com/viant/gmetric/stat"
	"testing"
	"time"
)

func TestOperation_Begin(t *testing.T) {
	var useCases = []struct {
		description string
		recentBuckets int
		recentUnit time.Duration
		unit        time.Duration
		counter.Provider
		stats []interface{}
		unknownStats []interface{}
	}{
		{
			description:"single metrics",
			unit:time.Microsecond,
		},
		{
			description:"stats metrics",
			unit:time.Microsecond,
			Provider:base.NewProvider("key1", "key2"),
			stats:[]interface{}{"key1", "key2", errors.New("test")},
		},
		{
			description:"stats metrics with unknown value",
			unit:time.Microsecond,
			Provider:base.NewProvider("key1", "key2"),
			stats:[]interface{}{"key1", "key2", errors.New("test")},
			unknownStats:[]interface{}{"key3"},
		},
	}

	for _, useCase := range useCases {
		op := NewOperation(useCase.recentBuckets, useCase.recentUnit, useCase.unit, useCase.Provider)

		for i := 0; i < 100; i++ {
			statsValues := stat.New()
			onDone := op.Begin(time.Now())

			if len(useCase.stats) > 0 {
				statsValues.AppendAll(useCase.stats)
			}
			if len(useCase.unknownStats) > 0 {
				statsValues.AppendAll(useCase.unknownStats)
			}

			onDone(time.Now(), statsValues)
		}
		assert.EqualValues(t, 100, op.CountValue(), useCase.description)
		if len(useCase.stats) == 0 {
			continue
		}
		for _, val := range useCase.stats {
			index := useCase.Provider.Map(val)
			if index == -1 {
				continue
			}
			assert.EqualValues(t, 100, op.Counters[index].CountValue(), useCase.description)
		}
	}
}
