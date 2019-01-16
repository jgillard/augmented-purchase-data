package httptransport

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	internal "github.com/jgillard/practising-go-tdd/internal"
	"github.com/julienschmidt/httprouter"
)

func (c *Server) questionListHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	questionList := c.questionStore.ListQuestionsForCategory(categoryID)

	payload := marshallResponse(questionList)

	res.Write(payload)
}

func (c *Server) questionGetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	questionID := ps.ByName("question")

	question := c.questionStore.GetQuestion(questionID)

	if reflect.DeepEqual(question, internal.Question{}) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionNotFound))
		return
	}

	payload := marshallResponse(question)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *Server) questionPostHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

	var got internal.QuestionPostRequest
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

	if !ensureJSONFieldsPresent(res, got, internal.QuestionPostRequest{}) {
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

	if !internal.IsValidOptionType(got.Type) {
		res.WriteHeader(http.StatusBadRequest)
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

	if c.questionStore.QuestionTitleExists(categoryID, got.Title) {
		fmt.Println(`"title" already exists`)
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateTitle))
		return
	}

	if c.categoryStore != nil && !c.categoryStore.CategoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	question := c.questionStore.AddQuestion(categoryID, got)

	payload := marshallResponse(question)

	res.Header().Set("Location", fmt.Sprintf("/categories/%s/questions/%s", categoryID, question.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)
}

func (c *Server) questionPatchHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

	if !internal.IsValidQuestionTitle(questionTitle) {
		fmt.Println(`"title" is not a valid string`)
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorInvalidTitle))
		return
	}

	if c.categoryStore != nil && !c.categoryStore.CategoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	if !c.questionStore.QuestionIDExists(questionID) {
		fmt.Println(`"questionID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionNotFound))
		return
	}

	if !c.questionStore.QuestionBelongsToCategory(questionID, categoryID) {
		fmt.Println(`"questionID" in path doesn't belong to "categoryID" in path`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionDoesntBelongToCategory))
		return
	}

	if c.questionStore.QuestionTitleExists(categoryID, got.Title) {
		fmt.Println(`"title" already exists`)
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateTitle))
		return
	}

	question := c.questionStore.RenameQuestion(questionID, questionTitle)
	payload := marshallResponse(question)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *Server) questionDeleteHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")
	questionID := ps.ByName("question")

	if c.categoryStore != nil && !c.categoryStore.CategoryIDExists(categoryID) {
		fmt.Println(`"categoryID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	if !c.questionStore.QuestionIDExists(questionID) {
		fmt.Println(`"questionID" in path doesn't exist`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionNotFound))
		return
	}

	if !c.questionStore.QuestionBelongsToCategory(questionID, categoryID) {
		fmt.Println(`"questionID" in path doesn't belong to "categoryID" in path`)
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorQuestionDoesntBelongToCategory))
		return
	}

	c.questionStore.DeleteQuestion(questionID)

	payload := marshallResponse(jsonStatus{statusDeleted})

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}
