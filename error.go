package transactioncategories

const (
	// Category
	ErrorCategoryNotFound      = "categoryID not found"
	ErrorInvalidJSON           = "request JSON invalid"
	ErrorFieldMissing          = "a required field is missing from the request"
	ErrorDuplicateCategoryName = "name is a duplicate"
	ErrorInvalidCategoryName   = "name is invalid"
	ErrorParentIDNotFound      = "parentID not found"
	ErrorCategoryTooNested     = "category would be too nested"

	//Question
	ErrorTitleEmpty                     = "title is empty"
	ErrorDuplicateTitle                 = "title is a duplicate"
	ErrorTypeEmpty                      = "type is empty"
	ErrorInvalidType                    = "type is invalid"
	ErrorOptionEmpty                    = "option is empty"
	ErrorOptionsInvalid                 = "options is invalid"
	ErrorDuplicateOption                = "options list has a duplicate"
	ErrorInvalidTitle                   = "title is invalid"
	ErrorQuestionNotFound               = "question not found"
	ErrorQuestionDoesntBelongToCategory = "question does not belong to category"
)
