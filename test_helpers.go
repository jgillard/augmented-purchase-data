package handlers

import (
	"net/http"
	"testing"
)

func NewGetRequest(t *testing.T, path string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, path, nil)
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
