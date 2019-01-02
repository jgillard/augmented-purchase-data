package transactioncategories

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rs/xid"
)

func TestListQuestionsForCategory(t *testing.T) {

	questionList := QuestionList{
		Questions: []Question{
			Question{ID: "1", Value: "how many nights?", CategoryID: "1234", Type: "number"},
			Question{ID: "2", Value: "which meal?", CategoryID: "5678", Type: "string", Options: OptionList{
				{ID: "1", Value: "brekkie"},
			}},
		},
	}
	questionStore := &stubQuestionStore{questionList}
	server := NewServer(nil, questionStore)

	t.Run("it returns a json question list for a category", func(t *testing.T) {
		req := newGetRequest(t, "/categories/1234/questions")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		bodyBytes := readBodyBytes(t, result.Body)

		var got QuestionList

		err := json.Unmarshal(bodyBytes, &got)
		// check for syntax error or type mismatch
		if err != nil {
			t.Fatal(err)
		}

		want := questionList.Questions[0]

		assertDeepEqual(t, got.Questions[0], want)
	})

	t.Run("it returns a json question list for a category with options", func(t *testing.T) {
		req := newGetRequest(t, "/categories/5678/questions")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		bodyBytes := readBodyBytes(t, result.Body)

		var got QuestionList

		err := json.Unmarshal(bodyBytes, &got)
		// check for syntax error or type mismatch
		if err != nil {
			t.Fatal(err)
		}

		want := questionList.Questions[1]

		assertDeepEqual(t, got.Questions[0], want)
	})

	t.Run("it returns an empty json question list when no questions exist for category", func(t *testing.T) {
		req := newGetRequest(t, "/categories/2345/questions")
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusOK)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		bodyBytes := readBodyBytes(t, result.Body)

		var got QuestionList

		err := json.Unmarshal(bodyBytes, &got)
		// check for syntax error or type mismatch
		if err != nil {
			t.Fatal(err)
		}

		assertNumbersEqual(t, len(got.Questions), 0)
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
				Question{ID: "1", Value: "how many nights?", CategoryID: "1234", Type: "number"},
				Question{ID: "2", Value: "which meal?", CategoryID: "5678", Type: "string", Options: OptionList{
					{ID: "1", Value: "brekkie"},
				}},
			},
		}
		categoryStore := &stubCategoryStore{categoryList}
		questionStore := &stubQuestionStore{questionList}
		server := NewServer(categoryStore, questionStore)

		cases := map[string]struct {
			path  string
			input string
			want  int
		}{
			"invalid json":             {path: "/categories/1234/questions", input: `{"foo":""}`, want: http.StatusBadRequest},
			"value is empty":           {path: "/categories/1234/questions", input: `{"value":"", "type":"number"}`, want: http.StatusBadRequest},
			"value is duplicate":       {path: "/categories/1234/questions", input: `{"value":"how many nights?", "type":"number"}`, want: http.StatusConflict},
			"type is empty":            {path: "/categories/1234/questions", input: `{"value":"foo", "type":""}`, want: http.StatusBadRequest},
			"type doesn't exist":       {path: "/categories/1234/questions", input: `{"value":"foo", "type":"foo"}`, want: http.StatusBadRequest},
			"options is not list type": {path: "/categories/1234/questions", input: `{"value":"foo", "type":"string", "options":""}`, want: http.StatusBadRequest},
			"options had duplicate":    {path: "/categories/1234/questions", input: `{"value":"foo", "type":"string", "options":["foo", "foo"]}`, want: http.StatusBadRequest},
			"category doesn't exist":   {path: "/categories/5678/questions", input: `{"value":"foo", "type":"string"}`, want: http.StatusNotFound},
		}

		for name, c := range cases {
			t.Run(name, func(t *testing.T) {
				body := strings.NewReader(c.input)
				req := newPostRequest(t, c.path, body)
				res := httptest.NewRecorder()

				server.ServeHTTP(res, req)
				result := res.Result()

				// check the response
				assertStatusCode(t, result.StatusCode, c.want)
				assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)
				assertBodyEmpty(t, result.Body)

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
		questionStore := &stubQuestionStore{questionList}
		server := NewServer(nil, questionStore)

		categoryID := "1"
		value := "how many nights?"
		optionType := "number"

		jsonReq := fmt.Sprintf(`{"value":"%s", "type":"%s"}`, value, optionType)
		body := strings.NewReader(jsonReq)
		path := fmt.Sprintf("/categories/%s/questions", categoryID)
		req := newPostRequest(t, path, body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		bodyBytes := readBodyBytes(t, result.Body)

		var got Question

		err := json.Unmarshal(bodyBytes, &got)
		// check for syntax error or type mismatch
		if err != nil {
			t.Fatal(err)
		}

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		// Options should not be present in response
		if got.Options != nil {
			t.Fatalf("Options should not have been set, got %v", got.Options)
		}

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s/questions/%s", categoryID, got.ID))
	})

	t.Run("add type:string question without options", func(t *testing.T) {
		questionList := QuestionList{
			Questions: []Question{},
		}
		questionStore := &stubQuestionStore{questionList}
		server := NewServer(nil, questionStore)

		categoryID := "1"
		value := "which meal?"
		optionType := "string"

		jsonReq := fmt.Sprintf(`{"value":"%s", "type":"%s"}`, value, optionType)
		body := strings.NewReader(jsonReq)
		path := fmt.Sprintf("/categories/%s/questions", categoryID)
		req := newPostRequest(t, path, body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		bodyBytes := readBodyBytes(t, result.Body)

		var got Question

		err := json.Unmarshal(bodyBytes, &got)
		// check for syntax error or type mismatch
		if err != nil {
			t.Fatal(err)
		}

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		// Options should be an empty list in response
		assertDeepEqual(t, got.Options, OptionList{})

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
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
		questionStore := &stubQuestionStore{questionList}
		server := NewServer(nil, questionStore)

		categoryID := "1"
		value := "which meal?"
		optionType := "string"
		options := []string{"brekkie", "lunch"}

		optionsJSON, err := json.Marshal(options)
		if err != nil {
			t.Fatal(err)
		}
		jsonReq := fmt.Sprintf(`{"value":"%s", "type":"%s", "options":%s}`, value, optionType, optionsJSON)

		body := strings.NewReader(jsonReq)
		path := fmt.Sprintf("/categories/%s/questions", categoryID)
		req := newPostRequest(t, path, body)
		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)
		result := res.Result()

		assertStatusCode(t, result.StatusCode, http.StatusCreated)
		assertContentType(t, result.Header.Get("Content-Type"), jsonContentType)

		bodyBytes := readBodyBytes(t, result.Body)

		var got Question

		err = json.Unmarshal(bodyBytes, &got)
		// check for syntax error or type mismatch
		if err != nil {
			t.Fatal(err)
		}

		// check the response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		assertNumbersEqual(t, len(got.Options), 2)
		assertIsXid(t, got.Options[0].ID)
		assertIsXid(t, got.Options[1].ID)
		assertStringsEqual(t, got.Options[0].Value, options[0])
		assertStringsEqual(t, got.Options[1].Value, options[1])

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)
		assertNumbersEqual(t, len(got.Options), 2)
		assertIsXid(t, got.Options[0].ID)
		assertIsXid(t, got.Options[1].ID)
		assertStringsEqual(t, got.Options[0].Value, options[0])
		assertStringsEqual(t, got.Options[1].Value, options[1])

		// get ID from store and check that's in returned Location header
		assertStringsEqual(t, result.Header.Get("Location"), fmt.Sprintf("/categories/%s/questions/%s", categoryID, got.ID))
	})
}

type stubQuestionStore struct {
	questionList QuestionList
}

func (s *stubQuestionStore) ListQuestionsForCategory(categoryID string) QuestionList {
	var questionList QuestionList
	for _, q := range s.questionList.Questions {
		if q.CategoryID == categoryID {
			questionList.Questions = append(questionList.Questions, q)
		}
	}
	return questionList
}

func (s *stubQuestionStore) AddQuestion(categoryID string, q QuestionPostRequest) Question {
	question := Question{
		ID:         xid.New().String(),
		Value:      q.Value,
		CategoryID: categoryID,
		Type:       q.Type,
	}

	if q.Type == "string" {
		question.Options = OptionList{}
		for _, value := range q.Options {
			option := Option{
				ID:    xid.New().String(),
				Value: value,
			}
			question.Options = append(question.Options, option)
		}
	}

	s.questionList.Questions = append(s.questionList.Questions, question)

	return question
}

func (s *stubQuestionStore) questionValueExists(categoryID, questionValue string) bool {
	alreadyExists := false
	for _, q := range s.questionList.Questions {
		if q.CategoryID == categoryID {
			if q.Value == questionValue {
				alreadyExists = true
			}
		}
	}
	return alreadyExists
}
