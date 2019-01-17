package internal

import "testing"

func TestNewInMemoryCategoryStore(t *testing.T) {
	got := NewInMemoryCategoryStore(nil)
	want := &InMemoryCategoryStore{}
	assertDeepEqual(t, got, want)
}

func TestInMemoryCategoryStore_ListCategories(t *testing.T) {
	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation"},
		},
	}
	store := NewInMemoryCategoryStore(&categoryList)

	got := store.ListCategories()
	want := categoryList
	assertDeepEqual(t, got, want)
}

func TestInMemoryCategoryStore_GetChildCategories(t *testing.T) {
	categoryList := CategoryList{
		Categories: []Category{
			Category{ID: "1234", Name: "accommodation", ParentID: ""},
			Category{ID: "1235", Name: "foo", ParentID: "1234"},
			Category{ID: "1236", Name: "bar", ParentID: "1235"},
		},
	}
	store := NewInMemoryCategoryStore(&categoryList)

	t.Run("has children", func(t *testing.T) {
		got := store.GetChildCategories("1234")
		want := []Category{
			categoryList.Categories[1],
		}
		assertDeepEqual(t, got, want)
	})

	t.Run("no children", func(t *testing.T) {
		got := store.GetChildCategories("1236")
		want := []Category{}
		assertDeepEqual(t, got, want)
	})
}

func TestInMemoryCategoryStore_GetCategory(t *testing.T) {
	category := Category{ID: "1234", Name: "accommodation"}
	categoryList := CategoryList{
		Categories: []Category{
			category,
		},
	}
	store := NewInMemoryCategoryStore(&categoryList)

	t.Run("ID doesn't exist", func(t *testing.T) {
		got := store.GetCategory("abcd")
		want := Category{}
		assertDeepEqual(t, got, want)
	})

	t.Run("ID exists", func(t *testing.T) {
		got := store.GetCategory("1234")
		want := category
		assertDeepEqual(t, got, want)
	})
}

func TestInMemoryCategoryStore_AddCategory(t *testing.T) {
	store := NewInMemoryCategoryStore(nil)

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
	store := NewInMemoryCategoryStore(&categoryList)

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
	store := NewInMemoryCategoryStore(&categoryList)

	store.DeleteCategory("1234")

	got := len(store.categories.Categories)
	want := 0
	assertNumbersEqual(t, got, want)
}
