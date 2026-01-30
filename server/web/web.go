package web

import (
	"html/template"
	"net/http"
	"path/filepath"
)

type Server struct {
	webappDirectory string
	jsonPort        string
	template        *template.Template
	staticDir       string
}

type PageData struct {
	JSONPort string
}

func NewServer(webappDirectory, jsonPort string) (*Server, error) {
	indexPath := filepath.Join(webappDirectory, "public", "index.html")
	tpl, err := template.ParseFiles(indexPath)
	if err != nil {
		return nil, err
	}

	staticDir := filepath.Join(webappDirectory, "public")

	return &Server{
		webappDirectory: webappDirectory,
		jsonPort:        jsonPort,
		template:        tpl,
		staticDir:       staticDir,
	}, nil
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(s.staticDir))))
	mux.HandleFunc("/", s.serveIndex)
	return mux
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_ = s.template.Execute(w, PageData{JSONPort: s.jsonPort})
}
