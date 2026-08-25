package httputil

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchURLBytes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world"))
	}))
	defer ts.Close()

	data, err := FetchURLBytes(context.Background(), ts.URL)
	if err != nil {
		t.Fatalf("FetchURLBytes error: %v", err)
	}
	if string(data) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(data))
	}
}
