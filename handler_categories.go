package transactioncategories

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
	categoryIDExists(categoryID string) bool
	categoryNameExists(categoryName string) bool
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

type jsonID struct {
	ID string `json:"id"`
}

type jsonName struct {
	Name string `json:"name"`
}

func (c *CategoryServer) CategoryGetHandler(res http.ResponseWriter, req *http.Request) {
	var payload []byte
	var err error

	if strings.HasPrefix(req.URL.Path, "/categories/") && len(req.URL.Path) > len("/categories/") {
		// GET a specific category
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
		// GET the list of categories
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

	var got jsonName
	UnmarshallRequest(requestBody, &got)

	if got == (jsonName{}) {
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	categoryName := got.Name

	if c.store.categoryNameExists(categoryName) {
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

func (c *CategoryServer) CategoryPatchHandler(res http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Path[len("/categories/"):]

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

	categoryName := got.Name

	if !c.store.categoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	if c.store.categoryNameExists(categoryName) {
		res.WriteHeader(http.StatusConflict)
		return
	}

	if !IsValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	category := c.store.RenameCategory(categoryID, categoryName)

	res.WriteHeader(http.StatusOK)
	payload, err := json.Marshal(category)
	if err != nil {
		log.Fatal(err)
	}

	res.Write(payload)
}

func (c *CategoryServer) CategoryDeleteHandler(res http.ResponseWriter, req *http.Request) {
	categoryID := req.URL.Path[len("/categories/"):]

	if !c.store.categoryIDExists(categoryID) {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	c.store.DeleteCategory(categoryID)

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

func UnmarshallRequest(body []byte, got interface{}) {
	err := json.Unmarshal(body, got)
	// json.unmarshall will not error if fields don't match
	// however got will be an empty struct, check that below
	if err != nil {
		log.Fatal(err)
		return
	}
}
