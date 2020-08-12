package gmetric_test

import (
	"errors"
	"github.com/viant/gmetric"
	"github.com/viant/gmetric/counter/base"
	"github.com/viant/gmetric/stat"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func ExampleService_Counter() {
	metrics := gmetric.New()
	handler := gmetric.NewHandler("/v1/metrics", metrics)
	http.Handle("/v1/metrics", handler)

	//basic single counter
	counter := metrics.OperationCounter("pkg.myapp", "mySingleCounter1", "my description", time.Microsecond, time.Minute, 2)
	go runBasicTasks(counter)

	go http.ListenAndServe(":8080", http.DefaultServeMux)
}

func runBasicTasks(counter *gmetric.Operation) {
	for i := 0; i < 1000; i++ {
		runBasicTask(counter)
	}
}

func runBasicTask(counter *gmetric.Operation) {
	onDone := counter.Begin(time.Now())
	defer func() {
		onDone(time.Now())
	}()
	time.Sleep(time.Nanosecond)

}

const (
	NoSuchKey = "noSuchKey"

	MyStatsCacheHit       = "cacheHit"
	MyStatsCacheCollision = "cacheCollision"
)

//MultiStateStatTestProvider represents multi stats value provider
type MultiStateStatTestProvider struct{ *base.Provider }

//Map maps value int slice index
func (p *MultiStateStatTestProvider) Map(key interface{}) int {
	textKey, ok := key.(string)
	if !ok {
		return p.Provider.Map(key)
	}
	switch textKey {
	case NoSuchKey:
		return 1
	case MyStatsCacheHit:
		return 2
	case MyStatsCacheCollision:
		return 3
	}
	return -1
}

func ExampleService_MultiOperationCounter() {
	metrics := gmetric.New()
	handler := gmetric.NewHandler("/v1/metrics", metrics)
	http.Handle("/v1/metrics", handler)
	counter := metrics.MultiOperationCounter("pkg.myapp", "myMultiCounter", "my description", time.Microsecond, time.Minute, 2, &MultiStateStatTestProvider{})
	go runMultiStateTasks(counter)

	err := http.ListenAndServe(":8080", http.DefaultServeMux)
	if err != nil {
		log.Fatal(err)
	}
}

func runMultiStateTasks(counter *gmetric.Operation) {
	for i := 0; i < 1000; i++ {
		runMultiStateTask(counter)
	}

}

func runMultiStateTask(counter *gmetric.Operation) {
	stats := stat.New()
	onDone := counter.Begin(time.Now())
	defer func() {
		onDone(time.Now(), stats)
	}()

	time.Sleep(time.Nanosecond)
	//simulate task metrics state
	rnd := rand.NewSource(time.Now().UnixNano())
	state := rnd.Int63() % 3
	switch state {
	case 0:
		stats.Append(NoSuchKey)
	case 1:
		stats.Append(MyStatsCacheHit)
	case 2:
		stats.Append(MyStatsCacheHit)
		stats.Append(MyStatsCacheCollision)
	}
	if rnd.Int63()%10 == 0 {
		stats.Append(errors.New("test error"))
	}
}
