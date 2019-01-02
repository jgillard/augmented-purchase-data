package transactioncategories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type QuestionStore interface {
	ListQuestionsForCategory(categoryID string) QuestionList
	AddQuestion(categoryID string, question QuestionPostRequest) Question
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID         string     `json:"id"`
	Value      string     `json:"value"`
	CategoryID string     `json:"categoryID"`
	Type       string     `json:"type"`
	Options    OptionList `json:"options"`
}

// is this a very odd thing to do?
type QuestionPostRequest struct {
	Value   string     `json:"value"`
	Type    string     `json:"type"`
	Options OptionList `json:"options"`
}

type OptionList []Option

type Option struct {
	ID    string `json:"id"`
	Value string `json:"value"`
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

func (c *Server) QuestionPostHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("Content-Type", jsonContentType)

	categoryID := ps.ByName("category")

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var got QuestionPostRequest
	err = json.Unmarshal(requestBody, &got)
	// json.unmarshall will not error if fields don't match
	// however got will be an empty struct, check that below
	if err != nil {
		// however it does blow if options is not the correct shape, so catch that here
		switch t := err.(type) {
		case *json.UnmarshalTypeError:
			if t.Field == "options" {
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			log.Fatal(err)
			return
		}
	}

	question := c.questionStore.AddQuestion(categoryID, got)

	payload := marshallResponse(question)

	res.Header().Set("Location", fmt.Sprintf("/categories/%s/questions/%s", categoryID, question.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)

}
