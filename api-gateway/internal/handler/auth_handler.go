package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
	"google.golang.org/grpc/status"
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
	UserID       string `json:"user_id"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type changePasswordRequest struct {
	UserID      string `json:"user_id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ── Shared auth helpers used by other handlers ───────────────────────────────

// verifyBearerToken extracts and verifies the Bearer token from the request.
// On failure, writes the appropriate HTTP error and returns nil.
func verifyBearerToken(w http.ResponseWriter, r *http.Request, authClient *grpcclient.AuthClient) *authpb.VerifyTokenResponse {
	token := bearerToken(r.Header.Get("Authorization"))
	if token == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "authentication required"})
		return nil
	}
	resp, err := authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{AccessToken: token})
	if err != nil || !resp.GetValid() {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "invalid or expired token"})
		return nil
	}
	return resp
}

// requireAdminOrManager verifies token and requires ADMIN or MANAGER role.
// Returns the token response on success, nil on failure (HTTP error already written).
func requireAdminOrManager(w http.ResponseWriter, r *http.Request, authClient *grpcclient.AuthClient) *authpb.VerifyTokenResponse {
	claims := verifyBearerToken(w, r, authClient)
	if claims == nil {
		return nil
	}
	for _, role := range claims.GetRoles() {
		if role == "ADMIN" || role == "MANAGER" {
			return claims
		}
	}
	writeJSON(w, http.StatusForbidden, map[string]any{"success": false, "message": "insufficient permissions"})
	return nil
}

// requireStaff verifies token and requires any staff role (ADMIN/MANAGER/CHEF/WAITER).
func requireStaff(w http.ResponseWriter, r *http.Request, authClient *grpcclient.AuthClient) *authpb.VerifyTokenResponse {
	claims := verifyBearerToken(w, r, authClient)
	if claims == nil {
		return nil
	}
	for _, role := range claims.GetRoles() {
		switch role {
		case "ADMIN", "MANAGER", "CHEF", "WAITER":
			return claims
		}
	}
	writeJSON(w, http.StatusForbidden, map[string]any{"success": false, "message": "staff access required"})
	return nil
}

// requireAdmin verifies token and requires ADMIN role only.
func requireAdmin(w http.ResponseWriter, r *http.Request, authClient *grpcclient.AuthClient) *authpb.VerifyTokenResponse {
	claims := verifyBearerToken(w, r, authClient)
	if claims == nil {
		return nil
	}
	for _, role := range claims.GetRoles() {
		if role == "ADMIN" {
			return claims
		}
	}
	writeJSON(w, http.StatusForbidden, map[string]any{"success": false, "message": "admin access required"})
	return nil
}

// authBizError writes a 400 if the gRPC response carries success=false.
// Returns true if an error was written (caller should return immediately).
func authBizError(w http.ResponseWriter, success bool, message string) bool {
	if !success {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": message})
		return true
	}
	return false
}

// grpcErrMsg extracts the human-readable description from a gRPC status error,
// stripping the "rpc error: code = X desc = " prefix.
func grpcErrMsg(err error) string {
	if s, ok := status.FromError(err); ok {
		return s.Message()
	}
	return err.Error()
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
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": grpcErrMsg(err)})
		return
	}
	if authBizError(w, resp.GetSuccess(), resp.GetMessage()) {
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
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": grpcErrMsg(err)})
		return
	}
	if authBizError(w, resp.GetSuccess(), resp.GetMessage()) {
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
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": grpcErrMsg(err)})
		return
	}
	if authBizError(w, resp.GetSuccess(), resp.GetMessage()) {
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
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": grpcErrMsg(err)})
		return
	}
	if !resp.GetValid() {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"success": false, "message": "token is invalid or expired"})
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
	if req.RefreshToken == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "refresh_token is required"})
		return
	}

	resp, err := h.authClient.Logout(r.Context(), &authpb.LogoutRequest{
		UserId:       req.UserID,
		AccessToken:  req.AccessToken,
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": grpcErrMsg(err)})
		return
	}
	if authBizError(w, resp.GetSuccess(), resp.GetMessage()) {
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
		writeJSON(w, http.StatusBadRequest, map[string]any{"success": false, "message": grpcErrMsg(err)})
		return
	}
	if authBizError(w, resp.GetSuccess(), resp.GetMessage()) {
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
