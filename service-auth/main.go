package main

import (
	"google.golang.org/grpc"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/env"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-auth/server"
)

func main() {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	defer serviceContext.Shutdown()
	grpcServer(serviceContext)
	serviceContext.Logger().Info("Service Auth")
}

func grpcServer(c *ctx.DefaultServiceContext) {
	createdService := server.NewAuthServer(c)
	c.ListenGRPC(env.GetEnvValue("GRPC_PORT"), func(s *grpc.Server) {
		proto.RegisterAuthServer(s, createdService)
	}, grpc.UnaryInterceptor(createdService.GrpcLogger))
}
