package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type CategoryStore interface {
	ListCategories() CategoryList
	GetCategory(categoryID string) Category
	AddCategory(categoryName string) Category
	RenameCategory(categoryID, categoryName string) Category
	DeleteCategory(categoryID string)
	CategoryIDExists(categoryID string) bool
	CategoryNameExists(categoryName string) bool
}

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type putRequestFormat struct {
	Name string `json:"name"`
}

func (c *CategoryServer) CategoryGetHandler(res http.ResponseWriter, req *http.Request) {
	var payload []byte
	var err error

	if strings.HasPrefix(req.URL.Path, "/categories/") && len(req.URL.Path) > len("/categories/") {
		categoryID := req.URL.Path[len("/categories/"):]
		category := c.store.GetCategory(categoryID)
		if category == (Category{}) {
			res.WriteHeader(http.StatusNotFound)
			return
		}

		payload, err = json.Marshal(category)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		categoryList := c.store.ListCategories()
		payload, err = json.Marshal(categoryList)
		if err != nil {
			log.Fatal(err)
		}
	}

	res.Write(payload)

}

func (c *CategoryServer) CategoryPostHandler(res http.ResponseWriter, req *http.Request) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	categoryName := string(requestBody)

	if c.store.CategoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		return
	}

	if !IsValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	category := c.store.AddCategory(categoryName)

	payload, err := json.Marshal(category)
	if err != nil {
		log.Fatal(err)
	}

	res.Header().Set("Location", fmt.Sprintf("/categories/%s", category.ID))
	res.WriteHeader(http.StatusCreated)
	res.Write(payload)
}

func (c *CategoryServer) CategoryPutHandler(res http.ResponseWriter, req *http.Request) {
	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var got Category

	err = json.Unmarshal(requestBody, &got)
	if err != nil {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryID := got.ID
	categoryName := got.Name

	if !c.store.CategoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if c.store.CategoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		return
	}

	if !IsValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	category := c.store.RenameCategory(got.ID, categoryName)

	res.WriteHeader(http.StatusOK)
	payload, err := json.Marshal(category)
	if err != nil {
		log.Fatal(err)
	}

	res.Write(payload)
}

func (c *CategoryServer) CategoryDeleteHandler(res http.ResponseWriter, req *http.Request) {
	type expectedFormat struct {
		ID string `json:"id"`
	}

	requestBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}

	var got expectedFormat

	err = json.Unmarshal(requestBody, &got)
	// json.unmarshall will not error if fields don't match
	// however got will be an empty struct, check that below
	if err != nil {
		log.Fatal(err)
		return
	}

	if got == (expectedFormat{}) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryID := got.ID

	if !c.store.CategoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	c.store.DeleteCategory(got.ID)

	res.WriteHeader(http.StatusNoContent)
}

func IsValidCategoryName(name string) bool {
	isValid := true

	if len(name) == 0 || len(name) > 32 {
		isValid = false
	}

	isLetterOrWhitespace := regexp.MustCompile(`^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+$`).MatchString
	if !isLetterOrWhitespace(name) {
		isValid = false
	}

	return isValid
}
