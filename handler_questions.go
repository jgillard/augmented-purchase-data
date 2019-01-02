package transactioncategories

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type QuestionStore interface {
	ListQuestionsForCategory(categoryID string) QuestionList
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID         string `json:"id"`
	Value      string `json:"value"`
	CategoryID string `json:"categoryID"`
	Type       string `json:"type"`
}

func (c *Server) QuestionListHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("Content-Type", jsonContentType)

	categoryID := ps.ByName("category")

	questionList := c.questionStore.ListQuestionsForCategory(categoryID)
	payload, err := json.Marshal(questionList)
	if err != nil {
		log.Fatal(err)
	}

	res.Write(payload)

}
