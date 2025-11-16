package web

import (
	"bytes"
	"embed"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"go_tutorials/internal/reverseinput"
	"go_tutorials/internal/visualiser"
)

//go:embed templates/index.html
var templateFS embed.FS

var indexTemplate = template.Must(template.ParseFS(templateFS, "templates/index.html"))

// Server exposes HTTP handlers for the visualiser project.
type Server struct {
	indexTmpl *template.Template
}

// NewServer builds a Server.
func NewServer() *Server {
	return &Server{
		indexTmpl: indexTemplate,
	}
}

// Routes wires HTTP handlers onto the provided ServeMux.
func (s *Server) Routes(mux *http.ServeMux) {
	mux.HandleFunc("/", s.handleHome)
	mux.HandleFunc("/api/visualise", s.handleVisualise)
	mux.HandleFunc("/api/download", s.handleDownload)
}

type visualiseRequest struct {
	Input string `json:"input"`
	Mode  string `json:"mode"`
}

type visualiseResponse struct {
	Items []visualiser.Result `json:"items"`
}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := s.indexTmpl.Execute(w, nil); err != nil {
		http.Error(w, "failed to render template", http.StatusInternalServerError)
		return
	}
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
	resolved, err := s.resolveInput(req.Mode, req.Input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := visualiser.AnalyseString(resolved)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := visualiseResponse{Items: results}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	query := r.URL.Query()
	format := strings.ToLower(query.Get("format"))
	if format == "" {
		format = "json"
	}
	input := query.Get("input")
	mode := query.Get("mode")

	resolved, err := s.resolveInput(mode, input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	results, err := visualiser.AnalyseString(resolved)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch format {
	case "json":
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Disposition", "attachment; filename=\"visualiser.json\"")
		json.NewEncoder(w).Encode(visualiseResponse{Items: results})
	case "csv":
		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)
		writer.Write([]string{
			"Character",
			"CodePointHex",
			"CodePointDec",
			"UTF8BytesHex",
			"UTF8BytesDec",
			"UTF8BytesBinary",
			"HTMLEntityDecimal",
			"HTMLEntityHex",
		})
		for _, item := range results {
			writer.Write([]string{
				item.Character,
				item.CodePointHex,
				fmt.Sprintf("%d", item.CodePointDec),
				strings.Join(item.UTF8BytesHex, " "),
				strings.Join(item.UTF8BytesDec, " "),
				strings.Join(item.UTF8BytesBinary, " "),
				item.HTMLEntityDecimal,
				item.HTMLEntityHex,
			})
		}
		writer.Flush()
		if err := writer.Error(); err != nil {
			http.Error(w, "failed to generate CSV", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", "attachment; filename=\"visualiser.csv\"")
		w.Write(buf.Bytes())
	default:
		http.Error(w, "unsupported format", http.StatusBadRequest)
	}
}

func (s *Server) resolveInput(mode, input string) (string, error) {
	if strings.TrimSpace(input) == "" {
		return "", errors.New("input is required")
	}
	normalized := strings.ToLower(strings.TrimSpace(mode))
	switch normalized {
	case "", "text":
		return input, nil
	case "codepoint", "codepoints", "cp":
		tokens := reverseinput.Tokenize(input)
		if len(tokens) == 0 {
			return "", errors.New("code point values required")
		}
		return reverseinput.BuildStringFromCodePoints(tokens)
	case "byte", "bytes":
		tokens := reverseinput.Tokenize(input)
		if len(tokens) == 0 {
			return "", errors.New("byte values required")
		}
		return reverseinput.BuildStringFromBytes(tokens)
	default:
		return "", fmt.Errorf("unknown mode %q (use text, codepoints, bytes)", mode)
	}
}
