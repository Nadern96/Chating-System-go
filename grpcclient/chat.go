package grpcclient

import (
	"log"
	"os"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/proto"
	"google.golang.org/grpc"
)

type ChatClient struct {
	ctx        ctx.ServiceContext
	connection *grpc.ClientConn
	client     proto.ChatClient
}

func NewClientChat(serviceContext ctx.ServiceContext) *ChatClient {
	connection, client := NewPlainChatClient(os.Getenv("SERVICE_CHAT_URL"))
	return &ChatClient{
		ctx:        serviceContext,
		connection: connection,
		client:     client,
	}
}

func NewPlainChatClient(link string) (*grpc.ClientConn, proto.ChatClient) {
	conn, err := grpc.Dial(link, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[%s] error [%s]", "grpcClient.auth.NewPlainChatClient.Dial", err.Error())
	}
	client := proto.NewChatClient(conn)
	return conn, client
}

func (cc *ChatClient) Client() proto.ChatClient {
	return cc.client
}

func (ac *ChatClient) Close() {
	err := ac.connection.Close()
	if err != nil {
		log.Fatalf("[%s] error [%s]", "grpcClient.chat.close.close", err.Error())
	}
}
