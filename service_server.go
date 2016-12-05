package gmetric

import (
	"golang.org/x/net/context"
)

type serviceServer struct {
	service CounterService
}

func (s *serviceServer) Query(context context.Context, request *QueryRequest) (response *QueryResponse, err error) {
	response = &QueryResponse{}
	metrics, err := s.service.Query(request.Query)
	if err != nil {
		response.Error = err.Error()
		return response, err
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
