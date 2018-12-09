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

	t.Run("check status response code", func(t *testing.T) {
		status := result.StatusCode
		desiredStatus := http.StatusOK
		assertNumbersEqual(t, status, desiredStatus)
	})

	t.Run("check status response body", func(t *testing.T) {
		body, err := ioutil.ReadAll(result.Body)
		if err != nil {
			t.Fatal(err)
		}
		bodyString := string(body)
		desiredBody := statusBodyString
		assertStringsEqual(t, bodyString, desiredBody)
	})

}
