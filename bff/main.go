package main

import (
	"github.com/gin-gonic/gin"
	"github.com/nadern96/Chating-System-go/bff/routes"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/env"
)

func main() {
	serviceContext := ctx.NewDefaultServiceContext()

	defer serviceContext.Shutdown()

	serviceContext.Logger().Info("Service bff")
	httpServer(serviceContext)
}

func httpServer(c *ctx.DefaultServiceContext) {
	engine := gin.Default()

	authRoute := routes.NewAuthRouter(c)
	authRoute.Install(engine.Group("/auth/"))

	chatRoute := routes.NewChatRouter(c)
	chatRoute.Install(engine.Group("/chat/"))

	c.ListenHTTP(env.GetEnvValue("PORT"), engine)
}
