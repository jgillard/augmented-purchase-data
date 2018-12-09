package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestListCategories(t *testing.T) {

	server := NewCategoryServer()

	t.Run("check response code", func(t *testing.T) {
		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		status := result.StatusCode
		desiredStatus := http.StatusOK
		assertNumbersEqual(t, status, desiredStatus)
	})

	t.Run("check content-type header is application/json", func(t *testing.T) {
		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		contentType := result.Header.Get("Content-Type")
		desiredContentType := jsonContentType
		assertStringsEqual(t, contentType, desiredContentType)
	})

	t.Run("return a list of categories with IDs & names", func(t *testing.T) {
		// this is kinda testing the marshalling/unmarshalling rather than the json responses themselves

		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		bodyBytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var categories CategoryList

		err = json.Unmarshal(bodyBytes, &categories)
		if err != nil {
			t.Fatal(err)
		}

		desiredBody := CategoryList{
			Categories: []Category{
				{ID: "a1b2", Name: "foo"},
			},
		}

		if !reflect.DeepEqual(categories, desiredBody) {
			t.Errorf("got '%v' wanted '%v'", categories, desiredBody)
		}

	})

	t.Run("test the json itself", func(t *testing.T) {
		// surely don't ned both of these, what's the best practice?

		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		bodyBytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		bodyString := string(bodyBytes)
		desiredBodyString := `{"categories":[{"id":"a1b2","name":"foo"}]}`

		assertStringsEqual(t, bodyString, desiredBodyString)

	})
}

func TestAddCategory(t *testing.T) {

	server := NewCategoryServer()

	t.Run("check response code", func(t *testing.T) {
		body := strings.NewReader("foo")
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		status := result.StatusCode
		desiredStatus := http.StatusCreated
		assertNumbersEqual(t, status, desiredStatus)
	})

	t.Run("check Content-Type header is application/json", func(t *testing.T) {
		body := strings.NewReader("foo")
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		contentType := result.Header.Get("Content-Type")
		desiredContentType := jsonContentType
		assertStringsEqual(t, contentType, desiredContentType)
	})

	t.Run("check response body contains a category with ID & name", func(t *testing.T) {
		categoryName := "accommodation"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		bodyBytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var category Category

		err = json.Unmarshal(bodyBytes, &category)
		if err != nil {
			t.Fatal(err)
		}

		if category.Name != categoryName {
			t.Errorf("got name '%s' wanted '%s'", category.Name, categoryName)
		}

		if !isXid(t, category.ID) {
			t.Errorf("got ID '%s' which isn't an xid", category.ID)
		}
	})
}
