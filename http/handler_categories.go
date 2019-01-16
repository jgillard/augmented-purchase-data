package httptransport

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"

	internal "github.com/jgillard/practising-go-tdd/internal"
)

// CategoryGetResponse returns a Category in addition to its immediate child categories
type CategoryGetResponse struct {
	internal.Category
	Children []internal.Category `json:"children"`
}

// CategoryPostRequest is a Category with no ID
// Used for sending new Categorys to the server
// ParentID is a string to allow "" to signify a top-level Category
type CategoryPostRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentID"`
}

func (c *Server) categoryListHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	categoryList := c.categoryStore.ListCategories()

	payload := marshallResponse(categoryList)

	res.Write(payload)
}

func (c *Server) categoryGetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	category := c.categoryStore.GetCategory(categoryID)

	if reflect.DeepEqual(category, internal.Category{}) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	children := c.categoryStore.GetChildCategories(categoryID)

	responseStruct := CategoryGetResponse{
		category,
		children,
	}

	payload := marshallResponse(responseStruct)

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

	if c.categoryStore.CategoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateCategoryName))
		return
	}

	if !internal.IsValidCategoryName(categoryName) {
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

	if !c.categoryStore.CategoryIDExists(parentID) && parentID != "" {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorParentIDNotFound))
		return
	}

	// checks for parent already a subcategory (depth zero indexed)
	// we currently confine to 2 levels of categories
	if c.categoryStore.GetCategoryDepth(parentID) == 1 {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorCategoryTooNested))
		return
	}

	category := c.categoryStore.AddCategory(categoryName, parentID)

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

	if !c.categoryStore.CategoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	if c.categoryStore.CategoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		res.Write(craftErrorPayload(errorDuplicateCategoryName))
		return
	}

	if !internal.IsValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		res.Write(craftErrorPayload(errorInvalidCategoryName))
		return
	}

	category := c.categoryStore.RenameCategory(categoryID, categoryName)

	payload := marshallResponse(category)

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}

func (c *Server) categoryDeleteHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	if !c.categoryStore.CategoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		res.Write(craftErrorPayload(errorCategoryNotFound))
		return
	}

	c.categoryStore.DeleteCategory(categoryID)

	payload := marshallResponse(jsonStatus{statusDeleted})

	res.WriteHeader(http.StatusOK)
	res.Write(payload)
}
