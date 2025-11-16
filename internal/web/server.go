package web

import (
	"encoding/json"
	"net/http"

	"go_tutorials/internal/visualiser"
)

// Server exposes HTTP handlers for the visualiser project.
type Server struct{}

// NewServer builds a Server.
func NewServer() *Server {
	return &Server{}
}

// Routes wires HTTP handlers onto the provided ServeMux.
func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/api/visualise", s.handleVisualise)
}

type visualiseRequest struct {
	Input string `json:"input"`
}

type visualiseResponse struct {
	Items []visualiser.Result `json:"items"`
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("String Visualiser API is running. POST /api/visualise for JSON visualisations."))
}

func (s *Server) handleVisualise(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req visualiseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}
	if req.Input == "" {
		http.Error(w, "input is required", http.StatusBadRequest)
		return
	}

	results, err := visualiser.AnalyseString(req.Input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := visualiseResponse{Items: results}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
