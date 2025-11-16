package web

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
	if got := w.Body.String(); !strings.Contains(got, "<!DOCTYPE html>") {
		t.Fatalf("expected HTML body, got %q", got)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Fatalf("expected text/html content type, got %s", ct)
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

func TestVisualiseHandlerReverseMode(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	payload := visualiseRequest{Input: "U+0041", Mode: "codepoints"}
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
	if len(resp.Items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(resp.Items))
	}
	if resp.Items[0].Character != "'A'" {
		t.Fatalf("expected 'A', got %s", resp.Items[0].Character)
	}
}

func TestVisualiseHandlerInvalidMode(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	payload := visualiseRequest{Input: "test", Mode: "unknown"}
	buf, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/api/visualise", bytes.NewReader(buf))
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
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

func TestDownloadHandlerJSON(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/download?format=json&input=Go", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "application/json") {
		t.Fatalf("expected json content type, got %s", ct)
	}
	if disp := w.Header().Get("Content-Disposition"); !strings.Contains(disp, "visualiser.json") {
		t.Fatalf("expected json attachment filename, got %s", disp)
	}
}

func TestDownloadHandlerCSV(t *testing.T) {
	srv := NewServer()
	mux := http.NewServeMux()
	srv.Routes(mux)

	req := httptest.NewRequest(http.MethodGet, "/api/download?format=csv&input=Go", nil)
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/csv") {
		t.Fatalf("expected csv content type, got %s", ct)
	}
	if disp := w.Header().Get("Content-Disposition"); !strings.Contains(disp, "visualiser.csv") {
		t.Fatalf("expected csv attachment filename, got %s", disp)
	}
	if !strings.Contains(w.Body.String(), "Character") {
		t.Fatalf("expected CSV header, got %s", w.Body.String())
	}
}
