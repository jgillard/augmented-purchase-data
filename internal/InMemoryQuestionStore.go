package internal

import (
	"github.com/rs/xid"
)

// InMemoryQuestionStore is a list of questions
// with methods for querying and manipluating those questions
type InMemoryQuestionStore struct {
	questionList QuestionList
}

// NewInMemoryQuestionStore returns an initialised InMemoryQuestionStore pointer
func NewInMemoryQuestionStore(q *QuestionList) *InMemoryQuestionStore {
	if q == nil {
		return &InMemoryQuestionStore{}
	}
	return &InMemoryQuestionStore{*q}
}

func (s *InMemoryQuestionStore) ListQuestions() QuestionList {
	return s.questionList
}

func (s *InMemoryQuestionStore) ListQuestionsForCategory(categoryID string) QuestionList {
	var questionList QuestionList
	for _, q := range s.questionList.Questions {
		if q.CategoryID == categoryID {
			questionList.Questions = append(questionList.Questions, q)
		}
	}
	return questionList
}

func (s *InMemoryQuestionStore) GetQuestion(questionID string) Question {
	var question = Question{}

	for _, q := range s.questionList.Questions {
		if q.ID == questionID {
			question = q
		}
	}

	return question
}

func (s *InMemoryQuestionStore) AddQuestion(categoryID string, q QuestionPostRequest) Question {
	question := Question{
		ID:         xid.New().String(),
		Title:      q.Title,
		CategoryID: categoryID,
		Type:       q.Type,
	}

	if q.Type == "string" {
		question.Options = OptionList{}
		for _, title := range *q.Options {
			option := Option{
				ID:    xid.New().String(),
				Title: title,
			}
			question.Options = append(question.Options, option)
		}
	}

	s.questionList.Questions = append(s.questionList.Questions, question)

	return question
}

func (s *InMemoryQuestionStore) RenameQuestion(questionID, questionTitle string) Question {
	index := 0

	for i, q := range s.questionList.Questions {
		if q.ID == questionID {
			index = i
			s.questionList.Questions[index].Title = questionTitle
			break
		}
	}

	return s.questionList.Questions[index]
}

func (s *InMemoryQuestionStore) DeleteQuestion(questionID string) {
	index := 0
	for i, q := range s.questionList.Questions {
		if q.ID == questionID {
			index = i
			break
		}
	}
	s.questionList.Questions = append(s.questionList.Questions[:index], s.questionList.Questions[index+1:]...)
}

func (s *InMemoryQuestionStore) QuestionIDExists(questionID string) bool {
	exists := false
	for _, q := range s.questionList.Questions {
		if q.ID == questionID {
			exists = true
		}
	}
	return exists
}

func (s *InMemoryQuestionStore) QuestionTitleExists(categoryID, questionTitle string) bool {
	alreadyExists := false
	for _, q := range s.questionList.Questions {
		if q.CategoryID == categoryID {
			if q.Title == questionTitle {
				alreadyExists = true
			}
		}
	}
	return alreadyExists
}

func (s *InMemoryQuestionStore) QuestionBelongsToCategory(questionID, categoryID string) bool {
	belongsToCategory := true
	for _, q := range s.questionList.Questions {
		if q.ID == questionID {
			if q.CategoryID != categoryID {
				belongsToCategory = false
			}
		}
	}
	return belongsToCategory
}
