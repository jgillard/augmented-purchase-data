package main

import (
	"log"
	"net"
	"os"

	"github.com/jgillard/practising-go-tdd/internal"
	api "github.com/jgillard/practising-go-tdd/rpc"

	"google.golang.org/grpc"
)

func main() {

	port := os.Getenv("RPC_PORT")
	if port == "" {
		log.Fatal("$RPC_PORT must be set")
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	categoryStore := internal.NewInMemoryCategoryStore(internal.CategoryList{})
	questionStore := internal.NewInMemoryQuestionStore(internal.QuestionList{})

	s := api.NewServer(categoryStore, questionStore)
	grpcServer := grpc.NewServer()
	api.RegisterPracticingGoTddServer(grpcServer, s)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
