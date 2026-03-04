package handler

import (
	"encoding/json"
	"net/http"

	"issue-tracker/backend/internal/api"
	"issue-tracker/backend/internal/model"
	"issue-tracker/backend/internal/service"
	"issue-tracker/backend/internal/store/mongodb"
)

// POST /issues
func HandlePostIssues(rt *api.Router) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodPost {
			api.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return nil
		}
		db := rt.DB()
		if db == nil {
			api.WriteError(w, http.StatusInternalServerError, "database not configured")
			return nil
		}
		var in model.CreateIssueInput
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			api.WriteError(w, http.StatusBadRequest, "invalid JSON")
			return nil
		}
		svc := service.NewIssueService(mongodb.NewIssueStore(db))
		issue, err := svc.Create(r.Context(), &in)
		if err != nil {
			api.WriteError(w, http.StatusBadRequest, err.Error())
			return nil
		}
		w.WriteHeader(http.StatusCreated)
		return json.NewEncoder(w).Encode(issue)
	}
}
