package document

import (
	"crypto/rsa"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"martinshaw.co/marnotes/server/crypto"
)

type Server struct {
	documentsDirectory string
	privateKey         *rsa.PrivateKey
	publicKey          *rsa.PublicKey
}

func NewServer(documentsDirectory string) *Server {
	return &Server{documentsDirectory: documentsDirectory}
}

func NewServerWithEncryption(documentsDirectory, keyDir string) (*Server, error) {
	kp, err := crypto.LoadOrGenerateKeyPair(keyDir, 2048)
	if err != nil {
		return nil, err
	}

	return &Server{
		documentsDirectory: documentsDirectory,
		privateKey:         kp.PrivateKey,
		publicKey:          kp.PublicKey,
	}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.corsMiddleware(s.healthCheck))
	mux.HandleFunc("/health", s.corsMiddleware(s.healthCheck))
	mux.HandleFunc("/documents", s.corsMiddleware(s.listDocuments))
	mux.HandleFunc("/doc/", s.corsMiddleware(s.serveDocument))
	mux.HandleFunc("/publickey", s.servePublicKey)
	return mux
}

func (s *Server) servePublicKey(w http.ResponseWriter, r *http.Request) {
	if s.publicKey == nil {
		http.Error(w, "Public key not available", http.StatusBadRequest)
		return
	}

	publicKeyBytes, err := crypto.GetPublicKeyPEM(s.publicKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to export public key: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Write(publicKeyBytes)
}

func (s *Server) encryptResponse(data []byte) ([]byte, error) {
	if s.publicKey == nil {
		return data, nil
	}

	encrypted, err := crypto.Encrypt(s.publicKey, data)
	if err != nil {
		return nil, err
	}

	response := map[string]string{
		"encrypted": encrypted,
	}

	return json.Marshal(response)
}

func (s *Server) decryptRequest(encryptedPayload string) ([]byte, error) {
	if s.privateKey == nil {
		return []byte(encryptedPayload), nil
	}

	return crypto.Decrypt(s.privateKey, encryptedPayload)
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

	// TODO: Implement proper client-side RSA decryption before enabling encryption
	// For now, send unencrypted responses
	w.Write(responseData)
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

	// TODO: Implement proper client-side RSA decryption before enabling encryption
	// For now, send unencrypted responses
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

	// TODO: Implement proper client-side RSA decryption before enabling encryption
	// For now, send unencrypted responses
	w.Write(responseData)
}
