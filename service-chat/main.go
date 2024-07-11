package main

import (
	"google.golang.org/grpc"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/env"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-chat/server"
)

func main() {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	defer serviceContext.Shutdown()
	grpcServer(serviceContext)
	serviceContext.Logger().Info("Service Chat")
}

func grpcServer(c *ctx.DefaultServiceContext) {
	createdService := server.NewChatServer(c)
	c.ListenGRPC(env.GetEnvValue("GRPC_PORT"), func(s *grpc.Server) {
		proto.RegisterChatServer(s, createdService)
	}, grpc.UnaryInterceptor(createdService.GrpcLogger))
}
