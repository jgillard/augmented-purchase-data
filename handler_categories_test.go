package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/rs/xid"
)

func TestListCategories(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation"},
			Category{ID: "5678", Name: "food and drink"},
		},
	}

	store := &StubCategoryStore{categoryList}
	server := NewCategoryServer(store)

	t.Run("check response code", func(t *testing.T) {
		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusOK
		assertNumbersEqual(t, got, want)
	})

	t.Run("check content-type header is application/json", func(t *testing.T) {
		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("Content-Type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("return the list of categories with IDs & names", func(t *testing.T) {
		// this is kinda testing the marshalling/unmarshalling rather than the json responses themselves?

		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		bodyBytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var got CategoryList

		err = json.Unmarshal(bodyBytes, &got)
		if err != nil {
			t.Fatal(err)
		}

		want := categoryList

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got '%v' wanted '%v'", got, want)
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

		got := string(bodyBytes)
		want := `{"categories":[{"id":"1234","name":"accommodation"},{"id":"5678","name":"food and drink"}]}`

		assertStringsEqual(t, got, want)

	})
}

func TestAddCategory(t *testing.T) {

	store := &StubCategoryStore{
		CategoryList{
			Categories: []Category{
				Category{ID: "1234", Name: "accommodation"},
			},
		},
	}
	server := NewCategoryServer(store)

	t.Run("check Content-Type header is application/json", func(t *testing.T) {
		body := strings.NewReader("foo")
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("Content-Type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("check response body contains a category with ID & correct name", func(t *testing.T) {
		categoryName := "newName"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		bodyBytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var got Category

		err = json.Unmarshal(bodyBytes, &got)
		if err != nil {
			t.Fatal(err)
		}

		assertStringsEqual(t, got.Name, categoryName)

		if !isXid(t, got.ID) {
			t.Errorf("got ID '%s' which isn't an xid", got.ID)
		}
	})

	t.Run("check can add new name", func(t *testing.T) {
		categoryName := "food and drink"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusCreated

		assertNumbersEqual(t, got, want)
	})

	t.Run("check cannot add a duplicate name", func(t *testing.T) {
		categoryName := "accommodation"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusConflict

		assertNumbersEqual(t, got, want)
	})

	t.Run("name must not be invalid", func(t *testing.T) {

		cases := map[string]string{
			"empty string":         "",
			"only whitespace":      "     ",
			"leading whitespace":   " foobar",
			"trailing whitespace":  "foobar ",
			"string over 32 chars": "abcdefhijklmnopqrstuvwxyzabcdefgh",
			"numeric":              "123",
			"punctuation chars":    "!@$",
		}

		for name, value := range cases {
			t.Run(fmt.Sprintf("name must not be %s", name), func(t *testing.T) {
				body := strings.NewReader(value)

				req := NewPostRequest(t, "/categories", body)
				res := httptest.NewRecorder()
				server.ServeHTTP(res, req)
				result := res.Result()

				got := result.StatusCode
				want := http.StatusUnprocessableEntity

				if got != want {
					t.Errorf("got status code %d expected %d, when posting categoryName '%s'", got, want, value)
				}
			})
		}

	})

}

func TestRenameCategory(t *testing.T) {

	existingCategory := Category{ID: "1234", Name: "accommodation"}
	store := &StubCategoryStore{
		CategoryList{
			Categories: []Category{
				existingCategory,
			},
		},
	}
	server := NewCategoryServer(store)

	t.Run("check Content-Type header is application/json", func(t *testing.T) {
		body := []byte("foo")
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("Content-Type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("rejects request if body not json category", func(t *testing.T) {
		body := []byte("foo")
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusBadRequest
		assertNumbersEqual(t, got, want)
	})

	t.Run("can rename existing category to valid name", func(t *testing.T) {
		newName := "newName"
		renamedCategory := Category{ID: "1234", Name: newName}
		body, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		gotStatus := result.StatusCode
		wantStatus := http.StatusOK
		assertNumbersEqual(t, gotStatus, wantStatus)

		gotBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}
		wantBody := body
		assertStringsEqual(t, string(gotBody), string(wantBody))

		if store.categories.Categories[0].Name != newName {
			t.Fatalf("store not updated: got '%s' want '%s'", gotBody, wantBody)
		}

	})
	t.Run("cannot rename existing category to invalid name", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "foo/*!bar"}
		body, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusUnprocessableEntity
		assertNumbersEqual(t, got, want)
	})
	t.Run("cannot rename non-existent category", func(t *testing.T) {
		nonExistentCategory := Category{ID: "5678", Name: "foobar"}
		body, err := json.Marshal(nonExistentCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusNotFound
		assertNumbersEqual(t, got, want)
	})

}

func TestRemoveCategory(t *testing.T) {

	existingCategory := Category{ID: "1234", Name: "accommodation"}
	store := &StubCategoryStore{
		CategoryList{
			Categories: []Category{
				existingCategory,
			},
		},
	}
	server := NewCategoryServer(store)

	t.Run("check Content-Type header is application/json", func(t *testing.T) {
		body := []byte("foo")
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("Content-Type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("rejects request if body not json with id field", func(t *testing.T) {})

	t.Run("can remove existing category", func(t *testing.T) {})

	t.Run("cannot remove non-existent category", func(t *testing.T) {})

}

type StubCategoryStore struct {
	categories CategoryList
}

func (s *StubCategoryStore) ListCategories() CategoryList {
	return s.categories
}

func (s *StubCategoryStore) AddCategory(categoryName string) Category {
	newCat := Category{
		ID:   xid.New().String(),
		Name: categoryName,
	}

	s.categories.Categories = append(s.categories.Categories, newCat)

	return newCat
}

func (s *StubCategoryStore) RenameCategory(id, name string) Category {
	index := 0

	for i, c := range s.categories.Categories {
		if c.ID == id {
			index = i
			s.categories.Categories[index].Name = name
			break
		}
	}

	return s.categories.Categories[index]
}

func (s *StubCategoryStore) CategoryIdExists(categoryID string) bool {
	alreadyExists := false

	for _, c := range s.categories.Categories {
		if c.ID == categoryID {
			alreadyExists = true
		}
	}

	return alreadyExists
}

func (s *StubCategoryStore) CategoryNameExists(categoryName string) bool {
	alreadyExists := false

	for _, c := range s.categories.Categories {
		if c.Name == categoryName {
			alreadyExists = true
		}
	}

	return alreadyExists
}
