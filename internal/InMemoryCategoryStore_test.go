package internal

import "testing"

func TestNewInMemoryCategoryStore(t *testing.T) {
	got := NewInMemoryCategoryStore(CategoryList{})
	want := &InMemoryCategoryStore{}
	assertDeepEqual(t, got, want)
}

func TestInMemoryCategoryStore_ListCategories(t *testing.T) {
	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation"},
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	got := store.ListCategories()
	want := categoryList
	assertDeepEqual(t, got, want)
}

func TestInMemoryCategoryStore_GetCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	t.Run("ID doesn't exist", func(t *testing.T) {
		got := store.GetCategory("abcd")
		want := CategoryGetResponse{}
		assertDeepEqual(t, got, want)
	})

	t.Run("ID exists", func(t *testing.T) {
		got := store.GetCategory("1234")
		want := CategoryGetResponse{category, []Category{}}
		assertDeepEqual(t, got, want)
	})
}

func TestInMemoryCategoryStore_AddCategory(t *testing.T) {
	categoryList := CategoryList{}
	store := NewInMemoryCategoryStore(categoryList)

	categoryName := "accomodation"
	parentID := "1234"

	got := store.AddCategory(categoryName, parentID)

	// assert response
	assertIsXid(t, got.ID)
	assertStringsEqual(t, got.Name, categoryName)
	assertStringsEqual(t, got.ParentID, parentID)

	// assert store
	got = store.categories.Categories[0]
	assertIsXid(t, got.ID)
	assertStringsEqual(t, got.Name, categoryName)
	assertStringsEqual(t, got.ParentID, parentID)
}

func TestInMemoryCategoryStore_RenameCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	newName := "new name"

	got := store.RenameCategory("1234", newName)

	// assert response
	assertStringsEqual(t, got.Name, newName)

	// assert store
	got = store.categories.Categories[0]
	assertStringsEqual(t, got.ID, "1234")
	assertStringsEqual(t, got.Name, newName)
}

func TestInMemoryCategoryStore_DeleteCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	store.DeleteCategory("1234")

	got := len(store.categories.Categories)
	want := 0
	assertNumbersEqual(t, got, want)
}
