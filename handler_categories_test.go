package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
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

		contentType := result.Header.Get("content-type")
		desiredContentType := "application/json"
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
