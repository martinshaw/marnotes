# Marnotes

A Go-based JSON document server with a React/Lexical web dashboard.

## Features

- RESTful API for listing and retrieving JSON documents
- React/TypeScript frontend with Blueprint UI
- Lexical-based rich text editor supporting headings, lists, and formatting
- Documents menu bar item: dropdown lists all available documents from the API
- Click a document to load it into the editor (supports both Lexical and plain JSON)
- Editor auto-focuses for immediate typing
- Navigation elements use `tabIndex={-1}` to keep keyboard focus in the editor
- CORS enabled for API endpoints

## API Endpoints

- `GET /health` — Server status and config
- `GET /documents` — List all JSON files in the documents directory
- `GET /documents/{filename}` — Retrieve a specific document (supports both plain JSON and Lexical state)

## Web Interface

- **Documents menu**: Dropdown with all available documents
- **Editor**: Loads and displays the selected document (as rich text if Lexical, or as formatted JSON)
- **Auto-focus**: Editor keeps focus for fast editing
- **Tab-index**: Navigation and menu items are not focusable by keyboard, so the editor always keeps focus

## Project Structure

```
.
├── main.go
├── go.mod
├── documents/
│ ├── example.json
│ ├── users.json
│ └── lexical-example.json
└── server/
├── server.go
└── document/
└── document.go
└── web/
├── web.go
├── package.json
├── tsconfig.json
└── public/
├── index.html.tpl
└── typescript/
├── app.tsx
└── components/
├── Editor.tsx
└── NavBar.tsx
```

## Usage

1. Build and run the Go server (`go build && ./marnotes`)
2. Open `http://localhost:8080` in your browser
3. Use the Documents menu to load and edit documents

---

All Rights Reserved, (c) 2026 marnotes
