package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/pseudoerr/auth-service/internal/models"
	"github.com/pseudoerr/auth-service/internal/service"
	"github.com/pseudoerr/auth-service/internal/validation"

	"log/slog"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Register user
	authResponse, err := h.authService.Register(&req)
	if err != nil {
		slog.Error("Registration failed", "error", err, "email", req.Email)
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.Info("User registered successfully", "user_id", authResponse.User.ID, "email", authResponse.User.Email)
	h.writeJSON(w, http.StatusCreated, authResponse)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Login user
	authResponse, err := h.authService.Login(&req)
	if err != nil {
		slog.Error("Login failed", "error", err, "email", req.Email)
		h.writeError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	slog.Info("User logged in successfully", "user_id", authResponse.User.ID, "email", authResponse.User.Email)
	h.writeJSON(w, http.StatusOK, authResponse)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req models.RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	// Validate request
	if err := validation.ValidateStruct(&req); err != nil {
		h.writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Refresh token
	authResponse, err := h.authService.RefreshToken(&req)
	if err != nil {
		slog.Error("Token refresh failed", "error", err)
		h.writeError(w, http.StatusUnauthorized, "Invalid refresh token")
		return
	}

	slog.Info("Token refreshed successfully", "user_id", authResponse.User.ID)
	h.writeJSON(w, http.StatusOK, authResponse)
}

func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT middleware context
	userID, err := getUserIDFromContext(r)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Get user profile
	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		slog.Error("Failed to get user profile", "error", err, "user_id", userID)
		h.writeError(w, http.StatusNotFound, "User not found")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// Get user ID from JWT middleware context
	userID, err := getUserIDFromContext(r)
	if err != nil {
		h.writeError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	// Optional: get refresh token from request body to logout specific session
	var req struct {
		RefreshToken string `json:"refresh_token,omitempty"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Logout user
	if err := h.authService.Logout(userID, req.RefreshToken); err != nil {
		slog.Error("Logout failed", "error", err, "user_id", userID)
		h.writeError(w, http.StatusInternalServerError, "Logout failed")
		return
	}

	slog.Info("User logged out successfully", "user_id", userID)
	h.writeJSON(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *AuthHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *AuthHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{"error": message})
}

func getUserIDFromContext(r *http.Request) (int, error) {
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		return 0, fmt.Errorf("user ID not found in context")
	}
	return strconv.Atoi(userIDStr)
}
