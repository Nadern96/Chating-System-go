package server

import (
	"context"
	"errors"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-chat/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type ChatServer struct {
	ctx            ctx.ServiceContext
	messageService *service.MessageService
}

func NewChatServer(serviceContext ctx.ServiceContext) *ChatServer {
	messageService := service.NewMessageService(serviceContext)
	return &ChatServer{
		ctx:            serviceContext,
		messageService: messageService,
	}
}

var ErrUnauthorized = errors.New(model.Unauthorized)

func (s *ChatServer) GrpcLogger(c context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	s.ctx.Logger().Debug(info.FullMethod)
	return handler(c, req)
}

func (s *ChatServer) SendMessage(c context.Context, req *proto.Message) (*proto.SendMessageResponse, error) {
	op := "ChatServer.SendMessage"

	if ok, _ := isAuthorized(c); !ok {
		s.ctx.Logger().Error(op + "." + model.Unauthorized)
		return nil, ErrUnauthorized
	}

	msg, err := model.MessageToModel(req)
	if err != nil {
		s.ctx.Logger().Error(op+".err: ", err)
		return nil, err
	}

	err = s.messageService.Send(c, msg)
	if err != nil {
		s.ctx.Logger().Error(op+".messageService.Send.err: ", err)
		return nil, err
	}

	return &proto.SendMessageResponse{
		Message: model.Success,
	}, nil
}

func (s *ChatServer) GetUserChats(c context.Context, req *proto.GetUserChatsRequest) (*proto.GetUserChatsResponse, error) {
	op := "ChatServer.GetUserChats"

	ok, userId := isAuthorized(c)
	if !ok {
		s.ctx.Logger().Error(op + "." + model.Unauthorized)
		return nil, ErrUnauthorized
	}

	chats, err := s.messageService.GetUserChats(c, userId)
	if err != nil {
		s.ctx.Logger().Error(op+".messageService.GetUserChats.err: ", err)
		return nil, err
	}

	res := &proto.GetUserChatsResponse{}
	for _, chat := range chats {
		res.Chats = append(res.Chats, chat.ToProto())
	}

	return res, nil
}

func isAuthorized(ctx context.Context) (bool, string) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		userId := md.Get("USER_ID")
		if len(userId) > 0 && userId[0] != "" {
			return true, userId[0]
		}
	}

	return false, ""
}
