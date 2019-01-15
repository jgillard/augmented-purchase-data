package transactioncategories

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"regexp"

	"github.com/julienschmidt/httprouter"
)

// CategoryStore is an interface that when implemented,
// provides methods for manipulating a store of categories,
// including some helper functions for querying the store
type CategoryStore interface {
	listCategories() CategoryList
	getCategory(categoryID string) CategoryGetResponse
	addCategory(categoryName, parentID string) Category
	renameCategory(categoryID, categoryName string) Category
	deleteCategory(categoryID string)
	categoryIDExists(categoryID string) bool
	categoryNameExists(categoryName string) bool
	categoryParentIDExists(categoryParentID string) bool
	getCategoryDepth(categoryID string) int
}

// CategoryList stores multiple categories
type CategoryList struct {
	Categories []Category `json:"categories"`
}

// Category stores all expected category attributes
// The structure implements the adjacency list pattern
type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentID"`
}

// CategoryGetResponse returns a Category in addition to its immediate child categories
type CategoryGetResponse struct {
	Category
	Children []Category `json:"children"`
}

// CategoryPostRequest is a Category with no ID
// Used for sending new Categorys to the server
// ParentID is a string to allow "" to signify a top-level Category
type CategoryPostRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentID"`
}

const categoryNameRegex = `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+$`

func (c *Server) categoryListHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	categoryList := c.categoryStore.listCategories()

	payload := marshallResponse(categoryList)

	res.Write(payload)
}

func (c *Server) categoryGetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	category := c.categoryStore.getCategory(categoryID)

	if reflect.DeepEqual(category, CategoryGetResponse{}) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	payload := marshallResponse(category)

	res.Write(payload)
}

func (c *Server) categoryPostHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if !jsonIsValid(requestBody) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(errorInvalidJSON))
		return
	}

	var got CategoryPostRequest
	unmarshallRequest(requestBody, &got)

	categoryName := got.Name

	if !ensureJSONFieldsPresent(res, got, CategoryPostRequest{}) {
		res.Write(craftErrorPayload(errorFieldMissing))
		return
	}

	if c.categoryStore.categoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateCategoryName))
		return
	}

	if !isValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorInvalidCategoryName))
		return
	}

	// parentID not supplied
	if got.ParentID == nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(errorFieldMissing))
		return
	}

	parentID := *got.ParentID

	if !c.categoryStore.categoryParentIDExists(parentID) && parentID != "" {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorParentIDNotFound))
		return
	}

	// checks for parent already a subcategory (depth zero indexed)
	// we currently confine to 2 levels of categories
	if c.categoryStore.getCategoryDepth(parentID) == 1 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorCategoryTooNested))
		return
	}

	category := c.categoryStore.addCategory(categoryName, parentID)

	payload := marshallResponse(category)

	res.Header().Set("Location", fmt.Sprintf("/categories/%s", category.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)
}

func (c *Server) categoryPatchHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

	var got jsonName
	unmarshallRequest(requestBody, &got)

	categoryName := got.Name

	if !ensureJSONFieldsPresent(res, got, jsonName{}) {
		res.Write(craftErrorPayload(errorFieldMissing))
		return
	}

	if !c.categoryStore.categoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	if c.categoryStore.categoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateCategoryName))
		return
	}

	if !isValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorInvalidCategoryName))
		return
	}

	category := c.categoryStore.renameCategory(categoryID, categoryName)

	payload := marshallResponse(category)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *Server) categoryDeleteHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	if !c.categoryStore.categoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	c.categoryStore.deleteCategory(categoryID)

	payload := marshallResponse(jsonStatus{statusDeleted})

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func isValidCategoryName(name string) bool {
	isValid := true

	if len(name) == 0 || len(name) > 32 {
		isValid = false
	}

	isLetterOrWhitespace := regexp.MustCompile(categoryNameRegex).MatchString
	if !isLetterOrWhitespace(name) {
		isValid = false
	}

	return isValid
}
