package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

}

func ListenGRPC(port string, registerFn func(*grpc.Server), opt ...grpc.ServerOption) {
	grpcServer := grpc.NewServer(opt...)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	registerFn(grpcServer)
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			log.Fatal(err)
		}
	}()
	return grpcServer
}
