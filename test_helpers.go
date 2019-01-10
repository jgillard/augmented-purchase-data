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

func assertIsXid(t *testing.T, s string) {
	t.Helper()
	_, err := xid.FromString(s)
	if err != nil {
		t.Fatalf("got ID '%s' which isn't an xid", s)
	}
}

func assertBodyJSONIsStatus(t *testing.T, got []byte, want string) {
	t.Helper()
	var body jsonStatus
	unmarshallInterfaceFromBody(t, got, &body)
	if body.Status != want {
		t.Errorf("wanted a json status '%s', got '%s'", want, got)
	}
}

func readBodyJSON(t *testing.T, b io.ReadCloser) []byte {
	t.Helper()
	body, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatal(err)
	}

	// this probably shouldn't live here...
	// checks that the body is JSON
	var js json.RawMessage
	if json.Unmarshal(body, &js) != nil {
		t.Fatalf("body is not json")
	}

	return js
}

func unmarshallInterfaceFromBody(t *testing.T, bodyBytes []byte, got interface{}) interface{} {
	t.Helper()
	err := json.Unmarshal(bodyBytes, got)
	// check for syntax error or type mismatch
	if err != nil {
		t.Fatal(err)
	}

	return got

}

func assertBodyErrorTitle(t *testing.T, bodyBytes []byte, title string) {
	t.Helper()
	var errors jsonErrors

	err := json.Unmarshal(bodyBytes, &errors)
	// check for syntax error or type mismatch
	if err != nil {
		t.Log("cannot unmarshall into jsonErrors")
		t.Fatal(err)
	}

	if len(errors.Errors) != 1 {
		t.Fatalf("expected %d errors in response, got %d", 1, len(errors.Errors))
	}

	assertStringsEqual(t, errors.Errors[0].Title, title)

}

func assertOptionsNil(t *testing.T, got OptionList) {
	t.Helper()
	if got != nil {
		t.Fatalf("Options should not have been set, got %v", got)
	}
}
