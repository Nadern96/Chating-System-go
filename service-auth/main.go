package main

import (
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/env"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-auth/server"
)

func main() {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra()
	defer serviceContext.Shutdown()
	grpcServer(serviceContext)
	serviceContext.Logger().Info("Service Auth")

	var query string = "INSERT INTO auth.user(id,username,created_at) VALUES(?,?,?)"

	if err := serviceContext.GetCassandra().Query(query, "1", "nader", time.Now()).Exec(); err != nil {
		log.Println("err: ", err)
		return
	}

	var createdAt time.Time
	var username string
	var id string

	serviceContext.GetCassandra().Query("select * from auth.user").Scan(&id, &createdAt, &username)
	log.Println("id : ", id, " username: ", username, "createdAt: ", createdAt)
}

func grpcServer(c *ctx.DefaultServiceContext) {
	createdService := server.NewAuthServer(c)
	c.ListenGRPC(env.GetEnvValue("GRPC_PORT"), func(s *grpc.Server) {
		proto.RegisterAuthServer(s, createdService)
	}, grpc.UnaryInterceptor(createdService.GrpcLogger))
}
