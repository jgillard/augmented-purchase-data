package transactioncategories

import "testing"

func TestInMemoryQuestionStore(t *testing.T) {
	got := NewInMemoryQuestionStore(QuestionList{})
	want := &InMemoryQuestionStore{}
	assertDeepEqual(t, got, want)
}

func TestInMemoryQuestionStore_ListQuestionsForCategory(t *testing.T) {
	questionList := QuestionList{
		Questions: []Question{
			Question{ID: "1", Title: "how many nights?", CategoryID: "1234", Type: "number"},
		},
	}
	store := NewInMemoryQuestionStore(questionList)

	got := store.ListQuestionsForCategory("1234")
	want := questionList
	assertDeepEqual(t, got, want)
}

func TestInMemoryQuestionStore_AddQuestion(t *testing.T) {
	questionList := QuestionList{}
	store := NewInMemoryQuestionStore(questionList)

	categoryID := "1234"

	t.Run("question with options", func(t *testing.T) {
		question := QuestionPostRequest{
			Title:   "foo",
			Type:    "string",
			Options: &[]string{"bar"},
		}

		got := store.AddQuestion(categoryID, question)

		// assert response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, question.Title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, question.Type)
		assertIsXid(t, got.Options[0].ID)
		assertStringsEqual(t, got.Options[0].Title, (*question.Options)[0])

		// assert store
		got = store.questionList.Questions[0]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, question.Title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, question.Type)
		assertIsXid(t, got.Options[0].ID)
		assertStringsEqual(t, got.Options[0].Title, (*question.Options)[0])
	})

	t.Run("question without options", func(t *testing.T) {
		question := QuestionPostRequest{
			Title:   "foo",
			Type:    "number",
			Options: nil,
		}

		got := store.AddQuestion(categoryID, question)

		// assert response
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, question.Title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, question.Type)
		if got.Options != nil {
			t.Fatal("options should not be present")
		}

		// assert store
		got = store.questionList.Questions[1]
		assertIsXid(t, got.ID)
		assertStringsEqual(t, got.Title, question.Title)
		assertStringsEqual(t, got.CategoryID, categoryID)
		assertStringsEqual(t, got.Type, question.Type)
		if got.Options != nil {
			t.Fatal("options should not be present")
		}
	})
}

func TestInMemoryQuestionStore_RenameQuestion(t *testing.T) {
	question := Question{ID: "1", Title: "how many nights?", CategoryID: "1234", Type: "number"}
	questionList := QuestionList{
		Questions: []Question{
			question,
		},
	}
	store := NewInMemoryQuestionStore(questionList)

	newTitle := "foobar"

	got := store.RenameQuestion(question.ID, newTitle)

	// assert response
	assertStringsEqual(t, got.ID, question.ID)
	assertStringsEqual(t, got.Title, newTitle)

	// assert store
	got = store.questionList.Questions[0]
	assertStringsEqual(t, got.ID, question.ID)
	assertStringsEqual(t, got.Title, newTitle)
}

func TestInMemoryQuestionStore_DeleteQuestion(t *testing.T) {
	question := Question{ID: "1", Title: "how many nights?", CategoryID: "1234", Type: "number"}
	questionList := QuestionList{
		Questions: []Question{
			question,
		},
	}
	store := NewInMemoryQuestionStore(questionList)

	store.DeleteQuestion("1")

	got := len(store.questionList.Questions)
	want := 0
	assertNumbersEqual(t, got, want)
}
