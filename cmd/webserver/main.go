package main

import (
	"log"
	"net/http"
	"os"

	transactioncategories "github.com/jgillard/practising-go-tdd"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	categoryStore := transactioncategories.NewInMemoryCategoryStore()
	questionStore := transactioncategories.NewInMemoryQuestionStore()

	server := transactioncategories.NewServer(categoryStore, questionStore)

	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
