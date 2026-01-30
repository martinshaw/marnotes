package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	documentsDirectory string
	webappDirectory    string
	port               string
	jsonOnly           bool
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

// serveWebApp serves the React web application
func (s *Server) serveWebApp(w http.ResponseWriter, r *http.Request) {
	indexPath := filepath.Join(s.webappDirectory, "index.html")
	http.ServeFile(w, r, indexPath)
}

func main() {
	webappDirectory := "./server/web"

	documentsDirectory := flag.String("dir", "./documents", "Directory containing JSON documents")
	port := flag.String("port", ":8080", "Port to listen on (e.g., :8080)")
	jsonOnly := flag.Bool("json-only", false, "Start only the JSON server")
	flag.Parse()

	server := &Server{
		documentsDirectory: *documentsDirectory,
		webappDirectory:    webappDirectory,
		port:               *port,
		jsonOnly:           *jsonOnly,
	}

	// Verify required directories exist
	if _, err := os.Stat(*documentsDirectory); os.IsNotExist(err) {
		createErr := os.MkdirAll(*documentsDirectory, 0755)
		if os.IsNotExist(createErr) {
			log.Fatalf("Failed to create documents directory: %s", *documentsDirectory)
		}

		log.Printf("Created documents directory: %s", *documentsDirectory)
	}

	if _, err := os.Stat(webappDirectory); os.IsNotExist(err) {
		log.Fatalf("Web app directory does not exist: %s", webappDirectory)
	}

	// Register routes based on flags
	mux := http.NewServeMux()

	if !*jsonOnly {
		// Only web app routes
		mux.HandleFunc("/", server.serveWebApp)
		log.Printf("Starting web app server on %s, serving from: %s", *port, webappDirectory)
	} else if *jsonOnly {
		// Only JSON server routes
		mux.HandleFunc("/", server.healthCheck)
		mux.HandleFunc("/health", server.healthCheck)
		mux.HandleFunc("/documents", server.listDocuments)
		mux.HandleFunc("/doc/", server.serveDocument)
		log.Printf("Starting JSON server on %s, serving documents from: %s", *port, *documentsDirectory)
	} else {
		// Both servers (default)
		// JSON API routes
		mux.HandleFunc("/health", server.healthCheck)
		mux.HandleFunc("/api/health", server.healthCheck)
		mux.HandleFunc("/documents", server.listDocuments)
		mux.HandleFunc("/doc/", server.serveDocument)

		// Web app route (default/root)
		mux.HandleFunc("/", server.serveWebApp)

		log.Printf("Starting combined server on %s", *port)
		log.Printf("  - JSON API serving documents from: %s", *documentsDirectory)
		log.Printf("  - Web app serving from: %s", webappDirectory)
	}

	if err := http.ListenAndServe(*port, mux); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
