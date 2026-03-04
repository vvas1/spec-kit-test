package contract

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"issue-tracker/backend/internal/api"
	"issue-tracker/backend/internal/api/handler"
	"issue-tracker/backend/internal/store/mongodb"
)

func setupUsersRouter(t *testing.T) *api.Router {
	ctx := context.Background()
	client, err := mongodb.NewClient(ctx)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	db := mongodb.Database(client)
	router := api.NewRouter(db, nil)
	router.Handle("/users", handler.HandleGetUsers(router))
	return router
}

func TestGetUsers(t *testing.T) {
	router := setupUsersRouter(t)
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("GET /api/users: got status %d, want 200", w.Code)
	}
	// Expect {"items": [...]}
	var out struct {
		Items []interface{} `json:"items"`
	}
	if err := json.NewDecoder(w.Body).Decode(&out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.Items == nil {
		t.Error("items should be non-nil array")
	}
}
