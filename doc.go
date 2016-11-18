package gmetric

/*

Package gmetric - Operation metric for go

This library comes with operational metric counters to meassure how application perform. Gmetric service exposes the counter via grpc or Rest endpoint.


It can be used to measure various aspects of the application, for instance execution time of  methods,  size of processed data, etc.



Usage:

    import (
       	"github.com/viant/gmetric"
    )

	var grpcPort, restPort = (8876, 8877)
	server, err := gmetric.NewServer(grpcPort, restPor)


	//register individual operation metrics counters
	someFuncLatency := server.Service().RegisterCounter("com/viant/app1", "someFuncLatency", "ns", "Time taken by some func in ns.", 10, nil)
	dataSizeProcessedByOtherFunc := server.Service().RegisterCounter("com/viant/app1, "otherFuncDataSize", "ns", ""Data size processed by otherFunc in bytes", 10, nil)



	func someFunction() (err error) {
		defer func(startTime time.Time) {
			someFuncLatency.AddLatency(startTime, err)
		}(time.Now())

		<<business logic comes herer>>
	}


	func otherFunction(payload []byte) (err error)  {
		someFuncLatency.Add(len(payload), err)
		<<business logic comes herer>>
	}




	//Dynamic Key Discovery use case
	//register individual keyed operation metrics counters

	someFuncLatencyByType := server.Service().RegisterKeyCounter("com/vinat/app1", "LatencyByType", "ns", "some desc", 10, nil, nil)
	dataSizeProcessedByOtherFuncByType := server.Service().RegisterKeyCounter("com/viant/app1, "otherFuncDataSizeByType", "ns", ""Data size processed by otherFunc in bytes", 10, nil, nil)




	func someFunction(type sting) (err error) {
		defer func(startTime time.Time) {
			someFuncLatencyByType.AddLatency(type, startTime, err)
		}(time.Now())

		<<business logic comes herer>>
	}


	func otherFunction(type string, payload []byte) (err error)  {
		someFuncLatencyByType.Add(type, len(payload), err)
		<<business logic comes herer>>
	}



*/
