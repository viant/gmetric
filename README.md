# Application operation performance metric (gmetric)


[![Application operation performance metric.](https://goreportcard.com/badge/github.com/viant/gmetric)](https://goreportcard.com/report/github.com/viant/gmetric)
[![GoDoc](https://godoc.org/github.com/viant/asc?status.svg)](https://godoc.org/github.com/viant/gmetric)

This library is compatible with Go 1.5+

Please refer to [`CHANGELOG.md`](CHANGELOG.md) if you encounter breaking changes.

- [Usage](#Usage)
- [License](#License)
- [Credits and Acknowledgements](#Credits-and-Acknowledgements)



## Usage:


This library comes with operational metric counters to measure how application perform. Gmetric service exposes the counter via grpc or Rest endpoint.
It can be used to measure various aspects of the application, for instance execution time of  methods,  size of processed data, etc.

```go


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
```


<a name="License"></a>
## License

The source code is made available under the terms of the Apache License, Version 2, as stated in the file `LICENSE`.

Individual files may be made available under their own specific license,
all compatible with Apache License, Version 2. Please see individual files for details.


<a name="Credits-and-Acknowledgements"></a>

##  Credits and Acknowledgements

**Library Author:** Adrian Witas

**Contributors:**