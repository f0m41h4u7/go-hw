package grpcapi

import (
	"net"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/api/grpcspec"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/config"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var app *calendar.Calendar

//go:generate protoc --proto_path=../../../api/grpcspec --go_out=plugins=grpc:../../../api/grpcspec ../../../api/grpcspec/event.proto
type Server struct {
	grpc *grpc.Server
}

func InitServer(cl *calendar.Calendar) *Server {
	app = cl
	s := &Server{}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zap.L())))
	grpcspec.RegisterCalendarServer(grpcServer, s)
	s.grpc = grpcServer

	return s
}

func (s *Server) Start() error {
	lis, err := net.Listen("tcp", net.JoinHostPort(config.Conf.GRPCServer.Host, config.Conf.GRPCServer.Port))
	if err != nil {
		return err
	}
	err = s.grpc.Serve(lis)

	return err
}

func (s *Server) Stop() {
	s.grpc.GracefulStop()
}
