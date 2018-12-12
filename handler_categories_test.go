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

	req := NewGetRequest(t, "/categories")
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)
	result := res.Result()

	t.Run("success status code", func(t *testing.T) {
		got := result.StatusCode
		want := http.StatusOK
		assertNumbersEqual(t, got, want)
	})

	t.Run("content-type header", func(t *testing.T) {
		got := result.Header.Get("content-type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("response format", func(t *testing.T) {
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

	t.Run("content-type header", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("content-type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("success status code", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusOK
		assertNumbersEqual(t, got, want)
	})

	t.Run("not-found status code", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/5678")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusNotFound
		assertNumbersEqual(t, got, want)
	})

	t.Run("success response format", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/1234")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

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

	t.Run("not-found response format", func(t *testing.T) {
		req := NewGetRequest(t, "/categories/5678")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		want := "{}"
		assertStringsEqual(t, got, want)
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

	t.Run("content-type header", func(t *testing.T) {
		body := strings.NewReader("foo")
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("content-type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("success status code", func(t *testing.T) {
		categoryName := "new category name"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusCreated
		assertNumbersEqual(t, got, want)
	})

	t.Run("duplicate failure status code", func(t *testing.T) {
		categoryName := "existing category name"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusConflict
		assertNumbersEqual(t, got, want)
	})

	t.Run("invalid name failure status code", func(t *testing.T) {
		categoryName := "abc123!@£"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusUnprocessableEntity
		assertNumbersEqual(t, got, want)

		responseBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		gotBody := string(responseBody)
		wantBody := "{}"
		assertStringsEqual(t, gotBody, wantBody)

	})

	t.Run("success response format", func(t *testing.T) {
		categoryName := "new category name"
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

		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)
	})

	t.Run("failure response format", func(t *testing.T) {
		categoryName := "existing category name"
		requestBody := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		responseBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(responseBody)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("success store updated", func(t *testing.T) {
		categoryName := "new category name"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := store.categories.Categories[1]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Name, categoryName)
	})

	t.Run("failure store unmodified", func(t *testing.T) {
		categoryName := "existing category name"
		body := strings.NewReader(categoryName)
		req := NewPostRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := store.categories.Categories
		want := stubCategories
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got '%v' wanted '%v'", got, want)
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

	t.Run("content-type header", func(t *testing.T) {
		body := []byte("foo")
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("content-type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("success status code", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "new category name"}
		body, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusOK
		assertNumbersEqual(t, got, want)
	})

	t.Run("invalid json failure status code", func(t *testing.T) {
		body := []byte("foo")
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusBadRequest
		assertNumbersEqual(t, got, want)
	})

	t.Run("invalid name failure status code", func(t *testing.T) {
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

	t.Run("not-found failure status code", func(t *testing.T) {
		renamedCategory := Category{ID: "5678", Name: "irrelevant"}
		body, err := json.Marshal(renamedCategory)
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

	t.Run("success response format", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "new category name"}
		requestBody, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		want := string(requestBody)
		assertStringsEqual(t, got, want)
	})

	t.Run("invalid json failure reponse format", func(t *testing.T) {
		body := []byte("foo")
		req := NewPutRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("invalid name failure reponse format", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "foo/*!bar"}
		requestBody, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("not-found failure reponse format", func(t *testing.T) {
		renamedCategory := Category{ID: "5678", Name: "irrelevant"}
		requestBody, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(body)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("success store updated", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "new category name"}
		requestBody, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := store.categories.Categories[0].Name
		want := renamedCategory.Name
		assertStringsEqual(t, got, want)
	})

	t.Run("failure store unmodified", func(t *testing.T) {
		renamedCategory := Category{ID: "1234", Name: "foo/*!bar"}
		requestBody, err := json.Marshal(renamedCategory)
		if err != nil {
			t.Fatal(err)
		}
		req := NewPutRequest(t, "/categories", requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := store.categories.Categories
		want := existingCategory
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got '%v' wanted '%v'", got, want)
		}
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

	t.Run("content-type header", func(t *testing.T) {
		body := strings.NewReader("{}")
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.Header.Get("content-type")
		want := jsonContentType
		assertStringsEqual(t, got, want)
	})

	t.Run("success status code", func(t *testing.T) {
		body := strings.NewReader(`{"id":"1234"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusOK
		assertNumbersEqual(t, got, want)
	})

	t.Run("invalid category json failure status code", func(t *testing.T) {
		body := strings.NewReader(`{"foo":"bar"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusBadRequest
		assertNumbersEqual(t, got, want)
	})

	t.Run("not-found failure status code", func(t *testing.T) {
		body := strings.NewReader(`{"id":"5678"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		got := result.StatusCode
		want := http.StatusNotFound
		assertNumbersEqual(t, got, want)
	})

	t.Run("success reponse format", func(t *testing.T) {
		body := strings.NewReader(`{"id":"1234"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		responseBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(responseBody)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("any failure reponse format", func(t *testing.T) {
		body := strings.NewReader(`{"id":"5678"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		responseBody, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}

		got := string(responseBody)
		want := "{}"
		assertStringsEqual(t, got, want)
	})

	t.Run("success store updated", func(t *testing.T) {
		body := strings.NewReader(`{"id":"1234"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := len(store.categories.Categories)
		want := 0
		assertNumbersEqual(t, got, want)
	})

	t.Run("failure store unmodified", func(t *testing.T) {
		body := strings.NewReader(`{"id":"5678"}`)
		req := NewDeleteRequest(t, "/categories", body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		got := len(store.categories.Categories)
		want := 1
		assertNumbersEqual(t, got, want)
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
