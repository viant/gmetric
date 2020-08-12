package gmetric

import (
	"fmt"
	"github.com/viant/toolbox"
)

func Router(URI string, service *Service) *toolbox.ServiceRouter {
	return toolbox.NewServiceRouter(
		toolbox.ServiceRouting{
			HTTPMethod: "GET",
			URI:        fmt.Sprintf("%voperations", URI),
			Handler:    service.OperationCounters,
			Parameters: []string{},
		},
		toolbox.ServiceRouting{
			HTTPMethod: "GET",
			URI:        fmt.Sprintf("%voperation/{name}", URI),
			Handler:    service.LookupOperation,
			Parameters: []string{"name"},
		},
		toolbox.ServiceRouting{
			HTTPMethod: "GET",
			URI:        fmt.Sprintf("%vcounters", URI),
			Handler:    service.Counters,
			Parameters: []string{},
		},
		toolbox.ServiceRouting{
			HTTPMethod: "GET",
			URI:        fmt.Sprintf("%vcounter/{name}", URI),
			Handler:    service.LookupCounter,
			Parameters: []string{"name"},
		},
		toolbox.ServiceRouting{
			HTTPMethod: "GET",
			URI:        fmt.Sprintf("%vcounter/{name}/cumulative/{metric}", URI),
			Handler:    service.LookupOperationCumulativeMetric,
			Parameters: []string{"name", "metric"},
		},
		toolbox.ServiceRouting{
			HTTPMethod: "GET",
			URI:        fmt.Sprintf("%vcounter/{name}/recent/{metric}", URI),
			Handler:    service.LookupOperationRecentMetric,
			Parameters: []string{"name", "metric"},
		},
	)
}
