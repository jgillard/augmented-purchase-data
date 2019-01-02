package transactioncategories

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListQuestionsForCategory(t *testing.T) {

	questionList := QuestionList{
		Questions: []Question{
			Question{ID: "1", Value: "how many nights?", CategoryID: "1234"},
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

		want := questionList

		assertDeepEqual(t, got, want)
		assertStringsEqual(t, got.Questions[0].ID, "1")
		assertStringsEqual(t, got.Questions[0].Value, "how many nights?")
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

type stubQuestionStore struct {
	questionList QuestionList
}

func (s *stubQuestionStore) ListQuestionsForCategory(categoryID string) QuestionList {
	return QuestionList{
		Questions: []Question{
			Question{ID: "1", Value: "how many nights?"},
		},
	}
}
