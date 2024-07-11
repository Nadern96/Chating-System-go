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
	exist, err := s.isChatExist(ctx, msg.ChatID)
	if err != nil {
		return err
	}

	if !exist {
		return model.ErrChatIdNotExist
	}

	query := `INSERT INTO message (chatid, messageid, fromUserId, toUserId, content, createdAt) VALUES (?, ?, ?, ?, ?, ?)`
	err = s.ctx.GetCassandra().Query(query, msg.ChatID, msg.MessageID, msg.FromUserID, msg.ToUserID, msg.Content, msg.CreatedAt).WithContext(ctx).Exec()
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

func (s *MessageService) StartChat(ctx context.Context, chat model.Chat) error {
	query := `INSERT INTO chat (chatid, fromUserId, toUserId) VALUES (?, ?, ?)`

	err := s.ctx.GetCassandra().Query(query, chat.ChatID, chat.FromUserID, chat.ToUserID).WithContext(ctx).Exec()
	if err != nil {
		return err
	}

	return nil
}

func (s *MessageService) isChatExist(ctx context.Context, chatId gocql.UUID) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM chat WHERE chatid = ?`

	if err := s.ctx.GetCassandra().Query(query, chatId).WithContext(ctx).Scan(&count); err != nil {
		if err == gocql.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

func (s *MessageService) GetChatMessages(ctx context.Context, chatId, startMsgId string) ([]model.Message, error) {
	op := "MessageService.GetChatMessages"
	chatIdUUID, err := gocql.ParseUUID(chatId)
	if err != nil {
		return nil, err
	}

	baseQuery := "SELECT chatId, messageId, fromUserId, toUserId, content, createdAt FROM message WHERE chatId = ?"
	var query string
	var iter *gocql.Iter

	if startMsgId != "" {
		msgUUID, err := gocql.ParseUUID(startMsgId)
		if err != nil {
			return nil, err
		}

		query = baseQuery + " AND messageId <= ? ORDER BY messageId DESC LIMIT 10"
		iter = s.ctx.GetCassandra().Query(query, chatIdUUID, msgUUID).WithContext(ctx).Iter()
	} else {
		// Construct query without messageId condition
		query = baseQuery + " ORDER BY messageId DESC LIMIT 10"
		iter = s.ctx.GetCassandra().Query(query, chatIdUUID).WithContext(ctx).Iter()
	}

	var messages []model.Message
	for {
		msg := model.Message{}
		if !iter.Scan(&msg.ChatID, &msg.MessageID, &msg.FromUserID, &msg.ToUserID, &msg.Content, &msg.CreatedAt) {
			break
		}
		messages = append(messages, msg)
	}

	if err := iter.Close(); err != nil {
		s.ctx.Logger().Errorln(op+".iter err: ", err)
		return nil, err
	}

	return messages, nil
}
