package transactioncategories

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
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

	router := httprouter.New()
	router.GET("/status", p.statusHandler)
	router.GET("/categories", p.categoriesHandler)
	router.GET("/categories/:category", p.categoriesHandler)
	router.POST("/categories", p.categoriesHandler)
	router.PATCH("/categories/:category", p.categoriesHandler)
	router.DELETE("/categories/:category", p.categoriesHandler)

	router.NotFound = http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusNotFound)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusMethodNotAllowed)
	})

	p.Handler = router

	return p
}

func (c *Server) categoriesHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	res.Header().Set("Content-Type", jsonContentType)

	switch req.Method {
	case http.MethodGet:
		c.CategoryGetHandler(res, req, nil)
	case http.MethodPost:
		c.CategoryPostHandler(res, req, nil)
	case http.MethodPatch:
		c.CategoryPatchHandler(res, req, nil)
	case http.MethodDelete:
		c.CategoryDeleteHandler(res, req, nil)
	default:
		res.WriteHeader(http.StatusMethodNotAllowed)
	}

}
