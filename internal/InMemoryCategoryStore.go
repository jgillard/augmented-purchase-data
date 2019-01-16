package internal

import "github.com/rs/xid"

// InMemoryCategoryStore is a list of categories
// with methods for querying and manipluating those categories
type InMemoryCategoryStore struct {
	categories CategoryList
}

// NewInMemoryCategoryStore returns an initialised InMemoryCategoryStore pointer
func NewInMemoryCategoryStore(c CategoryList) *InMemoryCategoryStore {
	return &InMemoryCategoryStore{c}
}

func (s *InMemoryCategoryStore) ListCategories() CategoryList {
	return s.categories
}

func (s *InMemoryCategoryStore) GetCategory(id string) CategoryGetResponse {
	category := Category{}
	for _, c := range s.categories.Categories {
		if c.ID == id {
			category = c
		}
	}

	if category == (Category{}) {
		return CategoryGetResponse{}
	}

	children := []Category{}
	for _, c := range s.categories.Categories {
		if c.ParentID == category.ID {
			children = append(children, c)
		}
	}

	response := CategoryGetResponse{
		category,
		children,
	}
	return response
}

func (s *InMemoryCategoryStore) AddCategory(categoryName, parentID string) Category {
	newCat := Category{
		ID:       xid.New().String(),
		Name:     categoryName,
		ParentID: parentID,
	}

	s.categories.Categories = append(s.categories.Categories, newCat)

	return newCat
}

func (s *InMemoryCategoryStore) RenameCategory(id, name string) Category {
	index := 0

	for i, c := range s.categories.Categories {
		if c.ID == id {
			index = i
			s.categories.Categories[index].Name = name
			break
		}
	}

	return s.categories.Categories[index]
}

func (s *InMemoryCategoryStore) DeleteCategory(id string) {
	index := 0

	for i, c := range s.categories.Categories {
		if c.ID == id {
			index = i
			break
		}
	}

	s.categories.Categories = append(s.categories.Categories[:index], s.categories.Categories[index+1:]...)
}

func (s *InMemoryCategoryStore) CategoryIDExists(categoryID string) bool {
	alreadyExists := false

	for _, c := range s.categories.Categories {
		if c.ID == categoryID {
			alreadyExists = true
		}
	}

	return alreadyExists
}

func (s *InMemoryCategoryStore) CategoryNameExists(categoryName string) bool {
	alreadyExists := false

	for _, c := range s.categories.Categories {
		if c.Name == categoryName {
			alreadyExists = true
		}
	}

	return alreadyExists
}

func (s *InMemoryCategoryStore) CategoryParentIDExists(parentID string) bool {
	exists := false

	for _, c := range s.categories.Categories {
		if c.ID == parentID {
			exists = true
		}
	}

	return exists
}

func (s *InMemoryCategoryStore) GetCategoryDepth(categoryID string) int {
	depth := 0

	for _, c := range s.categories.Categories {
		if c.ID == categoryID {
			// if already a subcategory
			if c.ParentID != "" {
				depth = 1
			}
		}
	}

	return depth
}
