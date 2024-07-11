package server

import (
	"context"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-auth/service"
	"github.com/nadern96/Chating-System-go/utils"
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

	err := s.authService.Register(c, model.User{
		Username: req.UserName,
		Password: req.Password,
		Email:    req.Email,
	})
	if err != nil {
		s.ctx.Logger().Error(op+".authService.Register err: ", err)
		return nil, err
	}

	return &proto.RegisterResponse{
		Message: "success",
	}, nil
}

func (s *AuthServer) Login(c context.Context, req *proto.LoginRequest) (*proto.LoginResponse, error) {
	op := "AuthServer.Login"

	token, err := s.authService.Login(c, model.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		s.ctx.Logger().Error(op+".authService.Login err: ", err)
		return nil, err
	}

	return &proto.LoginResponse{Token: token}, nil
}

func (s *AuthServer) Verify(c context.Context, req *proto.VerifyRequest) (*proto.VerifyResponse, error) {
	op := "AuthServer.Verify"

	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		s.ctx.Logger().Error(op+".authService.Verify.ParseToken err: ", err)
		return nil, err
	}

	email := claims.StandardClaims.Subject

	redisClient := s.ctx.Redis()
	redisVal, err := redisClient.Get(email).Result()
	if err != nil {
		s.ctx.Logger().Error(op+".authService.Verify.redisClient.Get err: ", err)
		return nil, model.ErrUnauthorized
	}

	return &proto.VerifyResponse{Message: redisVal}, nil
}
