package handlers

import (
	"bytes"
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

	t.Run("it returns a json category list", func(t *testing.T) {
		categoryList := CategoryList{
			Categories: []Category{
				Category{ID: "1234", Name: "accommodation"},
				Category{ID: "5678", Name: "food and drink"},
			},
		}
		store := &StubCategoryStore{categoryList}
		server := NewCategoryServer(store)

		req := NewGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("content-type"), jsonContentType)

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
}

func TestGetCategory(t *testing.T) {

	stubCategory := Category{ID: "1234", Name: "accommodation"}
	store := &StubCategoryStore{
		CategoryList{
			Categories: []Category{
				stubCategory,
			},
		},
	}
	server := NewCategoryServer(store)

	t.Run("not-found failure reponse", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/5678")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusNotFound)
		assertContentType(t, result.Header.Get("content-type"), jsonContentType)

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("success response", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("content-type"), jsonContentType)

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var got Category

		err = json.Unmarshal(body, &got)
		if err != nil {
			t.Fatal(err)
		}

		assertStringsEqual(t, got.ID, stubCategory.ID)
		assertStringsEqual(t, got.Name, stubCategory.Name)
	})
}

func TestAddCategory(t *testing.T) {

	stubCategories := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "existing category name"},
		},
	}
	store := &StubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := []struct {
			name  string
			value string
			want  int
		}{
			{name: "duplicate name", value: "existing category name", want: http.StatusConflict},
			{name: "invalid name", value: "abc123!@£", want: http.StatusUnprocessableEntity},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				body := strings.NewReader(c.value)
				req := NewPostRequest(t, "/categories", body)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("content-type"), jsonContentType)

				responseBody, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatal(err)
				}

				assertBodyString(t, string(responseBody), "{}")

				// check the store is unmodified
				got := store.categories
				want := stubCategories
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got '%v' wanted '%v'", got, want)
				}
			})
		}
	})

	t.Run("test success response & effect", func(t *testing.T) {
		categoryName := "new category name"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("content-type"), jsonContentType)

		bodyBytes, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		var got Category

		err = json.Unmarshal(bodyBytes, &got)
		if err != nil {
			t.Fatal(err)
		}

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)

		// check the store has been modified
		got = store.categories.Categories[1]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)
	})
}

func TestRenameCategory(t *testing.T) {

	stubCategories := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation"},
		},
	}
	store := &StubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := []struct {
			name  string
			value string
			want  int
		}{
			{name: "invalid json", value: "foo", want: http.StatusBadRequest},
			{name: "invalid name", value: categoryToString(t, Category{ID: "1234", Name: "foo/*!bar"}), want: http.StatusUnprocessableEntity},
			{name: "name not found", value: categoryToString(t, Category{ID: "5678", Name: "irrelevant"}), want: http.StatusNotFound},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				body := strings.NewReader(c.value)
				req := NewPutRequest(t, "/categories", body)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("content-type"), jsonContentType)

				responseBody, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatal(err)
				}

				assertBodyString(t, string(responseBody), "{}")

				// check the store is unmodified
				got := store.categories
				want := stubCategories
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got '%v' wanted '%v'", got, want)
				}
			})
		}
	})

	t.Run("test success responses & effect", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "new category name"}
		requestBody, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", bytes.NewReader(requestBody))
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("content-type"), jsonContentType)

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		assertBodyString(t, string(body), string(requestBody))

		// check the store is updated
		got := store.categories.Categories[0].Name
		want := renamedCategory.Name
		assertStringsEqual(t, got, want)
	})
}

func TestRemoveCategory(t *testing.T) {

	existingCategory := Category{ID: "1234", Name: "accommodation"}
	stubCategories := CategoryList{
		Categories: []Category{
			existingCategory,
		},
	}
	store := &StubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := []struct {
			name  string
			value string
			want  int
		}{
			{name: "invalid category json", value: `{"foo":"bar"}`, want: http.StatusBadRequest},
			{name: "category not found", value: `{"id":"5678"}`, want: http.StatusNotFound},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				body := strings.NewReader(c.value)
				req := NewDeleteRequest(t, "/categories", body)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("content-type"), jsonContentType)

				responseBody, err := ioutil.ReadAll(result.Body)
				if err != nil {
					t.Fatal(err)
				}

				assertBodyString(t, string(responseBody), "{}")

				// check the store is unmodified
				got := store.categories
				want := stubCategories
				if !reflect.DeepEqual(got, want) {
					t.Errorf("got '%v' wanted '%v'", got, want)
				}
			})
		}
	})

	t.Run("test success responses & effect", func(t *testing.T) {
		body := strings.NewReader(`{"id":"1234"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("content-type"), jsonContentType)

		responseBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(responseBody)
		want := "{}"
		assertStringsEqual(t, got, want)

		// check store is updated
		gotStore := store.categories
		wantStore := CategoryList{}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got '%v' wanted '%v'", gotStore, wantStore)
		}
	})
}

func TestIsValidCategoryName(t *testing.T) {
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
			got := IsValidCategoryName(value)
			want := false

			if got != want {
				t.Errorf("'%s' should be treated as an invalid category name", value)
			}
		})
	}
}

type StubCategoryStore struct {
	categories CategoryList
}

func (s *StubCategoryStore) ListCategories() CategoryList {
	return s.categories
}

func (s *StubCategoryStore) GetCategory(id string) Category {
	var category Category

	for _, c := range s.categories.Categories {
		if c.ID == id {
			category = c
		}
	}
	return category
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

func (s *StubCategoryStore) DeleteCategory(id string) {
	index := 0

	for i, c := range s.categories.Categories {
		if c.ID == id {
			index = i
			break
		}
	}

	s.categories.Categories = append(s.categories.Categories[:index], s.categories.Categories[index+1:]...)
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
