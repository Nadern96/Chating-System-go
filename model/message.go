package model

import (
	"log"
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
		log.Println("err from user. ", err)
		return Message{}, err
	}

	to, err := gocql.ParseUUID(in.ToUserId)
	if err != nil {
		log.Println("err touser. ", err)
		return Message{}, err
	}

	chatId, err := gocql.ParseUUID(in.ChatId)
	if err != nil {
		log.Println("err chat id . ", err)
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

func (m Message) ToProto() *proto.Message {
	return &proto.Message{
		ChatId:     m.ChatID.String(),
		ToUserId:   m.ToUserID.String(),
		FromUserId: m.FromUserID.String(),
		Content:    m.Content,
		CreatedAt:  m.CreatedAt.Format("2006-01-02 15:04:05"),
		MessageId:  m.MessageID.String(),
	}
}

func (c Chat) ToProto() *proto.Chat {
	return &proto.Chat{
		ChatId:     c.ChatID.String(),
		ToUserId:   c.ToUserID.String(),
		FromUserId: c.FromUserID.String(),
	}
}
