package handlers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
)

type CategoryStore interface {
	ListCategories() CategoryList
	AddCategory(categoryName string) Category
	CategoryNameExists(categoryName string) bool
}

type CategoryList struct {
	Categories []Category `json:"categories"`
}

type Category struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (c *CategoryServer) CategoryGetHandler(res http.ResponseWriter, req *http.Request) {
	categoryList := c.store.ListCategories()

	bytePayload, err := json.Marshal(categoryList)
	if err != nil {
		log.Fatal(err)
	}

	res.Write(bytePayload)
	res.WriteHeader(http.StatusOK)
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

	if !isValidCategoryName(categoryName) {
		res.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	category := c.store.AddCategory(categoryName)

	res.WriteHeader(http.StatusCreated)
	payload, err := json.Marshal(category)
	if err != nil {
		log.Fatal(err)
	}

	res.Write([]byte(payload))
}

func isValidCategoryName(name string) bool {
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
