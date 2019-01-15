package transactioncategories

const (
	// Generic
	errorInvalidJSON  = "request JSON invalid"
	errorFieldMissing = "a required field is missing from the request"

	// Category
	errorCategoryNotFound      = "categoryID not found"
	errorDuplicateCategoryName = "name is a duplicate"
	errorInvalidCategoryName   = "name is invalid"
	errorParentIDNotFound      = "parentID not found"
	errorCategoryTooNested     = "category would be too nested"

	//Question
	errorQuestionNotFound               = "question not found"
	errorTitleEmpty                     = "title is empty"
	errorInvalidTitle                   = "title is invalid"
	errorDuplicateTitle                 = "title is a duplicate"
	errorTypeEmpty                      = "type is empty"
	errorInvalidType                    = "type is invalid"
	errorOptionsInvalid                 = "options is invalid"
	errorOptionEmpty                    = "option is empty"
	errorDuplicateOption                = "options list has a duplicate"
	errorQuestionDoesntBelongToCategory = "question does not belong to category"
)
