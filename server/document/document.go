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
	mux.HandleFunc("/", s.healthCheck)
	mux.HandleFunc("/health", s.healthCheck)
	mux.HandleFunc("/documents", s.listDocuments)
	mux.HandleFunc("/doc/", s.serveDocument)
	return mux
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
	json.NewEncoder(w).Encode(map[string]interface{}{
		"documents": documents,
		"count":     len(documents),
	})
}

// serveDocument serves a specific JSON document by filename
func (s *Server) serveDocument(w http.ResponseWriter, r *http.Request) {
	filename := strings.TrimPrefix(r.URL.Path, "/doc/")
	filename = strings.TrimSuffix(filename, ".json")

	if filename == "" {
		http.Error(w, "No filename specified", http.StatusBadRequest)
		return
	}

	// Prevent directory traversal
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") {
		http.Error(w, "Invalid filename", http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(s.documentsDirectory, filename+".json")

	data, err := ioutil.ReadFile(filePath)
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
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"docsDir": s.documentsDirectory,
	})
}
