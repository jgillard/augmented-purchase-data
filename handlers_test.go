package handlers

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStatusHandler(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()

	handler := http.HandlerFunc(StatusHandler)

	handler.ServeHTTP(res, req)

	result := res.Result()

	status := result.StatusCode
	desiredStatus := http.StatusOK
	if status != desiredStatus {
		t.Errorf("got %d wanted %d", status, desiredStatus)
	}

	body, err := ioutil.ReadAll(result.Body)
	bodyString := string(body)
	if err != nil {
		t.Fatal(err)
	}
	desiredBody := statusString
	if bodyString != desiredBody {
		t.Errorf("got '%s' wanted '%s'", bodyString, desiredBody)
	}
}
