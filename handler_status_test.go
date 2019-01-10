package transactioncategories

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	server := NewServer(nil, nil)
	req := newGetRequest(t, "/status")
	res := httptest.NewRecorder()

	server.ServeHTTP(res, req)
	result := res.Result()
	body := readBodyJSON(t, result.Body)

	assertStatusCode(t, result.StatusCode, http.StatusOK)

	var got jsonStatus
	unmarshallInterfaceFromBody(t, body, &got)
	want := jsonStatus{"OK"}
	assertDeepEqual(t, got, want)

}
