package server

import (
	"context"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/proto"
	"google.golang.org/grpc"
)

type AuthServer struct {
	ctx ctx.ServiceContext
}

func NewAuthServer(serviceContext ctx.ServiceContext) *AuthServer {
	return &AuthServer{
		ctx: serviceContext,
	}
}

func (s *AuthServer) GrpcLogger(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s.ctx.Logger().Debug(info.FullMethod)
	return handler(c, req)
}

func (s *AuthServer) Register(c context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	return &proto.RegisterResponse{
		Message: "success",
	}, nil
}
