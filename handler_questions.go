package transactioncategories

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

type QuestionStore interface {
	ListQuestionsForCategory(categoryID string) QuestionList
	AddQuestion(categoryID string, question QuestionPostRequest) Question
	RenameQuestion(questionID, questionValue string) Question
	questionIDExists(questionID string) bool
	questionValueExists(categoryID, questionValue string) bool
	questionBelongsToCategory(questionID, categoryID string) bool
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
	Value   string   `json:"value"`
	Type    string   `json:"type"`
	Options []string `json:"options"`
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
				fmt.Println(`"options" object is wrong shape`)
				res.WriteHeader(http.StatusBadRequest)
				return
			}
		default:
			log.Fatal(err)
			return
		}
	}

	if reflect.DeepEqual(got, QuestionPostRequest{}) {
		fmt.Println("json field(s) missing from request")
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if got.Value == "" {
		fmt.Println(`"value" missing from request`)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if got.Type == "" {
		fmt.Println(`"type" missing from request`)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if got.Type != "string" && got.Type != "number" {
		fmt.Println(`"type" must be "string" or "number"`)
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	for _, opt := range got.Options {
		if opt == "" {
			fmt.Println(`"option" strings cannot be empty`)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	// detect duplicates in options list
	optionCounts := make(map[string]int)
	for _, opt := range got.Options {
		optionCounts[opt]++
	}
	for _, count := range optionCounts {
		if count > 1 {
			fmt.Println(`"option" contains duplicate strings`)
			res.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	if c.questionStore.questionValueExists(categoryID, got.Value) {
		fmt.Println(`"value" already exists`)
		res.WriteHeader(http.StatusConflict)
		return
	}

	if c.categoryStore != nil && !c.categoryStore.categoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	question := c.questionStore.AddQuestion(categoryID, got)

	payload := marshallResponse(question)

	res.Header().Set("Location", fmt.Sprintf("/categories/%s/questions/%s", categoryID, question.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)

}

func (c *Server) QuestionPatchHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("Content-Type", jsonContentType)

	categoryID := ps.ByName("category")
	questionID := ps.ByName("question")

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var got jsonName
	UnmarshallRequest(requestBody, &got)

	if got == (jsonName{}) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	questionValue := got.Name

	if !IsValidQuestionName(questionValue) {
		fmt.Println(`"name" is not a valid string`)
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if c.categoryStore != nil && !c.categoryStore.categoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if !c.questionStore.questionIDExists(questionID) {
		fmt.Println(`"questionID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if !c.questionStore.questionBelongsToCategory(questionID, categoryID) {
		fmt.Println(`"questionID" in path doesn't belong to "categoryID" in path`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if c.questionStore.questionValueExists(categoryID, got.Name) {
		fmt.Println(`"name" already exists`)
		res.WriteHeader(http.StatusConflict)
		return
	}

	question := c.questionStore.RenameQuestion(questionID, questionValue)

	res.WriteHeader(http.StatusOK)
	payload := marshallResponse(question)
	res.Write(payload)

}

func (c *Server) QuestionDeleteHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("Content-Type", jsonContentType)

	categoryID := ps.ByName("category")
	questionID := ps.ByName("question")

	if c.categoryStore != nil && !c.categoryStore.categoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if !c.questionStore.questionIDExists(questionID) {
		fmt.Println(`"questionID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if !c.questionStore.questionBelongsToCategory(questionID, categoryID) {
		fmt.Println(`"questionID" in path doesn't belong to "categoryID" in path`)
		res.WriteHeader(http.StatusNotFound)
		return
	}

}

func IsValidQuestionName(name string) bool {
	isValid := true

	if len(name) == 0 || len(name) > 32 {
		isValid = false
	}

	questionRegex := `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+\??$`
	isLetterOrWhitespace := regexp.MustCompile(questionRegex).MatchString
	if !isLetterOrWhitespace(name) {
		isValid = false
	}

	return isValid
}
