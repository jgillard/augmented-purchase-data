package handlers

import (
	"net/http"
)

type CategoryServer struct {
	store CategoryStore
	http.Handler
}

const jsonContentType = "application/json"

func NewCategoryServer(store CategoryStore) *CategoryServer {
	p := new(CategoryServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/status", http.HandlerFunc(p.statusHandler))
	router.Handle("/categories", http.HandlerFunc(p.categoriesHandler))
	router.Handle("/categories/", http.HandlerFunc(p.categoriesHandler))

	p.Handler = router

	return p
}

func (c *CategoryServer) categoriesHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", jsonContentType)

	switch req.Method {
	case http.MethodGet:
		c.CategoryGetHandler(res, req)
	case http.MethodPost:
		c.CategoryPostHandler(res, req)
	case http.MethodPut:
		c.CategoryPutHandler(res, req)
	case http.MethodDelete:
		c.CategoryDeleteHandler(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}

}
