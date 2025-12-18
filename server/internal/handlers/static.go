package handlers

import (
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// StaticHandler serves static files from multiple directories
type StaticHandler struct {
	webPublicDir string
	webDistDir   string
}

// NewStaticHandler creates a new static file handler
// webDir is the path to the web/ directory
func NewStaticHandler(webDir string) *StaticHandler {
	return &StaticHandler{
		webPublicDir: filepath.Join(webDir, "public"),
		webDistDir:   filepath.Join(webDir, "dist"),
	}
}

// ServeHTTP handles static file requests
func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	// Serve /static/* files
	if strings.HasPrefix(path, "/static/") {
		h.serveStatic(w, r, strings.TrimPrefix(path, "/static/"))
		return
	}

	// Serve index.html for root
	if path == "/" {
		h.serveFile(w, r, filepath.Join(h.webPublicDir, "index.html"))
		return
	}

	// Try to serve from public directory
	publicPath := filepath.Join(h.webPublicDir, path)
	if h.fileExists(publicPath) {
		h.serveFile(w, r, publicPath)
		return
	}

	// 404 for unknown paths
	http.NotFound(w, r)
}

// serveStatic serves files from /static/* path
// Looks in both public/ and dist/ directories
func (h *StaticHandler) serveStatic(w http.ResponseWriter, r *http.Request, name string) {
	// Try public directory first (CSS, images, etc.)
	publicPath := filepath.Join(h.webPublicDir, name)
	if h.fileExists(publicPath) {
		h.serveFile(w, r, publicPath)
		return
	}

	// Try dist directory (built JavaScript)
	distPath := filepath.Join(h.webDistDir, name)
	if h.fileExists(distPath) {
		h.serveFile(w, r, distPath)
		return
	}

	http.NotFound(w, r)
}

// serveFile serves a single file with appropriate content type
func (h *StaticHandler) serveFile(w http.ResponseWriter, r *http.Request, path string) {
	// Security: prevent directory traversal
	cleanPath := filepath.Clean(path)
	if !strings.HasPrefix(cleanPath, h.webPublicDir) && !strings.HasPrefix(cleanPath, h.webDistDir) {
		http.NotFound(w, r)
		return
	}

	// Set content type based on extension
	ext := filepath.Ext(path)
	switch ext {
	case ".html":
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case ".css":
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case ".json":
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".map":
		w.Header().Set("Content-Type", "application/json")
	}

	http.ServeFile(w, r, path)
}

// fileExists checks if a file exists and is not a directory
func (h *StaticHandler) fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !info.IsDir()
}

// EmbeddedStaticHandler serves static files from an embedded filesystem
// Useful for production builds where files are embedded in the binary
type EmbeddedStaticHandler struct {
	fs     fs.FS
	prefix string
}

// NewEmbeddedStaticHandler creates a handler for embedded files
func NewEmbeddedStaticHandler(filesystem fs.FS, prefix string) *EmbeddedStaticHandler {
	return &EmbeddedStaticHandler{
		fs:     filesystem,
		prefix: prefix,
	}
}

// ServeHTTP serves embedded static files
func (h *EmbeddedStaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, h.prefix)
	if path == "" || path == "/" {
		path = "index.html"
	}
	path = strings.TrimPrefix(path, "/")

	http.ServeFileFS(w, r, h.fs, path)
}
