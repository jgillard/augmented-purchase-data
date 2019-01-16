package main

import (
	"log"
	"net/http"
	"os"

	httptransport "github.com/jgillard/practising-go-tdd/http"
	internal "github.com/jgillard/practising-go-tdd/internal"
)

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	categoryStore := internal.NewInMemoryCategoryStore(internal.CategoryList{})
	questionStore := internal.NewInMemoryQuestionStore(internal.QuestionList{})

	server := httptransport.NewServer(categoryStore, questionStore)

	if err := http.ListenAndServe(":"+port, server); err != nil {
		log.Fatalf("could not listen on port %s %v", port, err)
	}
}
