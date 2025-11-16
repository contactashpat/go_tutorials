package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHomeHandler(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if got := w.Body.String(); got == "" {
		t.Fatalf("expected non-empty body")
	}
}

func TestVisualiseHandler(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	payload := visualiseRequest{Input: "Go"}
	buf, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/visualise", bytes.NewReader(buf))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	var resp visualiseResponse
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	if len(resp.Items) != len([]rune("Go")) {
		t.Fatalf("expected %d items, got %d", len([]rune("Go")), len(resp.Items))
	}
}

func TestVisualiseHandlerValidation(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	req := httptest.NewRequest(http.MethodPost, "/api/visualise", bytes.NewReader([]byte(`{"input":""}`)))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty input, got %d", w.Code)
	}
}
