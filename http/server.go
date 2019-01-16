package httptransport

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	internal "github.com/jgillard/practising-go-tdd/internal"
)

// Server matches the interface of http.Handler
// and adds a question and category store
type Server struct {
	categoryStore internal.CategoryStore
	questionStore internal.QuestionStore
	http.Handler
}

const jsonContentType = "application/json"

type middleware struct {
	handler http.Handler
}

func (m *middleware) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(contentTypeKey, jsonContentType)
	m.handler.ServeHTTP(res, req)
}

// NewServer returns a category & question server,
// with a router & middleware
func NewServer(cats internal.CategoryStore, questions internal.QuestionStore) *Server {
	p := new(Server)

	p.categoryStore = cats
	p.questionStore = questions

	router := httprouter.New()
	router.GET("/status", p.statusHandler)

	router.GET("/categories", p.categoryListHandler)
	router.GET("/categories/:category", p.categoryGetHandler)
	router.POST("/categories", p.categoryPostHandler)
	router.PATCH("/categories/:category", p.categoryPatchHandler)
	router.DELETE("/categories/:category", p.categoryDeleteHandler)

	router.GET("/categories/:category/questions", p.questionListHandler)
	router.GET("/categories/:category/questions/:question", p.questionGetHandler)
	router.POST("/categories/:category/questions", p.questionPostHandler)
	router.PATCH("/categories/:category/questions/:question", p.questionPatchHandler)
	router.DELETE("/categories/:category/questions/:question", p.questionDeleteHandler)

	router.NotFound = http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusNotFound)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(res http.ResponseWriter, _ *http.Request) {
		res.WriteHeader(http.StatusMethodNotAllowed)
	})

	p.Handler = &middleware{router}

	return p
}
