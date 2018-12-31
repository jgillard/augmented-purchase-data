package transactioncategories

type QuestionStore interface {
	ListQuestionsForCategory(categoryID string) QuestionList
}

type QuestionList struct {
	Questions []Question `json:"questions"`
}

type Question struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}
