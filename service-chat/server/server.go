package server

import (
	"context"
	"errors"

	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/proto"
	"github.com/nadern96/Chating-System-go/service-chat/service"
	"github.com/nadern96/Chating-System-go/utils"

	"google.golang.org/grpc"
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

	ok, fromUserId := utils.IsAuthorized(c)
	if !ok {
		s.ctx.Logger().Error(op + "." + model.Unauthorized)
		return nil, ErrUnauthorized
	}

	req.FromUserId = fromUserId

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

	ok, userId := utils.IsAuthorized(c)
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

func (s *ChatServer) StartChat(c context.Context, req *proto.StartChatRequest) (*proto.StartChatResponse, error) {
	op := "ChatServer.StartChat"

	ok, userId := utils.IsAuthorized(c)
	if !ok {
		s.ctx.Logger().Error(op + "." + model.Unauthorized)
		return nil, ErrUnauthorized
	}

	fromUserId, err := gocql.ParseUUID(userId)
	if err != nil {
		s.ctx.Logger().Error(op+".ParseUUID.fromUserId.err: ", err)
		return nil, err
	}

	toUserId, err := gocql.ParseUUID(req.ToUserId)
	if err != nil {
		s.ctx.Logger().Error(op+".ParseUUID.toUserId.err: ", err)
		return nil, err
	}

	chat := model.Chat{
		ChatID:     gocql.TimeUUID(),
		ToUserID:   fromUserId,
		FromUserID: toUserId,
	}

	err = s.messageService.StartChat(c, chat)
	if err != nil {
		s.ctx.Logger().Error(op+".messageService.StartChat.err: ", err)
		return nil, err
	}

	return &proto.StartChatResponse{
		ChatId: chat.ChatID.String(),
	}, nil
}

func (s *ChatServer) GetChatMessages(c context.Context, req *proto.GetChatMessageRequest) (*proto.GetChatMessageResponse, error) {
	op := "ChatServer.StartChat"

	ok, _ := utils.IsAuthorized(c)
	if !ok {
		s.ctx.Logger().Error(op + "." + model.Unauthorized)
		return nil, ErrUnauthorized
	}

	messages, err := s.messageService.GetChatMessages(c, req.ChatId, req.StartMsgId)
	if err != nil {
		s.ctx.Logger().Error(op+".messageService.GetChatMessages.err: ", err)
		return nil, err
	}

	res := &proto.GetChatMessageResponse{}
	for _, msg := range messages {
		res.Messages = append(res.Messages, msg.ToProto())
	}

	return res, nil
}
