package rpctransport

import (
	"context"
	"testing"
)

func TestGetStatus(t *testing.T) {
	server := NewServer(nil, nil)

	req := &EmptyRequest{}

	res, err := server.GetStatus(context.Background(), req)
	if err != nil {
		t.Errorf("GetStatus(%v) got unexpected error", req)
	}

	got := res.Status
	want := "OK"
	assertStringsEqual(t, got, want)
}
