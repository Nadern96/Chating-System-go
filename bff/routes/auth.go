package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/grpcclient"
	"github.com/nadern96/Chating-System-go/proto"
)

type AuthRouter struct {
	serviceContext ctx.ServiceContext
	authClient     proto.AuthClient
}

func NewAuthRouter(serviceContext ctx.ServiceContext) *AuthRouter {
	authClient := grpcclient.NewClientAuth(serviceContext)

	return &AuthRouter{
		serviceContext: serviceContext,
		authClient:     authClient.Client(),
	}
}

func (r *AuthRouter) Install(engine *gin.RouterGroup) {
	engine.POST("/login", r.Login)
	engine.POST("/register", r.Register)
	// r.GET("/logout", controllers.Logout)
}

func (r *AuthRouter) Register(ginCtx *gin.Context) {
	op := "authRouter.Register"

	req := &proto.RegisterRequest{}

	err := ginCtx.BindJSON(req)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.authClient.Register(ginCtx.Request.Context(), req)
	if err != nil {
		r.serviceContext.Logger().Error(op+".authClient.Register err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}

func (r *AuthRouter) Login(ginCtx *gin.Context) {
	op := "authRouter.Login"

	req := &proto.LoginRequest{}

	if err := ginCtx.BindJSON(&req); err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := r.authClient.Login(ginCtx.Request.Context(), req)
	if err != nil {
		r.serviceContext.Logger().Error(op+".authClient.Login err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}
