package routes

import (
	"net/http"
	"strings"

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
	engine.GET("/logout", r.Logout)
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

func (r *AuthRouter) Logout(ginCtx *gin.Context) {
	op := "authRouter.Logout"

	if len(ginCtx.Request.Header["Authorization"]) == 0 {
		r.serviceContext.Logger().Error("Invalid Headers, unauthorized")
		ginCtx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "UNAUTHORIZED"})
		return
	}

	token := strings.Split(ginCtx.Request.Header["Authorization"][0], " ")[1]
	res, err := r.authClient.Logout(ginCtx.Request.Context(), &proto.LogoutRequest{Token: token})
	if err != nil {
		r.serviceContext.Logger().Error(op+".authClient.Logout err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}
