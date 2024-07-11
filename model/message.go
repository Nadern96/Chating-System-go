package model

import (
	"time"

	"github.com/gocql/gocql"
	"github.com/nadern96/Chating-System-go/proto"
)

type Message struct {
	ChatID     gocql.UUID `json:"chatId"`
	MessageID  gocql.UUID `json:"messageId"`
	FromUserID gocql.UUID `json:"fromUserId"`
	ToUserID   gocql.UUID `json:"toUserId"`
	Content    string     `json:"content"`
	CreatedAt  time.Time  `json:"createdAt"`
}

type Chat struct {
	ChatID     gocql.UUID `json:"chatId"`
	FromUserID gocql.UUID `json:"fromUserId"`
	ToUserID   gocql.UUID `json:"toUserId"`
}

func MessageToModel(in *proto.Message) (Message, error) {
	from, err := gocql.ParseUUID(in.FromUserId)
	if err != nil {
		return Message{}, err
	}

	to, err := gocql.ParseUUID(in.ToUserId)
	if err != nil {
		return Message{}, err
	}

	chatId, err := gocql.ParseUUID(in.ChatId)
	if err != nil {
		return Message{}, err
	}

	return Message{
		ChatID:     chatId,
		MessageID:  gocql.TimeUUID(),
		ToUserID:   to,
		FromUserID: from,
		Content:    in.Content,
		CreatedAt:  time.Now(),
	}, nil
}

func (c Chat) ToProto() *proto.Chat {
	return &proto.Chat{
		ChatId:     c.ChatID.String(),
		ToUserId:   c.ToUserID.String(),
		FromUserId: c.FromUserID.String(),
	}
}
