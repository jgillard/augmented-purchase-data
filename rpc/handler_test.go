package rpctransport

import (
	"context"
	"testing"

	"github.com/jgillard/practising-go-tdd/internal"
)

func TestGetStatus(t *testing.T) {
	server := NewServer(nil, nil)

	req := &EmptyRequest{}

	res, err := server.GetStatus(context.Background(), req)
	if err != nil {
		t.Errorf("GetStatus(%v) got unexpected error", req)
	}

	got := res.Status
	want := "OK"
	assertStringsEqual(t, got, want)
}

func TestListCategories(t *testing.T) {
	categoryList := internal.CategoryList{
		Categories: []internal.Category{
			internal.Category{ID: "1234", Name: "accommodation", ParentID: ""},
		},
	}
	categoryStore := internal.NewInMemoryCategoryStore(&categoryList)
	server := NewServer(categoryStore, nil)

	req := &EmptyRequest{}

	res, err := server.ListCategories(context.Background(), req)
	if err != nil {
		t.Errorf("ListCategories(%v) got unexpected error", req)
	}

	got := res.Categories
	want := []*GetCategoryReply{
		&GetCategoryReply{ID: "1234", Name: "accommodation", ParentID: ""},
	}
	assertDeepEqual(t, got, want)
}
