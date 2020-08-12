package gmetric

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func flagMod(mod int) func() bool {
	counter := 0
	return func() bool {
		counter++
		return counter%mod == 0
	}
}

func errorMod(mod int) func() error {
	counter := 0
	return func() error {
		counter++
		if counter%mod == 0 {
			return errors.New("ttest")
		}
		return nil
	}
}



func TestNewCache(t *testing.T) {

	var useCases = []struct {
		description  string
		count        int
		execDuration time.Duration
		flag         func() bool
		error        func() error
	}{
		{
			description:  "basic test",
			count:        100,
			execDuration: 10 * time.Millisecond,
			flag:         flagMod(2),
			error:        errorMod(2),
		},
	}

	for _, useCase := range useCases {

		counter := NewOperation(useCase.description, "flag", "", 15, time.Minute, time.Millisecond, nil)
		for i := 0; i < useCase.count; i++ {
			at := time.Now()
			onDone := counter.Begin(at)
			onDone(time.Now(), useCase.error)
		}
		assert.EqualValues(t, counter.Count, useCase.count, useCase.description)
	}

}
