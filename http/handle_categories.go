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

func (c *Server) categoryListHandler(res http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	categoryList := c.categoryStore.ListCategories()

	payload := marshallResponse(categoryList)

	res.Write(payload)
}

func (c *Server) categoryGetHandler(res http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	categoryID := ps.ByName("category")

	category := c.categoryStore.GetCategory(categoryID)

	if reflect.DeepEqual(category, internal.CategoryGetResponse{}) {
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

	var got internal.CategoryPostRequest
	unmarshallRequest(requestBody, &got)

	categoryName := got.Name

	if !ensureJSONFieldsPresent(res, got, internal.CategoryPostRequest{}) {
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

	if !c.categoryStore.CategoryParentIDExists(parentID) && parentID != "" {
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
