package server

import (
	"context"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-auth/service"
	"google.golang.org/grpc"
)

type AuthServer struct {
	ctx         ctx.ServiceContext
	authService *service.AuthService
}

func NewAuthServer(serviceContext ctx.ServiceContext) *AuthServer {
	authService := service.NewAuthService(serviceContext)
	return &AuthServer{
		ctx:         serviceContext,
		authService: authService,
	}
}

func (s *AuthServer) GrpcLogger(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s.ctx.Logger().Debug(info.FullMethod)
	return handler(c, req)
}

func (s *AuthServer) Register(c context.Context, req *proto.RegisterRequest) (*proto.RegisterResponse, error) {
	op := "AuthServer.Register"

	err := s.authService.Register(model.User{
		Username: req.UserName,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		s.ctx.Logger().Error(op+"authService.Register err: ", err)
		return nil, err
	}

	return &proto.RegisterResponse{
		Message: "success",
	}, nil
}
