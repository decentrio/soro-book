package server

import (
	"context"
	"net"

	"github.com/decentrio/soro-book/lib/service"
	apitypes "github.com/decentrio/soro-book/types/api"
	"google.golang.org/grpc"
)

type GRPCServer struct {
	service.BaseService

	addr     string
	listener net.Listener
	server   *grpc.Server

	app *Application
}

func NewServer(addr string, app *Application) *GRPCServer {
	s := &GRPCServer{
		addr:     addr,
		listener: nil,
		app:      app,
	}
	s.BaseService = *service.NewBaseService("GRPC-Server", s)
	return s
}

// OnStart starts the gRPC service.
func (s *GRPCServer) OnStart() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.listener = ln
	s.server = grpc.NewServer()
	apitypes.RegisterAPIServicesServer(s.server, &gprcApplication{})

	go s.server.Serve(s.listener)
	return nil
}

// OnStop stops the gRPC server.
func (s *GRPCServer) OnStop() error {
	s.server.Stop()
	return nil
}

type Application interface {
}

type gprcApplication struct {
	Application
}

func (a *gprcApplication) Start(context.Context, *apitypes.StartRequest) (*apitypes.StartResponse, error) {
	return &apitypes.StartResponse{}, nil
}

func (a *gprcApplication) Stop(context.Context, *apitypes.StopRequest) (*apitypes.StopResponse, error) {
	return &apitypes.StopResponse{}, nil
}
