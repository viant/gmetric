package gmetric

import (
	"github.com/viant/toolbox"
	"net/http"
)

type handler struct {
	*toolbox.ServiceRouter
}

//ServeHTTP serves http traffic
func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	err := h.Route(writer, request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}
}


//New creates a new metrics
func Handler(URI string, metrics *Service) http.Handler {
	return &handler{ServiceRouter: Router(URI, metrics)}
}
