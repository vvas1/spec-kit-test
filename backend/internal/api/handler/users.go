package handler

import (
	"encoding/json"
	"net/http"

	"issue-tracker/backend/internal/api"
	"issue-tracker/backend/internal/store/mongodb"
)

// GET /users
func HandleGetUsers(rt *api.Router) api.Handler {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != http.MethodGet {
			api.WriteError(w, http.StatusMethodNotAllowed, "method not allowed")
			return nil
		}
		db := rt.DB()
		if db == nil {
			api.WriteError(w, http.StatusInternalServerError, "database not configured")
			return nil
		}
		store := mongodb.NewUserStore(db)
		users, err := store.List(r.Context())
		if err != nil {
			api.WriteError(w, http.StatusInternalServerError, err.Error())
			return nil
		}
		return json.NewEncoder(w).Encode(map[string]interface{}{"items": users})
	}
}
