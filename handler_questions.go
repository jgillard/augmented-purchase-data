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
	DeleteQuestion(questionID string)
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

var possibleOptionTypes = []string{"string", "number"}

func (c *Server) QuestionListHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	res.Header().Set("Content-Type", jsonContentType)

	categoryID := ps.ByName("category")

	questionList := c.questionStore.ListQuestionsForCategory(categoryID)

	payload := marshallResponse(questionList)

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

	if !ensureJSONFieldsPresent(res, got, QuestionPostRequest{}) {
		return
	}

	if !ensureStringFieldNonEmpty(res, "value", got.Value) {
		return
	}

	if !ensureStringFieldNonEmpty(res, "type", got.Type) {
		return
	}

	if !ensureStringFieldValue(res, "type", got.Type, possibleOptionTypes) {
		return
	}

	for _, opt := range got.Options {
		if !ensureStringFieldNonEmpty(res, "options", opt) {
			return
		}
	}

	if !ensureNoDuplicates(res, "options", got.Options) {
		return
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

	questionValue := got.Name

	if !ensureJSONFieldsPresent(res, got, jsonName{}) {
		return
	}

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
	payload := marshallResponse(question)

	res.WriteHeader(http.StatusOK)
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

	c.questionStore.DeleteQuestion(questionID)

	res.WriteHeader(http.StatusNoContent)
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

func ensureJSONFieldsPresent(res http.ResponseWriter, got, desired interface{}) bool {
	if reflect.DeepEqual(got, desired) {
		fmt.Println("json field(s) missing from request")
		res.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func ensureStringFieldNonEmpty(res http.ResponseWriter, key, value string) bool {
	if value == "" {
		fmt.Println(fmt.Sprintf(`"%s" missing from request`, key))
		res.WriteHeader(http.StatusBadRequest)
		return false
	}
	return true
}

func ensureStringFieldValue(res http.ResponseWriter, key, value string, possibleOptionTypes []string) bool {
	isValid := false

	for _, possible := range possibleOptionTypes {
		if value == possible {
			isValid = true
			break
		}
	}

	if !isValid {
		fmt.Printf(`"%s" must be one of %v`, key, possibleOptionTypes)
		res.WriteHeader(http.StatusBadRequest)
	}

	return isValid
}

func ensureNoDuplicates(res http.ResponseWriter, key string, strings []string) bool {
	noDuplicates := true

	counts := make(map[string]int)
	for _, str := range strings {
		counts[str]++
	}

	for _, count := range counts {
		if count > 1 {
			fmt.Printf(`"%s" contains duplicate strings`, key)
			res.WriteHeader(http.StatusBadRequest)
			noDuplicates = false
		}
	}

	return noDuplicates
}
