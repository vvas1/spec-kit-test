package api

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
)

// Handler is a function that handles a request and may return an error for centralized handling.
type Handler func(w http.ResponseWriter, r *http.Request) error

// Router wraps a mux with /api prefix, JSON middleware, CORS, error handling, and logging.
type Router struct {
	mux    *http.ServeMux   // top-level: /api/
	api    *http.ServeMux   // under /api: /issues, /users, etc.
	db     *mongo.Database
	logger *log.Logger
}

// NewRouter returns a new Router. Pass db as nil for tests; handlers will check.
func NewRouter(db *mongo.Database, logger *log.Logger) *Router {
	if logger == nil {
		logger = log.Default()
	}
	r := &Router{mux: http.NewServeMux(), api: http.NewServeMux(), db: db, logger: logger}
	r.mux.Handle("/api/", http.StripPrefix("/api", r.apiHandler()))
	return r
}

func (rt *Router) apiHandler() http.Handler {
	return rt.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// StripPrefix does not modify r.URL.Path; so path is still e.g. /api/issues. Rewrite for api mux.
		path := r.URL.Path
		if strings.HasPrefix(path, "/api") {
			path = path[len("/api"):]
			if path == "" {
				path = "/"
			}
		}
		r2 := r.Clone(r.Context())
		r2.URL = &url.URL{Path: path}
		rt.api.ServeHTTP(w, r2)
	}))
}

func (rt *Router) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// JSON
		if r.Header.Get("Content-Type") == "application/json" {
			// already set for reads
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		// Logging
		rt.logger.Printf("%s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

// ServeHTTP implements http.Handler. Routes are registered on /api/...
func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/api") {
		rt.writeError(w, http.StatusNotFound, "not found")
		return
	}
	rt.mux.ServeHTTP(w, r)
}

// writeError sends JSON { "error": "<message>" }.
func (rt *Router) writeError(w http.ResponseWriter, code int, message string) {
	WriteError(w, code, message)
}

// WriteError sends JSON { "error": "<message>" } for use by handlers.
func WriteError(w http.ResponseWriter, code int, message string) {
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// DB returns the database; may be nil in tests.
func (rt *Router) DB() *mongo.Database {
	return rt.db
}

// Handle registers a handler for the given pattern (relative to /api, e.g. "/issues").
func (rt *Router) Handle(pattern string, h Handler) {
	rt.api.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			rt.logger.Printf("handler error: %v", err)
			rt.writeError(w, http.StatusInternalServerError, err.Error())
		}
	})
}
