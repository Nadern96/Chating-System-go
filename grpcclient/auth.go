package grpcclient

import (
	"log"
	"os"

	"github.com/nadern96/Chating-System-go/ctx"
	"github.com/nadern96/Chating-System-go/proto"
	"google.golang.org/grpc"
)

type AuthClient struct {
	ctx        ctx.ServiceContext
	connection *grpc.ClientConn
	client     proto.AuthClient
}

func NewClientAuth(serviceContext ctx.ServiceContext) *AuthClient {
	connection, client := NewPlainAuthClient(os.Getenv("SERVICE_AUTH_URL"))
	return &AuthClient{
		ctx:        serviceContext,
		connection: connection,
		client:     client,
	}
}

func NewPlainAuthClient(link string) (*grpc.ClientConn, proto.AuthClient) {
	conn, err := grpc.Dial(link, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("[%s] error [%s]", "grpcClient.auth.NewPlainAuthClient.Dial", err.Error())
	}
	client := proto.NewAuthClient(conn)
	return conn, client
}

func (ac *AuthClient) Client() proto.AuthClient {
	return ac.client
}

func (ac *AuthClient) Close() {
	err := ac.connection.Close()
	if err != nil {
		log.Fatalf("[%s] error [%s]", "grpcClient.auth.close.close", err.Error())
	}
}
