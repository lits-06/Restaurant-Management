package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	staffpb "restaurant-management/proto/staff"
)

// StaffHandler exposes HTTP endpoints backed by staff-service gRPC methods.
type StaffHandler struct {
	staffClient *grpcclient.StaffClient
	authClient  *grpcclient.AuthClient
}

func NewStaffHandler(staffClient *grpcclient.StaffClient, authClient *grpcclient.AuthClient) *StaffHandler {
	return &StaffHandler{staffClient: staffClient, authClient: authClient}
}

type createStaffRequest struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Contact string `json:"contact"`
	Avatar  string `json:"avatar"`
}

type updateStaffRequest struct {
	Name    string `json:"name"`
	Role    string `json:"role"`
	Contact string `json:"contact"`
	Avatar  string `json:"avatar"`
}

func (h *StaffHandler) CreateStaff(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	var req createStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.staffClient.CreateStaff(r.Context(), &staffpb.CreateStaffRequest{
		Name:    req.Name,
		Role:    req.Role,
		Contact: req.Contact,
		Avatar:  req.Avatar,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *StaffHandler) ListStaff(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	resp, err := h.staffClient.ListStaff(r.Context(), &staffpb.ListStaffRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 20),
		Keyword:  strings.TrimSpace(r.URL.Query().Get("keyword")),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *StaffHandler) GetStaff(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	staffID := extractStaffIDFromPath(r.URL.Path)
	if staffID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "staff_id is required"})
		return
	}

	resp, err := h.staffClient.GetStaff(r.Context(), &staffpb.GetStaffRequest{StaffId: staffID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *StaffHandler) UpdateStaff(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	staffID := extractStaffIDFromPath(r.URL.Path)
	if staffID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "staff_id is required"})
		return
	}

	var req updateStaffRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.staffClient.UpdateStaff(r.Context(), &staffpb.UpdateStaffRequest{
		StaffId: staffID,
		Name:    req.Name,
		Role:    req.Role,
		Contact: req.Contact,
		Avatar:  req.Avatar,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *StaffHandler) DeleteStaff(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	staffID := extractStaffIDFromPath(r.URL.Path)
	if staffID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "staff_id is required"})
		return
	}

	resp, err := h.staffClient.DeleteStaff(r.Context(), &staffpb.DeleteStaffRequest{StaffId: staffID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// func (h *StaffHandler) verifyAuthorization(r *http.Request) (*authpb.VerifyTokenResponse, error) {
// 	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
// 	if authHeader == "" {
// 		return nil, errUnauthorized("missing Authorization header")
// 	}

// 	token := extractBearerToken(authHeader)
// 	if token == "" {
// 		return nil, errUnauthorized("invalid Authorization header, expected Bearer token")
// 	}

// 	resp, err := h.authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{AccessToken: token})
// 	if err != nil {
// 		return nil, errUnauthorized("failed to verify access token")
// 	}
// 	if !resp.Valid {
// 		return nil, errUnauthorized("invalid or expired access token")
// 	}

// 	return resp, nil
// }

func extractStaffIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/staff/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}
	return strings.TrimSpace(segments[0])
}
