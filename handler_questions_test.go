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

func TestListQuestionsForCategory(t *testing.T) {

	questionList := QuestionList{
		Questions: []Question{
			Question{ID: "1", Title: "how many nights?", CategoryID: "1234", Type: "number"},
			Question{ID: "2", Title: "which meal?", CategoryID: "5678", Type: "string", Options: OptionList{
				{ID: "1", Title: "brekkie"},
			}},
		},
	}
	questionStore := &InMemoryQuestionStore{questionList}
	server := NewServer(nil, questionStore)

	t.Run("it returns a json question list for a category", func(t *testing.T) {
		req := newGetRequest(t, "/categories/1234/questions")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallQuestionListFromBody(t, body)
		want := questionList.Questions[0]
		assertDeepEqual(t, got.Questions[0], want)
	})

	t.Run("it returns a json question list for a category with options", func(t *testing.T) {
		req := newGetRequest(t, "/categories/5678/questions")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallQuestionListFromBody(t, body)
		want := questionList.Questions[1]
		assertDeepEqual(t, got.Questions[0], want)
	})

	t.Run("it returns an empty json question list when no questions exist for category", func(t *testing.T) {
		req := newGetRequest(t, "/categories/2345/questions")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallQuestionListFromBody(t, body)
		want := len(got.Questions)
		assertNumbersEqual(t, want, 0)
	})
}

func TestAddQuestion(t *testing.T) {

	t.Run("test failure responses & effect", func(t *testing.T) {
		categoryList := CategoryList{
			Categories: []Category{
				Category{ID: "1234", Name: "foo", ParentID: ""},
			},
		}
		questionList := QuestionList{
			Questions: []Question{
				Question{ID: "1", Title: "how many nights?", CategoryID: "1234", Type: "number"},
				Question{ID: "2", Title: "which meal?", CategoryID: "5678", Type: "string", Options: OptionList{
					{ID: "1", Title: "brekkie"},
				}},
			},
		}
		categoryStore := &InMemoryCategoryStore{categoryList}
		questionStore := &InMemoryQuestionStore{questionList}
		server := NewServer(categoryStore, questionStore)

		cases := map[string]struct {
			path  string
			input string
			want  int
		}{
			"invalid json":             {path: "/categories/1234/questions", input: `{"foo":""}`, want: http.StatusBadRequest},
			"title is empty":           {path: "/categories/1234/questions", input: `{"title":"", "type":"number"}`, want: http.StatusBadRequest},
			"title is duplicate":       {path: "/categories/1234/questions", input: `{"title":"how many nights?", "type":"number"}`, want: http.StatusConflict},
			"type is empty":            {path: "/categories/1234/questions", input: `{"title":"foo", "type":""}`, want: http.StatusBadRequest},
			"type doesn't exist":       {path: "/categories/1234/questions", input: `{"title":"foo", "type":"foo"}`, want: http.StatusBadRequest},
			"options is not list type": {path: "/categories/1234/questions", input: `{"title":"foo", "type":"string", "options":""}`, want: http.StatusBadRequest},
			"options had duplicate":    {path: "/categories/1234/questions", input: `{"title":"foo", "type":"string", "options":["foo", "foo"]}`, want: http.StatusBadRequest},
			"category doesn't exist":   {path: "/categories/5678/questions", input: `{"title":"foo", "type":"string"}`, want: http.StatusNotFound},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				requestBody := strings.NewReader(c.input)
				req := newPostRequest(t, c.path, requestBody)
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
				got := questionStore.questionList
				want := questionList
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("add type:number question", func(t *testing.T) {
		questionList := QuestionList{
			Questions: []Question{},
		}
		questionStore := &InMemoryQuestionStore{questionList}
		server := NewServer(nil, questionStore)

		categoryID := "1"
		title := "how many nights?"
		optionType := "number"

		jsonReq := fmt.Sprintf(`{"title":"%s", "type":"%s"}`, title, optionType)
		requestBody := strings.NewReader(jsonReq)
		path := fmt.Sprintf("/categories/%s/questions", categoryID)
		req := newPostRequest(t, path, requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallQuestionFromBody(t, body)

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		// Options should not be present in response
		if got.Options != nil {
			t.Fatalf("Options should not have been set, got %v", got.Options)
		}

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s/questions/%s", categoryID, got.ID))
	})

	t.Run("add type:string question without options", func(t *testing.T) {
		questionList := QuestionList{
			Questions: []Question{},
		}
		questionStore := &InMemoryQuestionStore{questionList}
		server := NewServer(nil, questionStore)

		categoryID := "1"
		title := "which meal?"
		optionType := "string"

		jsonReq := fmt.Sprintf(`{"title":"%s", "type":"%s"}`, title, optionType)
		requestBody := strings.NewReader(jsonReq)
		path := fmt.Sprintf("/categories/%s/questions", categoryID)
		req := newPostRequest(t, path, requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallQuestionFromBody(t, body)

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		// Options should be an empty list in response
		assertDeepEqual(t, got.Options, OptionList{})

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		assertDeepEqual(t, got.Options, OptionList{})

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s/questions/%s", categoryID, got.ID))
	})

	t.Run("add type:string question with options", func(t *testing.T) {
		questionList := QuestionList{
			Questions: []Question{},
		}
		questionStore := &InMemoryQuestionStore{questionList}
		server := NewServer(nil, questionStore)

		categoryID := "1"
		title := "which meal?"
		optionType := "string"
		options := []string{"brekkie", "lunch"}

		optionsJSON, err := json.Marshal(options)
		if err != nil {
			t.Fatal(err)
		}
		jsonReq := fmt.Sprintf(`{"title":"%s", "type":"%s", "options":%s}`, title, optionType, optionsJSON)

		requestBody := strings.NewReader(jsonReq)
		path := fmt.Sprintf("/categories/%s/questions", categoryID)
		req := newPostRequest(t, path, requestBody)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		got := unmarshallQuestionFromBody(t, body)

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		assertNumbersEqual(t, len(got.Options), 2)
		assertIsXid(t, got.Options[0].ID)
		assertIsXid(t, got.Options[1].ID)
		assertStringsEqual(t, got.Options[0].Title, options[0])
		assertStringsEqual(t, got.Options[1].Title, options[1])

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		assertNumbersEqual(t, len(got.Options), 2)
		assertIsXid(t, got.Options[0].ID)
		assertIsXid(t, got.Options[1].ID)
		assertStringsEqual(t, got.Options[0].Title, options[0])
		assertStringsEqual(t, got.Options[1].Title, options[1])

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s/questions/%s", categoryID, got.ID))
	})
}

func TestRenameQuestion(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "foo", ParentID: ""},
			Category{ID: "2345", Name: "bar", ParentID: ""},
		},
	}
	questionList := QuestionList{
		Questions: []Question{
			Question{ID: "1", Title: "how many nuggets?", CategoryID: "1234", Type: "number"},
			Question{ID: "2", Title: "how much nougat?", CategoryID: "1234", Type: "number"},
			Question{ID: "3", Title: "how much nougat?", CategoryID: "2345", Type: "number"},
		},
	}
	categoryStore := &InMemoryCategoryStore{categoryList}
	questionStore := &InMemoryQuestionStore{questionList}
	server := NewServer(categoryStore, questionStore)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			path  string
			input string
			want  int
		}{
			"invalid json":                         {path: "/categories/1234/questions/1", input: `{"foo":""}`, want: http.StatusBadRequest},
			"invalid title":                        {path: "/categories/1234/questions/1", input: `{"title":"foo/*!bar"}`, want: http.StatusUnprocessableEntity},
			"duplicate title":                      {path: "/categories/1234/questions/1", input: `{"title":"how much nougat?"}`, want: http.StatusConflict},
			"category doesn't exist":               {path: "/categories/5678/questions/1", input: `{"title":"irrelevant"}`, want: http.StatusNotFound},
			"ID not found":                         {path: "/categories/1234/questions/4", input: `{"title":"irrelevant"}`, want: http.StatusNotFound},
			"question does not belong to category": {path: "/categories/1234/questions/3", input: `{"title":"irrelevant"}`, want: http.StatusNotFound},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				requestBody := strings.NewReader(c.input)
				req := newPatchRequest(t, c.path, requestBody)
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
				got := questionStore.questionList
				want := questionList
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success responses & effect", func(t *testing.T) {
		newQuestionTitle := "whattup world?"
		requestBody, err := json.Marshal(jsonTitle{Title: newQuestionTitle})
		if err != nil {
			t.Fatal(err)
		}
		req := newPatchRequest(t, "/categories/1234/questions/1", bytes.NewReader(requestBody))
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		// check the response
		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)

		responseBody := unmarshallQuestionFromBody(t, body)

		renamedQuestion := Question{ID: "1", Title: newQuestionTitle, CategoryID: "1234", Type: "number"}
		assertStringsEqual(t, responseBody.ID, renamedQuestion.ID)
		assertStringsEqual(t, responseBody.Title, renamedQuestion.Title)
		assertStringsEqual(t, responseBody.CategoryID, renamedQuestion.CategoryID)
		assertStringsEqual(t, responseBody.Type, renamedQuestion.Type)

		// check the store is updated
		got := questionStore.questionList.Questions[0].Title
		want := renamedQuestion.Title
		assertStringsEqual(t, got, want)
	})
}

func TestRemoveQuestion(t *testing.T) {

	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "foo", ParentID: ""},
			Category{ID: "2345", Name: "bar", ParentID: ""},
		},
	}
	questionList := QuestionList{
		Questions: []Question{
			Question{ID: "1", Title: "how many nuggets?", CategoryID: "1234", Type: "number"},
		},
	}
	categoryStore := &InMemoryCategoryStore{categoryList}
	questionStore := &InMemoryQuestionStore{questionList}
	server := NewServer(categoryStore, questionStore)

	t.Run("test failure responses & effect", func(t *testing.T) {
		cases := map[string]struct {
			path string
			want int
		}{
			"category doesn't exist":              {path: "/categories/5678/questions/2", want: http.StatusNotFound},
			"question doesn't exist":              {path: "/categories/1234/questions/2", want: http.StatusNotFound},
			"question doesn't belong to category": {path: "/categories/2345/questions/1", want: http.StatusNotFound},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				req := newDeleteRequest(t, c.path)
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
				got := questionStore.questionList
				want := questionList
				assertDeepEqual(t, got, want)
			})
		}
	})

	t.Run("test success response & effect", func(t *testing.T) {
		req := newDeleteRequest(t, "/categories/1234/questions/1")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()
		body := readBodyBytes(t, result.Body)

		// check response
		assertStatusCode(t, result.StatusCode, http.StatusNoContent)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
		assertBodyIsJSON(t, body)
		assertBodyEmptyJSON(t, body)

		// check store is updated
		got := len(questionStore.questionList.Questions)
		want := 0
		assertNumbersEqual(t, got, want)
	})
}
