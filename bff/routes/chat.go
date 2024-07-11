package routes

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/grpcclient"
	"github.com/nadern96/Chating-System-go/proto"
	"google.golang.org/grpc/metadata"
)

type ChatRouter struct {
	serviceContext ctx.ServiceContext
	authClient     proto.AuthClient
	chatClient     proto.ChatClient
}

func NewChatRouter(serviceContext ctx.ServiceContext) *ChatRouter {
	authClient := grpcclient.NewClientAuth(serviceContext)
	chatClient := grpcclient.NewClientChat(serviceContext)

	return &ChatRouter{
		serviceContext: serviceContext,
		authClient:     authClient.Client(),
		chatClient:     chatClient.Client(),
	}
}

func (r *ChatRouter) Install(engine *gin.RouterGroup) {
	engine.POST("/send", r.AuthVerify(), r.Send)
	engine.GET("", r.AuthVerify(), r.GetUserChats)
	engine.POST("/start", r.AuthVerify(), r.StartChat)
	engine.GET("/:chatId", r.AuthVerify(), r.GetChatMessages)
}

func (r *ChatRouter) AuthVerify() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("authorization = ", c.Request.Header["Authorization"])

		if len(c.Request.Header["Authorization"]) == 0 {
			r.serviceContext.Logger().Error("Invalid Headers, unauthorized")
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "UNAUTHORIZED"})
			return
		}

		token := strings.Split(c.Request.Header["Authorization"][0], " ")[1]

		res, err := r.authClient.Verify(c, &proto.VerifyRequest{Token: token})
		if err != nil {
			r.serviceContext.Logger().Error("AuthVerify err: ", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.Request.Header.Set("USER_ID", res.Message)
		c.Next()
	}
}

func (r *ChatRouter) Send(ginCtx *gin.Context) {
	op := "chatRouter.Send"

	req := &proto.Message{}
	err := ginCtx.BindJSON(req)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCtx := metadata.AppendToOutgoingContext(ginCtx, "USER_ID", ginCtx.Request.Header.Get("USER_ID"))
	res, err := r.chatClient.SendMessage(newCtx, req)
	if err != nil {
		r.serviceContext.Logger().Error(op+".chatClient.SendMessage err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}

func (r *ChatRouter) GetUserChats(ginCtx *gin.Context) {
	op := "chatRouter.GetUserChats"

	newCtx := metadata.AppendToOutgoingContext(ginCtx, "USER_ID", ginCtx.Request.Header.Get("USER_ID"))
	res, err := r.chatClient.GetUserChats(newCtx, &proto.GetUserChatsRequest{})

	if err != nil {
		r.serviceContext.Logger().Error(op+".chatClient.GetUserChats err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}

func (r *ChatRouter) StartChat(ginCtx *gin.Context) {
	op := "chatRouter.StartChat"

	req := &proto.StartChatRequest{}
	err := ginCtx.BindJSON(req)
	if err != nil {
		ginCtx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCtx := metadata.AppendToOutgoingContext(ginCtx, "USER_ID", ginCtx.Request.Header.Get("USER_ID"))
	res, err := r.chatClient.StartChat(newCtx, req)
	if err != nil {
		r.serviceContext.Logger().Error(op+".chatClient.StartChat err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}

func (r *ChatRouter) GetChatMessages(ginCtx *gin.Context) {
	op := "chatRouter.GetChatMessages"

	chatId := ginCtx.Param("chatId")

	startMsgId := ginCtx.Query("startMsgId")

	newCtx := metadata.AppendToOutgoingContext(ginCtx, "USER_ID", ginCtx.Request.Header.Get("USER_ID"))

	res, err := r.chatClient.GetChatMessages(newCtx, &proto.GetChatMessageRequest{ChatId: chatId, StartMsgId: startMsgId})
	if err != nil {
		r.serviceContext.Logger().Error(op+".chatClient.GetChatMessages err: ", err)
		ginCtx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ginCtx.JSON(http.StatusOK, res)
}
