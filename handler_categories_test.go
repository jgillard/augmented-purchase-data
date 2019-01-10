package transactioncategories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListCategories(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "abcdef", Name: "hostel", ParentID: "1234"},
			Category{ID: "ghijkm", Name: "apartment", ParentID: "1234"},
		},
	}
	store := NewInMemoryCategoryStore(categoryList)
	server := NewServer(store, nil)

	t.Run("it returns a json category list", func(t *testing.T) {
		req := newGetRequest(t, "/categories")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallCategoryListFromBody(t, body)
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
	store := NewInMemoryCategoryStore(categoryList)
	server := NewServer(store, nil)

	t.Run("not-found failure reponse", func(t *testing.T) {
		req := newGetRequest(t, "/categories/5678")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusNotFound)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)
		assertBodyErrorTitle(t, body, "categoryID not found")
	})

	t.Run("get category with children", func(t *testing.T) {
		req := newGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallCategoryGetResponseFromBody(t, body)
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
		body := readBodyBytes(t, result.Body)

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallCategoryGetResponseFromBody(t, body)
		assertStringsEqual(t, got.ID, categoryList.Categories[1].ID)
		assertStringsEqual(t, got.Name, categoryList.Categories[1].Name)
		assertStringsEqual(t, got.ParentID, categoryList.Categories[1].ParentID)
		assertNumbersEqual(t, len(got.Children), 0)
	})
}

func TestAddCategory(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "existing category name", ParentID: ""},
			Category{ID: "2345", Name: "existing subcategory name", ParentID: "1234"},
		},
	}
	store := NewInMemoryCategoryStore(categoryList)
	server := NewServer(store, nil)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			input      string
			want       int
			errorTitle string
		}{
			"invalid json": {
				input:      `"foo"`,
				want:       http.StatusBadRequest,
				errorTitle: ErrorInvalidJSON,
			},
			"name missing": {
				input:      `{}`,
				want:       http.StatusBadRequest,
				errorTitle: ErrorFieldMissing,
			},
			"duplicate name": {
				input:      `{"name":"existing category name"}`,
				want:       http.StatusConflict,
				errorTitle: ErrorDuplicateCategoryName,
			},
			"invalid name": {
				input:      `{"name":"abc123!@£"}`,
				want:       http.StatusUnprocessableEntity,
				errorTitle: ErrorInvalidCategoryName,
			},
			"parentID missing": {
				input:      `{"name":"valid name"}`,
				want:       http.StatusBadRequest,
				errorTitle: ErrorFieldMissing,
			},
			"parentID doesn't exist": {
				input:      `{"name":"foo", "parentID":"5678"}`,
				want:       http.StatusUnprocessableEntity,
				errorTitle: ErrorParentIDNotFound,
			},
			"category would be >2 levels deep": {
				input:      `{"name":"foo", "parentID":"2345"}`,
				want:       http.StatusUnprocessableEntity,
				errorTitle: ErrorCategoryTooNested,
			},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				requestBody := strings.NewReader(c.input)
				req := newPostRequest(t, "/categories", requestBody)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()
				body := readBodyBytes(t, result.Body)

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyIsJSON(t, body)
				assertBodyErrorTitle(t, body, c.errorTitle)

				// check the store is unmodified
				got := store.categories
				want := categoryList
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success response & effect without parentID", func(t *testing.T) {
		categoryName := "new category name"
		parentID := ""
		requestBody := strings.NewReader(fmt.Sprintf(`{"name":"%s", "parentID":""}`, categoryName))
		req := newPostRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallCategoryFromBody(t, body)

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

	t.Run("test success response & effect with parentID", func(t *testing.T) {
		categoryName := "another new category name"
		parentID := "1234"
		requestBody := strings.NewReader(fmt.Sprintf(`{"name":"%s", "parentID":"%s"}`, categoryName, parentID))
		req := newPostRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallCategoryFromBody(t, body)

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)
		assertStringsEqual(t, got.ParentID, parentID)

		// check the store has been modified
		got = store.categories.Categories[3]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)
		assertStringsEqual(t, got.ParentID, parentID)

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s", got.ID))
	})
}

func TestRenameCategory(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation", ParentID: ""},
		},
	}
	store := NewInMemoryCategoryStore(categoryList)
	server := NewServer(store, nil)

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
				requestBody := strings.NewReader(c.body)
				req := newPatchRequest(t, fmt.Sprintf("/categories/%s", c.ID), requestBody)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()
				body := readBodyBytes(t, result.Body)

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyIsJSON(t, body)
				assertBodyEmptyJSON(t, body)

				// check the store is unmodified
				got := store.categories
				want := categoryList
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
		body := readBodyBytes(t, result.Body)

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		responseBody := unmarshallCategoryFromBody(t, body)

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
	categoryList := CategoryList{
		Categories: []Category{
			existingCategory,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)
	server := NewServer(store, nil)

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
				body := readBodyBytes(t, result.Body)

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyIsJSON(t, body)
				assertBodyEmptyJSON(t, body)

				// check the store is unmodified
				got := store.categories
				want := categoryList
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success response & effect", func(t *testing.T) {
		req := newDeleteRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		// check response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)
		assertBodyJSONIsStatus(t, body, "deleted")

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
