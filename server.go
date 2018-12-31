package transactioncategories

import (
	"net/http"
)

type Server struct {
	categoryStore CategoryStore
	questionStore QuestionStore
	http.Handler
}

const jsonContentType = "application/json"

func NewServer(cats CategoryStore, questions QuestionStore) *Server {
	p := new(Server)

	p.categoryStore = cats
	p.questionStore = questions

	router := http.NewServeMux()
	router.Handle("/status", http.HandlerFunc(p.statusHandler))
	router.Handle("/categories", http.HandlerFunc(p.categoriesHandler))
	router.Handle("/categories/", http.HandlerFunc(p.categoriesHandler))

	p.Handler = router

	return p
}

func (c *Server) categoriesHandler(res http.ResponseWriter, req *http.Request) {

	res.Header().Set("Content-Type", jsonContentType)

	switch req.Method {
	case http.MethodGet:
		c.CategoryGetHandler(res, req)
	case http.MethodPost:
		c.CategoryPostHandler(res, req)
	case http.MethodPatch:
		c.CategoryPatchHandler(res, req)
	case http.MethodDelete:
		c.CategoryDeleteHandler(res, req)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}

}
