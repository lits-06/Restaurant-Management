package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
)

// AuthHandler exposes HTTP endpoints backed by auth-service gRPC methods.
type AuthHandler struct {
	authClient *grpcclient.AuthClient
}

func NewAuthHandler(authClient *grpcclient.AuthClient) *AuthHandler {
	return &AuthHandler{authClient: authClient}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type verifyTokenRequest struct {
	AccessToken string `json:"access_token"`
}

type logoutRequest struct {
	UserID      string `json:"user_id"`
	AccessToken string `json:"access_token"`
}

type changePasswordRequest struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	var req registerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.Register(r.Context(), &authpb.RegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Username: req.Username,
		FullName: req.FullName,
		Phone:    req.Phone,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.Login(r.Context(), &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	var req refreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.RefreshToken(r.Context(), &authpb.RefreshTokenRequest{RefreshToken: req.RefreshToken})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) VerifyToken(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	var req verifyTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{AccessToken: req.AccessToken})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	var req logoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	if req.AccessToken == "" {
		req.AccessToken = bearerToken(r.Header.Get("Authorization"))
	}
	if req.UserID == "" && req.AccessToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id or access_token is required"})
		return
	}

	resp, err := h.authClient.Logout(r.Context(), &authpb.LogoutRequest{UserId: req.UserID, AccessToken: req.AccessToken})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	var req changePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.authClient.ChangePassword(r.Context(), &authpb.ChangePasswordRequest{
		UserId:      req.UserID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func allowPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		writeCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		return false
	}

	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return false
	}

	return true
}

func bearerToken(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	const prefix = "Bearer "
	if strings.HasPrefix(value, prefix) {
		return strings.TrimSpace(value[len(prefix):])
	}
	return ""
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	writeCORSHeaders(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if statusCode >= http.StatusBadRequest {
		if m, ok := payload.(map[string]any); ok {
			if _, exists := m["success"]; !exists {
				m["success"] = false
			}
		}
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.Header().Set("Access-Control-Allow-Methods", strings.Join([]string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions}, ", "))
}

func allowGet(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		writeCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		return false
	}

	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return false
	}

	return true
}

func allowPut(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		writeCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		return false
	}

	if r.Method != http.MethodPut {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return false
	}

	return true
}

func allowDelete(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		writeCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		return false
	}

	if r.Method != http.MethodDelete {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return false
	}

	return true
}

func allowPatch(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		writeCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		return false
	}

	if r.Method != http.MethodPatch {
		writeJSON(w, http.StatusMethodNotAllowed, map[string]any{"error": "method not allowed"})
		return false
	}

	return true
}
