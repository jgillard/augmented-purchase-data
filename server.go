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

type Middleware struct {
	handler http.Handler
}

func (m *Middleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(ContentTypeKey, jsonContentType)
	m.handler.ServeHTTP(res, req)
}

func NewServer(cats CategoryStore, questions QuestionStore) *Server {
	p := new(Server)

	p.categoryStore = cats
	p.questionStore = questions

	router := httprouter.New()
	router.GET("/status", p.statusHandler)

	router.GET("/categories", p.CategoryListHandler)
	router.GET("/categories/:category", p.CategoryGetHandler)
	router.POST("/categories", p.CategoryPostHandler)
	router.PATCH("/categories/:category", p.CategoryPatchHandler)
	router.DELETE("/categories/:category", p.CategoryDeleteHandler)

	router.GET("/categories/:category/questions", p.QuestionListHandler)
	router.POST("/categories/:category/questions", p.QuestionPostHandler)
	router.PATCH("/categories/:category/questions/:question", p.QuestionPatchHandler)
	router.DELETE("/categories/:category/questions/:question", p.QuestionDeleteHandler)

	router.NotFound = http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusNotFound)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusMethodNotAllowed)
	})

	p.Handler = &Middleware{router}

	return p
}
