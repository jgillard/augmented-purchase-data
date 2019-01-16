package internal

import "regexp"

// CategoryStore is an interface that when implemented,
// provides methods for manipulating a store of categories,
// including some helper functions for querying the store
type CategoryStore interface {
	ListCategories() CategoryList
	GetCategory(categoryID string) Category
	GetChildCategories(categoryID string) []Category
	AddCategory(categoryName, parentID string) Category
	RenameCategory(categoryID, categoryName string) Category
	DeleteCategory(categoryID string)

	CategoryIDExists(categoryID string) bool
	CategoryNameExists(categoryName string) bool
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

const categoryNameRegex = `^[a-zA-Z]+[a-zA-Z ]+?[a-zA-Z]+$`

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
