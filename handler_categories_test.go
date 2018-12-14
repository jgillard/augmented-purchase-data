package transactioncategories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
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
	store := &stubCategoryStore{categoryList}
	server := NewCategoryServer(store)

	t.Run("it returns a json category list", func(t *testing.T) {
		req := newGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		got := unmarshallCategoryListFromBody(t, result.Body)
		want := categoryList
		assertDeepEqual(t, got, want)
	})
}

func TestGetCategory(t *testing.T) {

	stubCategory := Category{ID: "1234", Name: "accommodation"}
	store := &stubCategoryStore{
		CategoryList{
			Categories: []Category{
				stubCategory,
			},
		},
	}
	server := NewCategoryServer(store)

	t.Run("not-found failure reponse", func(t *testing.T) {
		req := newGetRequest(t, "/categories/5678")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusNotFound)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyEmpty(t, result.Body)
	})

	t.Run("success response", func(t *testing.T) {
		req := newGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		got := unmarshallCategoryFromBody(t, result.Body)
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
	store := &stubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			input string
			want  int
		}{
			"invalid json":   {input: `{"foo":""}`, want: http.StatusBadRequest},
			"duplicate name": {input: `{"name":"existing category name"}`, want: http.StatusConflict},
			"invalid name":   {input: `{"name":"abc123!@£"}`, want: http.StatusUnprocessableEntity},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				body := strings.NewReader(c.input)
				req := newPostRequest(t, "/categories", body)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyEmpty(t, result.Body)

				// check the store is unmodified
				got := store.categories
				want := stubCategories
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success response & effect", func(t *testing.T) {
		categoryName := "new category name"
		body := strings.NewReader(fmt.Sprintf(`{"name":"%s"}`, categoryName))
		req := newPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		got := unmarshallCategoryFromBody(t, result.Body)

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)

		// check the store has been modified
		got = store.categories.Categories[1]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s", got.ID))
	})
}

func TestRenameCategory(t *testing.T) {

	stubCategories := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation"},
		},
	}
	store := &stubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			ID   string
			body string
			want int
		}{
			"invalid json":   {ID: "1234", body: `{"foo":""}`, want: http.StatusBadRequest},
			"invalid name":   {ID: "1234", body: `{"name":"foo/*!bar"}`, want: http.StatusUnprocessableEntity},
			"duplicate name": {ID: "1234", body: `{"name":"accommodation"}`, want: http.StatusConflict},
			"ID not found":   {ID: "5678", body: `{"name":"irrelevant"}`, want: http.StatusNotFound},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				body := strings.NewReader(c.body)
				req := newPatchRequest(t, fmt.Sprintf("/categories/%s", c.ID), body)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyEmpty(t, result.Body)

				// check the store is unmodified
				got := store.categories
				want := stubCategories
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success responses & effect", func(t *testing.T) {
		newCatName := jsonName{Name: "new category name"}
		requestBody, err := json.Marshal(newCatName)
		if err != nil {
			t.Fatal(err)
		}
		req := newPatchRequest(t, "/categories/1234", bytes.NewReader(requestBody))
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		responseBody := readBodyBytes(t, result.Body)

		renamedCategory := Category{ID: "1234", Name: "new category name"}
		renamedCategoryBytes, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		assertBodyString(t, string(responseBody), string(renamedCategoryBytes))

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
	store := &stubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			input string
			want  int
		}{
			"category not found": {input: "5678", want: http.StatusNotFound},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				req := newDeleteRequest(t, fmt.Sprintf("/categories/%s", c.input))
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyEmpty(t, result.Body)

				// check the store is unmodified
				got := store.categories
				want := stubCategories
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success response & effect", func(t *testing.T) {
		req := newDeleteRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check response
		assertStatusCode(t, result.StatusCode, http.StatusNoContent)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyEmpty(t, result.Body)

		// check store is updated
		got := len(store.categories.Categories)
		want := 0
		assertNumbersEqual(t, got, want)
	})
}

func testIsValidCategoryName(t *testing.T) {
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

type stubCategoryStore struct {
	categories CategoryList
}

func (s *stubCategoryStore) ListCategories() CategoryList {
	return s.categories
}

func (s *stubCategoryStore) GetCategory(id string) Category {
	var category Category

	for _, c := range s.categories.Categories {
		if c.ID == id {
			category = c
		}
	}
	return category
}

func (s *stubCategoryStore) AddCategory(categoryName string) Category {
	newCat := Category{
		ID:   xid.New().String(),
		Name: categoryName,
	}

	s.categories.Categories = append(s.categories.Categories, newCat)

	return newCat
}

func (s *stubCategoryStore) RenameCategory(id, name string) Category {
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

func (s *stubCategoryStore) DeleteCategory(id string) {
	index := 0

	for i, c := range s.categories.Categories {
		if c.ID == id {
			index = i
			break
		}
	}

	s.categories.Categories = append(s.categories.Categories[:index], s.categories.Categories[index+1:]...)
}

func (s *stubCategoryStore) categoryIDExists(categoryID string) bool {
	alreadyExists := false

	for _, c := range s.categories.Categories {
		if c.ID == categoryID {
			alreadyExists = true
		}
	}

	return alreadyExists
}

func (s *stubCategoryStore) categoryNameExists(categoryName string) bool {
	alreadyExists := false

	for _, c := range s.categories.Categories {
		if c.Name == categoryName {
			alreadyExists = true
		}
	}

	return alreadyExists
}
