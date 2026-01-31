package server

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"martinshaw.co/marnotes/server/document"
	"martinshaw.co/marnotes/server/web"
)

func Run() {
	documentsDirectory := flag.String("dir", "./documents", "Directory containing JSON documents")
	documentPort := flag.String("document-port", ":8080", "Port for the JSON document server (e.g., :8080)")
	webPort := flag.String("web-port", ":3000", "Port for the web app server (e.g., :3000)")
	webappDirectory := flag.String("web-dir", "./server/web", "Directory containing the web application")
	jsonOnly := flag.Bool("json-only", false, "Start only the JSON server")
	webOnly := flag.Bool("web-only", false, "Start only the web app server")
	flag.Parse()

	if *jsonOnly && *webOnly {
		log.Fatal("Cannot specify both -json-only and -web-only")
	}

	if !*webOnly {
		if _, err := os.Stat(*documentsDirectory); os.IsNotExist(err) {
			if createErr := os.MkdirAll(*documentsDirectory, 0755); createErr != nil {
				log.Fatalf("Failed to create documents directory: %s", *documentsDirectory)
			}
			log.Printf("Created documents directory: %s", *documentsDirectory)
		}
	}

	if !*jsonOnly {
		if _, err := os.Stat(*webappDirectory); os.IsNotExist(err) {
			log.Fatalf("Web app directory does not exist: %s", *webappDirectory)
		}
	}

	documentPort = ensureAvailablePort(*documentPort)
	jsonPortValue := strings.TrimPrefix(*documentPort, ":")
	if jsonPortValue == "" {
		jsonPortValue = "8080"
	}

	if *jsonOnly {
		docServer := document.NewServer(*documentsDirectory)
		log.Printf("Starting JSON server on %s, serving documents from: %s", *documentPort, *documentsDirectory)
		if err := runHTTPServer(*documentPort, docServer.Handler()); err != nil {
			log.Fatalf("JSON server failed to start: %v", err)
		}
		return
	}

	if *webOnly {
		if err := buildWebAssets(*webappDirectory); err != nil {
			log.Fatalf("Failed to build web app assets: %v", err)
		}
		webPort = ensureAvailablePort(*webPort)
		webServer, err := web.NewServer(*webappDirectory, jsonPortValue)
		if err != nil {
			log.Fatalf("Failed to load web app template: %v", err)
		}
		log.Printf("Starting web app server on %s, serving from: %s", *webPort, *webappDirectory)
		if err := runHTTPServer(*webPort, webServer.Handler()); err != nil {
			log.Fatalf("Web app server failed to start: %v", err)
		}
		return
	}

	if err := buildWebAssets(*webappDirectory); err != nil {
		log.Fatalf("Failed to build web app assets: %v", err)
	}

	webPort = ensureAvailablePort(*webPort)
	docServer := document.NewServer(*documentsDirectory)
	webServer, err := web.NewServer(*webappDirectory, jsonPortValue)
	if err != nil {
		log.Fatalf("Failed to load web app template: %v", err)
	}

	log.Printf("Starting JSON server on %s, serving documents from: %s", *documentPort, *documentsDirectory)
	log.Printf("Starting web app server on %s, serving from: %s", *webPort, *webappDirectory)

	var wg sync.WaitGroup
	errCh := make(chan error, 2)

	wg.Add(2)
	go func() {
		defer wg.Done()
		errCh <- runHTTPServer(*documentPort, docServer.Handler())
	}()

	go func() {
		defer wg.Done()
		errCh <- runHTTPServer(*webPort, webServer.Handler())
	}()

	wg.Wait()
	close(errCh)

	for err := range errCh {
		if err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}
}

func runHTTPServer(addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	return server.ListenAndServe()
}

func buildWebAssets(webappDirectory string) error {
	packageJSON := filepath.Join(webappDirectory, "package.json")
	if _, err := os.Stat(packageJSON); err != nil {
		return err
	}

	nodeModules := filepath.Join(webappDirectory, "node_modules")
	if _, err := os.Stat(nodeModules); os.IsNotExist(err) {
		if err := runWebCommand(webappDirectory, "npm", "install"); err != nil {
			return err
		}
		log.Println("✓ Web dependencies installed")
	}

	// if err := runWebCommand(webappDirectory, "npm", "run", "build"); err != nil {
	if err := runWebCommand(webappDirectory, "npm", "run", "build:dev"); err != nil {
		return err
	}
	log.Println("✓ Web assets compiled successfully")
	return nil
}

func runWebCommand(dir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Run()
}

func isPortAvailable(port string) bool {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

func findAvailablePort(basePort int) string {
	for port := basePort; port < basePort+100; port++ {
		addr := fmt.Sprintf(":%d", port)
		if isPortAvailable(addr) {
			return addr
		}
	}
	return fmt.Sprintf(":%d", basePort)
}

func ensureAvailablePort(port string) *string {
	if isPortAvailable(port) {
		return &port
	}

	basePort := 8080
	if port != ":8080" {
		fmt.Sscanf(port, ":%d", &basePort)
	}

	availablePort := findAvailablePort(basePort)
	log.Printf("Port %s is in use, using %s instead", port, availablePort)
	return &availablePort
}
