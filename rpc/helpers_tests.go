package rpctransport

import (
	"testing"
)

func assertStringsEqual(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got '%s' wanted '%s'", got, want)
	}
}
