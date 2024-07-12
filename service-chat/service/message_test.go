package service

import (
	"context"
	"errors"
	"testing"

	"slices"

	"github.com/bmizerany/assert"
	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
	"github.com/nadern96/Chating-System-go/proto"
)

func TestSendMessagetoANonExistentChat(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	messageService := NewMessageService(serviceContext)

	req := &proto.Message{}

	req.ChatId = gocql.TimeUUID().String()
	req.FromUserId = gocql.TimeUUID().String()
	req.ToUserId = gocql.TimeUUID().String()

	msg, err := model.MessageToModel(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = messageService.Send(context.Background(), msg)
	if err != nil && !errors.Is(err, model.ErrChatIdNotExist) {
		t.Errorf(err.Error())
	}

	assert.Equalf(t, errors.Is(err, model.ErrChatIdNotExist), true, "expected err: chat id does not exist")
}

func TestSendMessage(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	messageService := NewMessageService(serviceContext)

	chat := model.Chat{
		ChatID:     gocql.TimeUUID(),
		ToUserID:   gocql.TimeUUID(), // assume registering a user and use his id here
		FromUserID: gocql.TimeUUID(), // assume registering a user and use his id here
	}

	err := messageService.StartChat(context.Background(), chat)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := &proto.Message{}

	req.ChatId = chat.ChatID.String()
	req.FromUserId = chat.FromUserID.String()
	req.ToUserId = chat.ToUserID.String()
	req.Content = "Hi"

	msg, err := model.MessageToModel(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = messageService.Send(context.Background(), msg)
	if err != nil {
		t.Errorf(err.Error())
	}

	assert.Equalf(t, errors.Is(err, nil), true, "expected err: nil")
}

func TestGetChatMessages(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	messageService := NewMessageService(serviceContext)

	chat := model.Chat{
		ChatID:     gocql.TimeUUID(),
		ToUserID:   gocql.TimeUUID(), // assume registering a user and use his id here
		FromUserID: gocql.TimeUUID(), // assume registering a user and use his id here
	}

	err := messageService.StartChat(context.Background(), chat)
	if err != nil {
		t.Errorf(err.Error())
	}

	req := &proto.Message{}

	req.ChatId = chat.ChatID.String()
	req.FromUserId = chat.FromUserID.String()
	req.ToUserId = chat.ToUserID.String()
	req.Content = "Hi"

	msg, err := model.MessageToModel(req)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = messageService.Send(context.Background(), msg)
	if err != nil {
		t.Errorf(err.Error())
	}

	messages, err := messageService.GetChatMessages(context.Background(), req.ChatId, "")
	if err != nil {
		t.Errorf(err.Error())
	}

	messagesContent := []string{}
	for _, message := range messages {
		// Hi
		messagesContent = append(messagesContent, message.Content)
	}

	assert.Equalf(t, slices.Contains(messagesContent, "Hi"), true, "expected messages content to contain Hi")
}

func TestStartChatThatAlreadyExist(t *testing.T) {
	serviceContext := ctx.NewDefaultServiceContext().WithCassandra().WithRedis()
	messageService := NewMessageService(serviceContext)

	chat := model.Chat{
		ChatID:     gocql.TimeUUID(),
		ToUserID:   gocql.TimeUUID(), // assume registering a user and use his id here
		FromUserID: gocql.TimeUUID(), // assume registering a user and use his id here
	}

	err := messageService.StartChat(context.Background(), chat)
	if err != nil {
		t.Errorf(err.Error())
	}

	err = messageService.StartChat(context.Background(), chat)
	if err != nil && !errors.Is(err, model.ErrChatAlreadyExists) {
		t.Errorf(err.Error())
	}

	assert.Equalf(t, errors.Is(err, model.ErrChatAlreadyExists), true, "expected chat already exist")
}
