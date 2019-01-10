package transactioncategories

const (
	// Generic
	ErrorInvalidJSON  = "request JSON invalid"
	ErrorFieldMissing = "a required field is missing from the request"

	// Category
	ErrorCategoryNotFound      = "categoryID not found"
	ErrorDuplicateCategoryName = "name is a duplicate"
	ErrorInvalidCategoryName   = "name is invalid"
	ErrorParentIDNotFound      = "parentID not found"
	ErrorCategoryTooNested     = "category would be too nested"

	//Question
	ErrorQuestionNotFound               = "question not found"
	ErrorTitleEmpty                     = "title is empty"
	ErrorInvalidTitle                   = "title is invalid"
	ErrorDuplicateTitle                 = "title is a duplicate"
	ErrorTypeEmpty                      = "type is empty"
	ErrorInvalidType                    = "type is invalid"
	ErrorOptionsInvalid                 = "options is invalid"
	ErrorOptionEmpty                    = "option is empty"
	ErrorDuplicateOption                = "options list has a duplicate"
	ErrorQuestionDoesntBelongToCategory = "question does not belong to category"
)
