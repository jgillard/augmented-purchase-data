package handlers

import (
	"net/http"
)

type CategoryServer struct {
	http.Handler
}

const jsonContentType = "application/json"

func NewCategoryServer() *CategoryServer {
	p := new(CategoryServer)

	router := http.NewServeMux()
	router.Handle("/status", http.HandlerFunc(StatusHandler))
	router.Handle("/categories", http.HandlerFunc(ListCategoriesHandler))

	p.Handler = router

	return p
}
