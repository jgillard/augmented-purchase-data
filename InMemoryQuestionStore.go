package transactioncategories

import (
	"github.com/rs/xid"
)

type InMemoryQuestionStore struct {
	questionList QuestionList
}

func NewInMemoryQuestionStore(q QuestionList) *InMemoryQuestionStore {
	return &InMemoryQuestionStore{q}
}

func (s *InMemoryQuestionStore) listQuestionsForCategory(categoryID string) QuestionList {
	var questionList QuestionList
	for _, q := range s.questionList.Questions {
		if q.CategoryID == categoryID {
			questionList.Questions = append(questionList.Questions, q)
		}
	}
	return questionList
}

func (s *InMemoryQuestionStore) getQuestion(questionID string) Question {
	var question = Question{}

	for _, q := range s.questionList.Questions {
		if q.ID == questionID {
			question = q
		}
	}

	return question
}

func (s *InMemoryQuestionStore) addQuestion(categoryID string, q QuestionPostRequest) Question {
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

func (s *InMemoryQuestionStore) renameQuestion(questionID, questionTitle string) Question {
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

func (s *InMemoryQuestionStore) deleteQuestion(questionID string) {
	index := 0
	for i, q := range s.questionList.Questions {
		if q.ID == questionID {
			index = i
			break
		}
	}
	s.questionList.Questions = append(s.questionList.Questions[:index], s.questionList.Questions[index+1:]...)
}

func (s *InMemoryQuestionStore) questionIDExists(questionID string) bool {
	exists := false
	for _, q := range s.questionList.Questions {
		if q.ID == questionID {
			exists = true
		}
	}
	return exists
}

func (s *InMemoryQuestionStore) questionTitleExists(categoryID, questionTitle string) bool {
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

func (s *InMemoryQuestionStore) questionBelongsToCategory(questionID, categoryID string) bool {
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
