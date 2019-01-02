package transactioncategories

import (
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

	questionList := QuestionList{
		Questions: []Question{},
	}
	questionStore := &stubQuestionStore{questionList}
	server := NewServer(nil, questionStore)

	t.Run("add type:number question", func(t *testing.T) {
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

		// check the store has been modified
		got = questionStore.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Value, value)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, optionType)

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
