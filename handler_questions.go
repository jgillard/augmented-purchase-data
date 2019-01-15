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
	listQuestionsForCategory(categoryID string) QuestionList
	getQuestion(questionID string) Question
	addQuestion(categoryID string, question QuestionPostRequest) Question
	renameQuestion(questionID, questionTitle string) Question
	deleteQuestion(questionID string)
	questionIDExists(questionID string) bool
	questionTitleExists(categoryID, questionTitle string) bool
	questionBelongsToCategory(questionID, categoryID string) bool
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	CategoryID string     `json:"categoryID"`
	Type       string     `json:"type"`
	Options    OptionList `json:"options"`
}

// is this a very odd thing to do?
type QuestionPostRequest struct {
	Title   string    `json:"title"`
	Type    string    `json:"type"`
	Options *[]string `json:"options"`
}

type OptionList []Option

type Option struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

const questionTitleRegex = `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+\??$`

var possibleOptionTypes = []string{"string", "number"}

func (c *server) questionListHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	questionList := c.questionStore.listQuestionsForCategory(categoryID)

	payload := marshallResponse(questionList)

	res.Write(payload)
}

func (c *server) questionGetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	questionID := ps.ByName("question")

	question := c.questionStore.getQuestion(questionID)

	if reflect.DeepEqual(question, Question{}) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionNotFound))
		return
	}

	payload := marshallResponse(question)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *server) questionPostHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if !jsonIsValid(requestBody) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(errorInvalidJSON))
		return
	}

	var got QuestionPostRequest
	err = json.Unmarshal(requestBody, &got)
	if err != nil {
		// json.unmarshall explodes if options is not the correct shape, so catch that here
		switch t := err.(type) {
		case *json.UnmarshalTypeError:
			if t.Field == "options" {
				fmt.Println(`"options" object is wrong shape`)
				res.WriteHeader(http.StatusBadRequest)
				res.Write(craftErrorPayload(errorOptionsInvalid))
			}
		default:
			log.Fatal(err)
		}
		return
	}

	if !ensureJSONFieldsPresent(res, got, QuestionPostRequest{}) {
		res.Write(craftErrorPayload(errorFieldMissing))
		return
	}

	if !ensureStringFieldNonEmpty(res, "title", got.Title) {
		res.Write(craftErrorPayload(errorTitleEmpty))
		return
	}

	if !ensureStringFieldNonEmpty(res, "type", got.Type) {
		res.Write(craftErrorPayload(errorTypeEmpty))
		return
	}

	if !ensureStringFieldTitle(res, "type", got.Type, possibleOptionTypes) {
		res.Write(craftErrorPayload(errorInvalidType))
		return
	}

	// perform checks on nested options object
	if got.Options != nil {
		for _, opt := range *got.Options {
			if !ensureStringFieldNonEmpty(res, "options", opt) {
				res.Write(craftErrorPayload(errorOptionEmpty))
				return
			}
		}

		if !ensureNoDuplicates(res, "options", *got.Options) {
			res.Write(craftErrorPayload(errorDuplicateOption))
			return
		}
	}

	if c.questionStore.questionTitleExists(categoryID, got.Title) {
		fmt.Println(`"title" already exists`)
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateTitle))
		return
	}

	if c.categoryStore != nil && !c.categoryStore.categoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	question := c.questionStore.addQuestion(categoryID, got)

	payload := marshallResponse(question)

	res.Header().Set("Location", fmt.Sprintf("/categories/%s/questions/%s", categoryID, question.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)
}

func (c *server) questionPatchHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")
	questionID := ps.ByName("question")

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if !jsonIsValid(requestBody) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(errorInvalidJSON))
		return
	}

	var got jsonTitle
	unmarshallRequest(requestBody, &got)

	questionTitle := got.Title

	if !ensureJSONFieldsPresent(res, got, jsonTitle{}) {
		res.Write(craftErrorPayload(errorFieldMissing))
		return
	}

	if !isValidQuestionTitle(questionTitle) {
		fmt.Println(`"title" is not a valid string`)
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorInvalidTitle))
		return
	}

	if c.categoryStore != nil && !c.categoryStore.categoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	if !c.questionStore.questionIDExists(questionID) {
		fmt.Println(`"questionID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionNotFound))
		return
	}

	if !c.questionStore.questionBelongsToCategory(questionID, categoryID) {
		fmt.Println(`"questionID" in path doesn't belong to "categoryID" in path`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionDoesntBelongToCategory))
		return
	}

	if c.questionStore.questionTitleExists(categoryID, got.Title) {
		fmt.Println(`"title" already exists`)
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateTitle))
		return
	}

	question := c.questionStore.renameQuestion(questionID, questionTitle)
	payload := marshallResponse(question)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *server) questionDeleteHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")
	questionID := ps.ByName("question")

	if c.categoryStore != nil && !c.categoryStore.categoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	if !c.questionStore.questionIDExists(questionID) {
		fmt.Println(`"questionID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionNotFound))
		return
	}

	if !c.questionStore.questionBelongsToCategory(questionID, categoryID) {
		fmt.Println(`"questionID" in path doesn't belong to "categoryID" in path`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionDoesntBelongToCategory))
		return
	}

	c.questionStore.deleteQuestion(questionID)

	payload := marshallResponse(jsonStatus{statusDeleted})

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func isValidQuestionTitle(title string) bool {
	isValid := true

	if len(title) == 0 || len(title) > 32 {
		isValid = false
	}

	isLetterOrWhitespace := regexp.MustCompile(questionTitleRegex).MatchString
	if !isLetterOrWhitespace(title) {
		isValid = false
	}

	return isValid
}
