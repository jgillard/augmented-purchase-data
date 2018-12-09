package handlers

import (
	"net/http"
	"testing"
)

func NewListCategoriesRequest(t *testing.T) *http.Request {
	req, err := http.NewRequest(http.MethodGet, "/categories", nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func assertNumbersEqual(t *testing.T, a, b int) {
	t.Helper()
	if a != b {
		t.Errorf("got %d wanted %d", a, b)
	}
}

func assertStringsEqual(t *testing.T, a, b string) {
	t.Helper()
	if a != b {
		t.Errorf("got '%s' wanted '%s'", a, b)
	}
}
