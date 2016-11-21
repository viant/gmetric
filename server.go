package gmetric

import (
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

//Server represents operation counter server.
type Server struct {
	grpcPort        int
	restPort        int
	service         CounterService
	serverService   ServiceServer
	server          *grpc.Server
	isRunning       int32
	netListener     net.Listener
	restNetListener net.Listener
	cancelFunc      context.CancelFunc
}

//Service returns a operation counter service
func (s *Server) Service() CounterService {
	return s.service
}

//Start starts server.
func (s *Server) Start() (err error) {

	s.server = grpc.NewServer()
	RegisterServiceServer(s.server, s.serverService)
	atomic.StoreInt32(&s.isRunning, 1)
	s.netListener, err = net.Listen("tcp", ":"+strconv.Itoa(s.grpcPort))
	if err != nil {
		return fmt.Errorf("failed to grcp listen: %v", err)
	}

	go func() {

		fmt.Printf("%v(%v) GRCP started listening on %v\n", ApplicationName, ApplicationVersion, s.grpcPort)
		if err = s.server.Serve(s.netListener); err != nil {
			isRunning := atomic.LoadInt32(&s.isRunning) == 1
			if isRunning {
				log.Fatalf("failed to listen: %v", err)
			}
		}
	}()

	s.restNetListener, err = net.Listen("tcp", ":"+strconv.Itoa(s.restPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go func() {
		ctx := context.Background()
		ctx, s.cancelFunc = context.WithCancel(ctx)

		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		endpoint := "localhost:" + strconv.Itoa(s.grpcPort)
		err := RegisterServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
		if err != nil {
			log.Fatalf("Failed to connection to grpc endpoint: %v", err)
		}

		fmt.Printf("%v(%v) REST started listening on %v\n", ApplicationName, ApplicationVersion, s.restPort)
		err = http.Serve(s.restNetListener, mux)
		if err != nil {
			log.Fatal(err)
		}
	}()

	return err

}

//Stop stops server
func (s *Server) Stop() (err error) {
	fmt.Printf("%v(%v) stopping\n", ApplicationName, ApplicationVersion)
	atomic.StoreInt32(&s.isRunning, 0)
	s.cancelFunc()

	if s.cancelFunc != nil {
		s.cancelFunc()
	}
	time.Sleep(100 * time.Millisecond)
	if s.netListener != nil {
		err = s.netListener.Close()
		s.server.Stop()
	}
	return err
}

//NewServer create a new server with specified grpc and rest endpoint.
func NewServer(grpcPort, restPort int) (*Server, error) {
	service := NewCounterService()
	return &Server{
		grpcPort:      grpcPort,
		restPort:      restPort,
		service:       service,
		serverService: mewServiceServer(service),
	}, nil
}
