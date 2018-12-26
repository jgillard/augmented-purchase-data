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
			Category{ID: "abcdef", Name: "hostel", ParentID: "1234"},
			Category{ID: "ghijkm", Name: "apartment", ParentID: "1234"},
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
		assertStringsEqual(t, got.Categories[0].ID, "abcdef")
		assertStringsEqual(t, got.Categories[0].Name, "hostel")
		assertStringsEqual(t, got.Categories[0].ParentID, "1234")
		assertStringsEqual(t, got.Categories[1].ID, "ghijkm")
		assertStringsEqual(t, got.Categories[1].Name, "apartment")
		assertStringsEqual(t, got.Categories[1].ParentID, "1234")
	})
}

func TestGetCategory(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation", ParentID: ""},
			Category{ID: "2345", Name: "food and drink", ParentID: ""},
			Category{ID: "abcdef", Name: "hostel", ParentID: "1234"},
			Category{ID: "ghijkm", Name: "apartment", ParentID: "1234"},
		},
	}
	store := &stubCategoryStore{categoryList}
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

	t.Run("get category with children", func(t *testing.T) {
		req := newGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		got := unmarshallCategoryGetResponseFromBody(t, result.Body)
		assertStringsEqual(t, got.ID, categoryList.Categories[0].ID)
		assertStringsEqual(t, got.Name, categoryList.Categories[0].Name)
		assertStringsEqual(t, got.ParentID, categoryList.Categories[0].ParentID)

		accomodationChildren := []Category{
			categoryList.Categories[2],
			categoryList.Categories[3],
		}
		assertDeepEqual(t, got.Children, accomodationChildren)
	})

	t.Run("get category without children", func(t *testing.T) {
		req := newGetRequest(t, "/categories/2345")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		got := unmarshallCategoryGetResponseFromBody(t, result.Body)
		assertStringsEqual(t, got.ID, categoryList.Categories[1].ID)
		assertStringsEqual(t, got.Name, categoryList.Categories[1].Name)
		assertStringsEqual(t, got.ParentID, categoryList.Categories[1].ParentID)
		assertNumbersEqual(t, len(got.Children), 0)
	})
}

func TestAddCategory(t *testing.T) {

	stubCategories := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "existing category name", ParentID: ""},
			Category{ID: "2345", Name: "existing subcategory name", ParentID: "1234"},
		},
	}
	store := &stubCategoryStore{stubCategories}
	server := NewCategoryServer(store)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			input string
			want  int
		}{
			"invalid json":                     {input: `{"foo":""}`, want: http.StatusBadRequest},
			"name missing":                     {input: `{}`, want: http.StatusBadRequest},
			"duplicate name":                   {input: `{"name":"existing category name"}`, want: http.StatusConflict},
			"invalid name":                     {input: `{"name":"abc123!@£"}`, want: http.StatusUnprocessableEntity},
			"parentID missing":                 {input: `{"name":"valid name"}`, want: http.StatusBadRequest},
			"parentID doesn't exist":           {input: `{"parentID":"5678"}`, want: http.StatusUnprocessableEntity},
			"category would be >2 levels deep": {input: `{"name":"foo", "parentID":"2345"}`, want: http.StatusUnprocessableEntity},
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
		parentID := ""
		body := strings.NewReader(fmt.Sprintf(`{"name":"%s", "parentID":""}`, categoryName))
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
		assertStringsEqual(t, got.ParentID, parentID)

		// check the store has been modified
		got = store.categories.Categories[2]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)
		assertStringsEqual(t, got.ParentID, parentID)

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s", got.ID))
	})
}

func TestRenameCategory(t *testing.T) {

	stubCategories := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation", ParentID: ""},
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
		newCatName := "new category name"
		requestBody, err := json.Marshal(jsonName{Name: newCatName})
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

		responseBody := unmarshallCategoryFromBody(t, result.Body)

		renamedCategory := Category{ID: "1234", Name: newCatName, ParentID: ""}
		assertStringsEqual(t, responseBody.ID, renamedCategory.ID)
		assertStringsEqual(t, responseBody.Name, renamedCategory.Name)
		assertStringsEqual(t, responseBody.ParentID, renamedCategory.ParentID)

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

func (s *stubCategoryStore) GetCategory(id string) CategoryGetResponse {
	var category Category
	for _, c := range s.categories.Categories {
		if c.ID == id {
			category = c
		}
	}

	if category == (Category{}) {
		return CategoryGetResponse{}
	}

	var children []Category
	for _, c := range s.categories.Categories {
		if c.ParentID == category.ID {
			children = append(children, c)
		}
	}

	response := CategoryGetResponse{
		category,
		children,
	}
	return response
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

func (s *stubCategoryStore) categoryParentIDExists(parentID string) bool {
	exists := false

	for _, c := range s.categories.Categories {
		if c.ID == parentID {
			exists = true
		}
	}

	return exists
}

func (s *stubCategoryStore) getCategoryDepth(categoryID string) int {
	depth := 0

	for _, c := range s.categories.Categories {
		if c.ID == categoryID {
			// if already a subcategory
			if c.ParentID != "" {
				depth = 1
			}
		}
	}

	return depth
}
