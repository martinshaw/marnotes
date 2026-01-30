# JSON Document Server

A Go web server that serves JSON documents from a configurable directory and includes a React-based web dashboard.

## Features

- Dual-mode HTTP server: JSON API + Web Dashboard
- Configurable document directory
- RESTful endpoints for listing and retrieving documents
- React web interface displaying server status and port
- Health check endpoint
- Input validation to prevent directory traversal attacks
- JSON validation on serve
- Flexible deployment: run both servers together or independently

## Building

```bash
go build
```

## Running

### Run Both Servers (Default)

```bash
# Both JSON API and web app on the same port
./jsonserver
./jsonserver -port :8080
```

### Run JSON Server Only

```bash
./jsonserver -json-only
./jsonserver -json-only -dir /path/to/json/files -port :3000
```

### Run Web App Only

```bash
./jsonserver -webapp-only
./jsonserver -webapp-only -port :3000
```

### All Available Flags

```bash
-dir string       # Directory containing JSON documents (default "./documents")
-webapp string    # Directory containing web application (default "./webapp")
-port string      # Port to listen on (default ":8080")
-json-only        # Start only the JSON server
-webapp-only      # Start only the web app server
```

## API Endpoints

### Web Dashboard

- `GET /` - React web dashboard (when running in default or webapp-only mode)

### JSON API

- `GET /health` or `GET /api/health` - Returns server status
- `GET /documents` - Returns list of all JSON files in the directory
- `GET /doc/{filename}` - Returns a specific JSON document

**Examples:**

```bash
# Health check
curl http://localhost:8080/health

# List all documents
curl http://localhost:8080/documents

# Get specific document
curl http://localhost:8080/doc/users
curl http://localhost:8080/doc/example
```

## Example Usage

### Combined Mode (Default)

1. Build and run:

```bash
go build
./jsonserver
```

2. Open your browser to `http://localhost:8080` to see the React dashboard

3. Access the JSON API:

```bash
curl http://localhost:8080/documents
curl http://localhost:8080/doc/users
curl http://localhost:8080/doc/example
```

### JSON-Only Mode

1. Run with JSON server only:

```bash
./jsonserver -json-only -port :3000
```

2. Access documents:

```bash
curl http://localhost:3000/documents
curl http://localhost:3000/doc/users
```

## Project Structure

```
.
├── main.go          # Main application with dual-server support
├── go.mod           # Go module definition
├── documents/       # Directory containing JSON files
│   ├── example.json
│   └── users.json
└── webapp/          # React web application
    └── index.html
```
