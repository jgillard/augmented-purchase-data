package main

import (
	"fmt"
	"log"
	"net"

	api "github.com/jgillard/practising-go-tdd/rpc"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 7777))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := api.Server{}

	grpcServer := grpc.NewServer()

	api.RegisterStatusServer(grpcServer, &s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
