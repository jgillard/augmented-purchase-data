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

type CategoryStore interface {
	ListCategories() CategoryList
	GetCategory(categoryID string) CategoryGetResponse
	AddCategory(categoryName, parentID string) Category
	RenameCategory(categoryID, categoryName string) Category
	DeleteCategory(categoryID string)
	categoryIDExists(categoryID string) bool
	categoryNameExists(categoryName string) bool
	categoryParentIDExists(categoryParentID string) bool
	getCategoryDepth(categoryID string) int
}

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentID"`
}

// is this a very odd thing to do?
type CategoryGetResponse struct {
	Category
	Children []Category `json:"children"`
}

// is this a very odd thing to do?
type CategoryPostRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentID"`
}

func (c *Server) CategoryListHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// GET the list of categories
	categoryList := c.categoryStore.ListCategories()
	payload := marshallResponse(categoryList)

	res.Write(payload)
}

func (c *Server) CategoryGetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	// GET a specific category
	categoryID := ps.ByName("category")
	category := c.categoryStore.GetCategory(categoryID)

	if reflect.DeepEqual(category, CategoryGetResponse{}) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(ErrorCategoryNotFound))
		return
	}

	payload := marshallResponse(category)

	res.Write(payload)
}

func (c *Server) CategoryPostHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if !jsonIsValid(requestBody) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(ErrorInvalidJSON))
		return
	}

	var got CategoryPostRequest
	unmarshallRequest(requestBody, &got)

	categoryName := got.Name

	if !ensureJSONFieldsPresent(res, got, CategoryPostRequest{}) {
		res.Write(craftErrorPayload(ErrorFieldMissing))
		return
	}

	if c.categoryStore.categoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(ErrorDuplicateCategoryName))
		return
	}

	if !isValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(ErrorInvalidCategoryName))
		return
	}

	// parentID not supplied
	if got.ParentID == nil {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(ErrorFieldMissing))
		return
	}

	parentID := *got.ParentID

	if !c.categoryStore.categoryParentIDExists(parentID) && parentID != "" {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(ErrorParentIDNotFound))
		return
	}

	// checks for parent already a subcategory (depth zero indexed)
	// we currently confine to 2 levels of categories
	if c.categoryStore.getCategoryDepth(parentID) == 1 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(ErrorCategoryTooNested))
		return
	}

	category := c.categoryStore.AddCategory(categoryName, parentID)

	payload := marshallResponse(category)

	res.Header().Set("Location", fmt.Sprintf("/categories/%s", category.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)
}

func (c *Server) CategoryPatchHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	if !jsonIsValid(requestBody) {
		res.WriteHeader(http.StatusBadRequest)
		res.Write(craftErrorPayload(ErrorInvalidJSON))
		return
	}

	var got jsonName
	unmarshallRequest(requestBody, &got)

	categoryName := got.Name

	if !ensureJSONFieldsPresent(res, got, jsonName{}) {
		res.Write(craftErrorPayload(ErrorFieldMissing))
		return
	}

	if !c.categoryStore.categoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(ErrorCategoryNotFound))
		return
	}

	if c.categoryStore.categoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(ErrorDuplicateCategoryName))
		return
	}

	if !isValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(ErrorInvalidCategoryName))
		return
	}

	category := c.categoryStore.RenameCategory(categoryID, categoryName)

	payload := marshallResponse(category)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *Server) CategoryDeleteHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	if !c.categoryStore.categoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(ErrorCategoryNotFound))
		return
	}

	c.categoryStore.DeleteCategory(categoryID)

	payload := marshallResponse(jsonStatus{"deleted"})

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func isValidCategoryName(name string) bool {
	isValid := true

	if len(name) == 0 || len(name) > 32 {
		isValid = false
	}

	stringRegex := `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+$`
	isLetterOrWhitespace := regexp.MustCompile(stringRegex).MatchString
	if !isLetterOrWhitespace(name) {
		isValid = false
	}

	return isValid
}
