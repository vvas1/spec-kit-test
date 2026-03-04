package contract

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"issue-tracker/backend/internal/api"
	"issue-tracker/backend/internal/api/handler"
	"issue-tracker/backend/internal/store/mongodb"
)

func setupRouter(t *testing.T) *api.Router {
	ctx := context.Background()
	client, err := mongodb.NewClient(ctx)
	if err != nil {
		t.Skipf("MongoDB not available: %v", err)
	}
	db := mongodb.Database(client)
	router := api.NewRouter(db, nil)
	router.Handle("/issues", handler.HandlePostIssues(router))
	return router
}

func TestPostIssues(t *testing.T) {
	router := setupRouter(t)
	// 400 missing title
	body := bytes.NewBufferString(`{}`)
	req := httptest.NewRequest(http.MethodPost, "/api/issues", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("POST /api/issues {}: got status %d, want 400", w.Code)
	}
	// 400 title too long
	longTitle := string(make([]byte, 201))
	body2 := bytes.NewBufferString(`{"title":"` + longTitle + `"}`)
	req2 := httptest.NewRequest(http.MethodPost, "/api/issues", body2)
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	if w2.Code != http.StatusBadRequest {
		t.Errorf("POST /api/issues long title: got status %d, want 400", w2.Code)
	}
	// 201 with body
	body3 := bytes.NewBufferString(`{"title":"Test","description":"D"}`)
	req3 := httptest.NewRequest(http.MethodPost, "/api/issues", body3)
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	if w3.Code != http.StatusCreated {
		t.Errorf("POST /api/issues valid: got status %d, want 201", w3.Code)
	}
	var issue map[string]interface{}
	if err := json.NewDecoder(w3.Body).Decode(&issue); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if issue["title"] != "Test" || issue["status"] != "To Do" {
		t.Errorf("unexpected body: %v", issue)
	}
}
