package grpcapi

import (
	"net"

	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/api/grpcspec"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/app/calendar"
	"github.com/f0m41h4u7/go-hw/hw12_13_14_15_calendar/internal/pkg/config"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var app *calendar.Calendar

//go:generate protoc --proto_path=../../../api/grpcspec --go_out=plugins=grpc:../../../api/grpcspec ../../../api/grpcspec/event.proto
type Server struct {
	Grpc *grpc.Server
}

func InitServer(cl *calendar.Calendar) *Server {
	app = cl
	s := &Server{}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(grpc_zap.UnaryServerInterceptor(zap.L())))
	grpcspec.RegisterCalendarServer(grpcServer, s)
	s.Grpc = grpcServer
	return s
}

func StartServer(cl *calendar.Calendar) {
	s := InitServer(cl)
	lis, err := net.Listen("tcp", net.JoinHostPort(config.Conf.Grpcserver.Host, config.Conf.Grpcserver.Port))
	if err != nil {
		zap.L().Error("failed to run grpc server", zap.Error(err))
		return
	}
	err = s.Grpc.Serve(lis)
	if err != nil {
		zap.L().Error("failed to run grpc server", zap.Error(err))
		return
	}
}
