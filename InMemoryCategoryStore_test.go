package transactioncategories

import "testing"

func TestNewInMemoryCategoryStore(t *testing.T) {
	got := NewInMemoryCategoryStore(CategoryList{})
	want := &InMemoryCategoryStore{}
	assertDeepEqual(t, got, want)
}

func TestInMemoryCategoryStore_listCategories(t *testing.T) {
	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation"},
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	got := store.listCategories()
	want := categoryList
	assertDeepEqual(t, got, want)
}

func TestInMemoryCategoryStore_getCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	t.Run("ID doesn't exist", func(t *testing.T) {
		got := store.getCategory("abcd")
		want := CategoryGetResponse{}
		assertDeepEqual(t, got, want)
	})

	t.Run("ID exists", func(t *testing.T) {
		got := store.getCategory("1234")
		want := CategoryGetResponse{category, []Category{}}
		assertDeepEqual(t, got, want)
	})
}

func TestInMemoryCategoryStore_addCategory(t *testing.T) {
	categoryList := CategoryList{}
	store := NewInMemoryCategoryStore(categoryList)

	categoryName := "accomodation"
	parentID := "1234"

	got := store.addCategory(categoryName, parentID)

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

func TestInMemoryCategoryStore_renameCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	newName := "new name"

	got := store.renameCategory("1234", newName)

	// assert response
	assertStringsEqual(t, got.Name, newName)

	// assert store
	got = store.categories.Categories[0]
	assertStringsEqual(t, got.ID, "1234")
	assertStringsEqual(t, got.Name, newName)
}

func TestInMemoryCategoryStore_deleteCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(categoryList)

	store.deleteCategory("1234")

	got := len(store.categories.Categories)
	want := 0
	assertNumbersEqual(t, got, want)
}
