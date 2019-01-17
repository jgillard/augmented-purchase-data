package rpctransport

import (
	"reflect"
	"testing"
)

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
