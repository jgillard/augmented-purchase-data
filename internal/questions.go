package internal

import "regexp"

// QuestionStore is an interface that when implemented,
// provides methods for manipulating a store of questions,
// including some helper functions for querying the store
type QuestionStore interface {
	ListQuestionsForCategory(categoryID string) QuestionList
	GetQuestion(questionID string) Question
	AddQuestion(categoryID string, question QuestionPostRequest) Question
	RenameQuestion(questionID, questionTitle string) Question
	DeleteQuestion(questionID string)
	QuestionIDExists(questionID string) bool
	QuestionTitleExists(categoryID, questionTitle string) bool
	QuestionBelongsToCategory(questionID, categoryID string) bool
}

// QuestionList stores multiple Categorys
type QuestionList struct {
	Questions []Question `json:"questions"`
}

// Question stores all possible question attributes
// The structure implements the adjacency list pattern
// and also has a Type field (currently only "number" or "string"),
// and Options for string Questions
type Question struct {
	ID         string     `json:"id"`
	Title      string     `json:"title"`
	CategoryID string     `json:"categoryID"`
	Type       string     `json:"type"`
	Options    OptionList `json:"options"`
}

// QuestionPostRequest is a Question with no ID or CategoryID,
// as CategoryID is obtained from the reqwuest path
// Used for sending new Questions to the server
// Options is a string to allow an empty list to be sent
type QuestionPostRequest struct {
	Title   string    `json:"title"`
	Type    string    `json:"type"`
	Options *[]string `json:"options"`
}

// OptionList stores multiple Options
type OptionList []Option

// Option stores all expected option attributes
type Option struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

const questionTitleRegex = `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+\??$`

var PossibleOptionTypes = []string{"string", "number"}

func IsValidQuestionTitle(title string) bool {
	isValid := true

	if len(title) == 0 || len(title) > 32 {
		isValid = false
	}

	isLetterOrWhitespace := regexp.MustCompile(questionTitleRegex).MatchString
	if !isLetterOrWhitespace(title) {
		isValid = false
	}

	return isValid
}
