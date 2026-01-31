package document

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	documentsDirectory string
}

func NewServer(documentsDirectory string) *Server {
	return &Server{documentsDirectory: documentsDirectory}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.corsMiddleware(s.healthCheck))
	mux.HandleFunc("/documents", s.corsMiddleware(s.listDocuments))
	mux.HandleFunc("/documents/{filePath...}", s.corsMiddleware(s.serveDocument))
	return mux
}

func (s *Server) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

// listDocuments lists all JSON files in the configured directory
func (s *Server) listDocuments(w http.ResponseWriter, r *http.Request) {
	files, err := ioutil.ReadDir(s.documentsDirectory)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read directory: %v", err), http.StatusInternalServerError)
		return
	}

	var documents []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			documents = append(documents, file.Name())
		}
	}

	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"documents": documents,
		"count":     len(documents),
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(responseData)
}

// serveDocument serves a specific JSON document by filePath
func (s *Server) serveDocument(w http.ResponseWriter, r *http.Request) {
	filePath := r.PathValue("filePath")

	if filePath == "" {
		http.Error(w, "No file path specified", http.StatusBadRequest)
		return
	}

	// Prevent directory traversal
	if strings.Contains(filePath, "..") || strings.Contains(filePath, "/") {
		http.Error(w, "Invalid file path", http.StatusBadRequest)
		return
	}

	filePath = filepath.Join(s.documentsDirectory, filePath)

	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			http.Error(w, "Document not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to read file: %v", err), http.StatusInternalServerError)
		}
		return
	}

	// Validate JSON
	var jsonData interface{}
	if err := json.Unmarshal(data, &jsonData); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.Write(data)
}

// healthCheck provides a simple health check endpoint
func (s *Server) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{
		"status":  "healthy",
		"docsDir": s.documentsDirectory,
	}

	responseData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to marshal response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(responseData)
}
