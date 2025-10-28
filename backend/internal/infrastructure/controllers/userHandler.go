package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/domain/user"
	"github.com/mattadlerpdx/estatesalefinder-ai/backend/internal/infrastructure/middleware"
)

type UserHandler struct {
	service *user.Service
}

// NewUserHandler returns a pointer to a new UserHandler
func NewUserHandler(service *user.Service) *UserHandler {
	return &UserHandler{service: service}
}

type EnsureUserResponse struct {
	Message string `json:"message"`
}

// EnsureUser handles POST /auth/ensureUser
func (h *UserHandler) EnsureUser(w http.ResponseWriter, r *http.Request) {
	uid, ok := r.Context().Value(middleware.ContextKeyUID).(string)
	if !ok || uid == "" {
		log.Println("Missing or invalid Firebase UID in context.")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	log.Printf("Validating Firebase user with UID: %s", uid)

	err := h.service.EnsureUserExists(uid)
	if err != nil {
		log.Printf("Failed to ensure user with UID %s: %v", uid, err)
		http.Error(w, "Failed to ensure user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Firebase user with UID %s ensured in database.", uid)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(EnsureUserResponse{Message: "User ensured in DB."})
}
