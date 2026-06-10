package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	userpb "restaurant-management/proto/user"
)

type UserHandler struct {
	userClient *grpcclient.UserClient
	authClient *grpcclient.AuthClient
}

func NewUserHandler(userClient *grpcclient.UserClient, authClient *grpcclient.AuthClient) *UserHandler {
	return &UserHandler{userClient: userClient, authClient: authClient}
}

type createUserRequest struct {
	Email    string   `json:"email"`
	Username string   `json:"username"`
	FullName string   `json:"full_name"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Roles    []string `json:"roles"`
}

type updateUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Status   string `json:"status"`
}

type assignRoleRequest struct {
	Roles []string `json:"roles"`
}

type userChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.userClient.CreateUser(r.Context(), &userpb.CreateUserRequest{
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Phone:    req.Phone,
		Password: req.Password,
		Roles:    parseProtoRoles(req.Roles),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	resp, err := h.userClient.ListUsers(r.Context(), &userpb.ListUsersRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 20),
		Keyword:  strings.TrimSpace(r.URL.Query().Get("keyword")),
		Status:   parseProtoUserStatus(r.URL.Query().Get("status")),
		Role:     parseProtoUserRole(r.URL.Query().Get("role")),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	userID := extractUserIDFromPath(r.URL.Path)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id is required"})
		return
	}
	resp, err := h.userClient.GetUser(r.Context(), &userpb.GetUserRequest{UserId: userID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	email := strings.TrimSpace(r.URL.Query().Get("email"))
	if email == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "email query param is required"})
		return
	}
	resp, err := h.userClient.GetUserByEmail(r.Context(), &userpb.GetUserByEmailRequest{Email: email})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}
	userID := extractUserIDFromPath(r.URL.Path)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id is required"})
		return
	}
	var req updateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.userClient.UpdateUser(r.Context(), &userpb.UpdateUserRequest{
		UserId:   userID,
		Email:    req.Email,
		Username: req.Username,
		FullName: req.FullName,
		Phone:    req.Phone,
		Status:   parseProtoUserStatus(req.Status),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	userID := extractUserIDFromPath(r.URL.Path)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id is required"})
		return
	}
	resp, err := h.userClient.DeleteUser(r.Context(), &userpb.DeleteUserRequest{UserId: userID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) AssignRole(w http.ResponseWriter, r *http.Request) {
	if !allowPatch(w, r) {
		return
	}
	userID := extractUserIDFromPath(r.URL.Path)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id is required"})
		return
	}
	var req assignRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.userClient.AssignRole(r.Context(), &userpb.AssignRoleRequest{
		UserId: userID,
		Roles:  parseProtoRoles(req.Roles),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) GetUserRoles(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	userID := extractUserIDFromPath(r.URL.Path)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id is required"})
		return
	}
	resp, err := h.userClient.GetUserRoles(r.Context(), &userpb.GetUserRolesRequest{UserId: userID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if !allowPatch(w, r) {
		return
	}
	userID := extractUserIDFromPath(r.URL.Path)
	if userID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "user_id is required"})
		return
	}
	var req userChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.userClient.ChangePassword(r.Context(), &userpb.ChangePasswordRequest{
		UserId:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// ── helpers ───────────────────────────────────────────────────

func extractUserIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/users/"), " /")
	if trimmed == "" {
		return ""
	}
	return strings.Split(trimmed, "/")[0]
}

func parseProtoRoles(roles []string) []userpb.UserRole {
	out := make([]userpb.UserRole, 0, len(roles))
	for _, r := range roles {
		switch strings.ToUpper(r) {
		case "USER":
			out = append(out, userpb.UserRole_ROLE_USER)
		case "MANAGER":
			out = append(out, userpb.UserRole_ROLE_MANAGER)
		case "CHEF":
			out = append(out, userpb.UserRole_ROLE_CHEF)
		case "WAITER":
			out = append(out, userpb.UserRole_ROLE_WAITER)
		case "ADMIN":
			out = append(out, userpb.UserRole_ROLE_ADMIN)
		}
	}
	return out
}

func parseProtoUserRole(r string) userpb.UserRole {
	switch strings.ToUpper(r) {
	case "USER":
		return userpb.UserRole_ROLE_USER
	case "MANAGER":
		return userpb.UserRole_ROLE_MANAGER
	case "CHEF":
		return userpb.UserRole_ROLE_CHEF
	case "WAITER":
		return userpb.UserRole_ROLE_WAITER
	case "ADMIN":
		return userpb.UserRole_ROLE_ADMIN
	default:
		return userpb.UserRole_ROLE_UNKNOWN
	}
}

func parseProtoUserStatus(s string) userpb.UserStatus {
	switch strings.ToUpper(s) {
	case "ACTIVE":
		return userpb.UserStatus_STATUS_ACTIVE
	case "INACTIVE":
		return userpb.UserStatus_STATUS_INACTIVE
	case "SUSPENDED":
		return userpb.UserStatus_STATUS_SUSPENDED
	default:
		return userpb.UserStatus_STATUS_UNKNOWN
	}
}
