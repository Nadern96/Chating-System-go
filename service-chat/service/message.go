package service

import (
	"context"

	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/model"
)

type MessageService struct {
	ctx ctx.ServiceContext
}

func NewMessageService(ctx ctx.ServiceContext) *MessageService {
	return &MessageService{
		ctx: ctx,
	}
}

func (s *MessageService) Send(ctx context.Context, msg model.Message) error {
	query := `INSERT INTO message (chatid, messageid, fromUserId, toUserId, content, createdAt) VALUES (?, ?, ?, ?, ?, ?)`

	err := s.ctx.GetCassandra().Query(query, msg.ChatID, msg.MessageID, msg.FromUserID, msg.ToUserID, msg.Content, msg.CreatedAt).WithContext(ctx).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) GetUserChats(ctx context.Context, userId string) ([]model.Chat, error) {
	userIdUUID, err := gocql.ParseUUID(userId)
	if err != nil {
		return nil, err
	}

	var chats []model.Chat
	iter := s.ctx.GetCassandra().Query(`SELECT chatid, fromuserid, touserid FROM chat WHERE fromuserid = ?`, userIdUUID).WithContext(ctx).Iter()

	for {
		chat := model.Chat{}
		if !iter.Scan(&chat.ChatID, &chat.FromUserID, &chat.ToUserID) {
			break
		}
		chats = append(chats, chat)
	}

	iter = s.ctx.GetCassandra().Query(`SELECT chatid, fromuserid, touserid FROM chat WHERE touserid = ?`, userIdUUID).WithContext(ctx).Iter()
	for {
		chat := model.Chat{}
		if !iter.Scan(&chat.ChatID, &chat.FromUserID, &chat.ToUserID) {
			break
		}
		chats = append(chats, chat)
	}

	if err := iter.Close(); err != nil {
		return nil, err
	}

	return chats, nil
}
