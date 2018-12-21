package transactioncategories

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"

	"github.com/rs/xid"
)

func newGetRequest(t *testing.T, path string) *http.Request {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func newPostRequest(t *testing.T, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(http.MethodPost, path, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func newPatchRequest(t *testing.T, path string, body io.Reader) *http.Request {
	req, err := http.NewRequest(http.MethodPatch, path, body)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func newDeleteRequest(t *testing.T, path string) *http.Request {
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		t.Fatal(err)
	}
	return req
}

func assertNumbersEqual(t *testing.T, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("got %d wanted %d", got, want)
	}
}

func assertStringsEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' wanted '%s'", got, want)
	}
}

func assertDeepEqual(t *testing.T, got, want interface{}) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got '%v' wanted '%v'", got, want)
	}
}

var assertStatusCode = assertNumbersEqual
var assertContentType = assertStringsEqual
var assertBodyString = assertStringsEqual

func assertIsXid(t *testing.T, s string) {
	t.Helper()
	_, err := xid.FromString(s)
	if err != nil {
		t.Fatalf("got ID '%s' which isn't an xid", s)
	}
}

func assertBodyEmpty(t *testing.T, b io.ReadCloser) {
	t.Helper()
	got := readBodyBytes(t, b)
	if len(got) != 0 {
		t.Errorf("wanted an empty response body, got '%s'", got)
	}
}

func readBodyBytes(t *testing.T, b io.ReadCloser) []byte {
	t.Helper()
	body, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}
	return body
}

func unmarshallCategoryListFromBody(t *testing.T, b io.ReadCloser) CategoryList {
	bodyBytes := readBodyBytes(t, b)

	var got CategoryList

	err := json.Unmarshal(bodyBytes, &got)
	// check for syntax error or type mismatch
	if err != nil {
		t.Fatal(err)
	}

	return got
}

func unmarshallCategoryFromBody(t *testing.T, b io.ReadCloser) Category {
	bodyBytes := readBodyBytes(t, b)

	var got Category

	err := json.Unmarshal(bodyBytes, &got)
	// check for syntax error or type mismatch
	if err != nil {
		t.Fatal(err)
	}

	return got
}

func unmarshallCategoryGetResponseFromBody(t *testing.T, b io.ReadCloser) CategoryGetResponse {
	bodyBytes := readBodyBytes(t, b)

	var got CategoryGetResponse

	err := json.Unmarshal(bodyBytes, &got)
	// check for syntax error or type mismatch
	if err != nil {
		t.Fatal(err)
	}

	return got
}
