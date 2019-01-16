package internal

import "regexp"

const categoryNameRegex = `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+$`

// CategoryStore is an interface that when implemented,
// provides methods for manipulating a store of categories,
// including some helper functions for querying the store
type CategoryStore interface {
	ListCategories() CategoryList
	GetCategory(categoryID string) CategoryGetResponse
	AddCategory(categoryName, parentID string) Category
	RenameCategory(categoryID, categoryName string) Category
	DeleteCategory(categoryID string)
	CategoryIDExists(categoryID string) bool
	CategoryNameExists(categoryName string) bool
	CategoryParentIDExists(categoryParentID string) bool
	GetCategoryDepth(categoryID string) int
}

// CategoryList stores multiple categories
type CategoryList struct {
	Categories []Category `json:"categories"`
}

// Category stores all expected category attributes
// The structure implements the adjacency list pattern
type Category struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parentID"`
}

// CategoryGetResponse returns a Category in addition to its immediate child categories
type CategoryGetResponse struct {
	Category
	Children []Category `json:"children"`
}

// CategoryPostRequest is a Category with no ID
// Used for sending new Categorys to the server
// ParentID is a string to allow "" to signify a top-level Category
type CategoryPostRequest struct {
	Name     string  `json:"name"`
	ParentID *string `json:"parentID"`
}

func IsValidCategoryName(name string) bool {
	isValid := true

	if len(name) == 0 || len(name) > 32 {
		isValid = false
	}

	isLetterOrWhitespace := regexp.MustCompile(categoryNameRegex).MatchString
	if !isLetterOrWhitespace(name) {
		isValid = false
	}

	return isValid
}
