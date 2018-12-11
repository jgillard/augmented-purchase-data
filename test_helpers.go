package handlers

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/rs/xid"
)

func NewGetRequest(t *testing.T, path string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func NewPostRequest(t *testing.T, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func NewPutRequest(t *testing.T, path string, b []byte) *http.Request {
	body := bytes.NewBuffer(b)
	req, err := http.NewRequest(http.MethodPut, path, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func NewDeleteRequest(t *testing.T, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(http.MethodDelete, path, body)
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

func isXid(t *testing.T, str string) bool {
	t.Helper()
	_, err := xid.FromString(str)
	if err != nil {
		return false
	} else {
		return true
	}
}
