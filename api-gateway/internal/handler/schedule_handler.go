package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
	schedulepb "restaurant-management/proto/schedule"
)

type ScheduleHandler struct {
	scheduleClient *grpcclient.ScheduleClient
	authClient     *grpcclient.AuthClient
}

func NewScheduleHandler(scheduleClient *grpcclient.ScheduleClient, authClient *grpcclient.AuthClient) *ScheduleHandler {
	return &ScheduleHandler{scheduleClient: scheduleClient, authClient: authClient}
}

type createShiftRequest struct {
	UserID    string `json:"user_id"`
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Role      string `json:"role"`
	Notes     string `json:"notes"`
}

type updateShiftRequest struct {
	Date      string `json:"date"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Notes     string `json:"notes"`
}

func (h *ScheduleHandler) verifyCaller(r *http.Request) (*callerInfo, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return nil, errUnauthorized("missing Authorization header")
	}
	token := extractBearerToken(authHeader)
	if token == "" {
		return nil, errUnauthorized("invalid Authorization header")
	}
	resp, err := h.authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{AccessToken: token})
	if err != nil || !resp.Valid {
		return nil, errUnauthorized("invalid or expired token")
	}
	return &callerInfo{UserID: resp.UserId, Roles: resp.Roles}, nil
}

// ListShifts — GET /schedule/shifts — auth required, staff only
func (h *ScheduleHandler) ListShifts(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}
	if !hasStaffRole(caller.Roles) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "staff role required"})
		return
	}

	resp, err := h.scheduleClient.ListShifts(r.Context(), &schedulepb.ListShiftsRequest{
		Month:    r.URL.Query().Get("month"),
		UserId:   r.URL.Query().Get("user_id"),
		Role:     r.URL.Query().Get("role"),
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 50),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// CreateShift — POST /schedule/shifts — auth required; can create for self always; for others requires ADMIN/MANAGER
func (h *ScheduleHandler) CreateShift(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	var req createShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	// If user_id omitted, default to caller
	if req.UserID == "" {
		req.UserID = caller.UserID
	}

	// Creating for another user requires ADMIN or MANAGER
	if req.UserID != caller.UserID && !hasAdminOrManager(caller.Roles) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "only ADMIN or MANAGER can create shifts for other users"})
		return
	}

	resp, err := h.scheduleClient.CreateShift(r.Context(), &schedulepb.CreateShiftRequest{
		UserId:    req.UserID,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Role:      req.Role,
		Notes:     req.Notes,
		CreatedBy: caller.UserID,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetShift — GET /schedule/shifts/{id} — auth required; owner or ADMIN/MANAGER
func (h *ScheduleHandler) GetShift(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	shiftID := extractShiftIDFromPath(r.URL.Path)
	if shiftID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "shift_id is required"})
		return
	}

	resp, err := h.scheduleClient.GetShift(r.Context(), &schedulepb.GetShiftRequest{ShiftId: shiftID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	if resp.Shift.UserId != caller.UserID && !hasAdminOrManager(caller.Roles) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "access denied"})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateShift — PUT /schedule/shifts/{id} — auth required; owner or ADMIN/MANAGER
func (h *ScheduleHandler) UpdateShift(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	shiftID := extractShiftIDFromPath(r.URL.Path)
	if shiftID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "shift_id is required"})
		return
	}

	// Fetch first to check ownership
	getResp, err := h.scheduleClient.GetShift(r.Context(), &schedulepb.GetShiftRequest{ShiftId: shiftID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if getResp.Shift.UserId != caller.UserID && !hasAdminOrManager(caller.Roles) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "access denied"})
		return
	}

	var req updateShiftRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.scheduleClient.UpdateShift(r.Context(), &schedulepb.UpdateShiftRequest{
		ShiftId:   shiftID,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Notes:     req.Notes,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// DeleteShift — DELETE /schedule/shifts/{id} — auth required; owner or ADMIN/MANAGER
func (h *ScheduleHandler) DeleteShift(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	shiftID := extractShiftIDFromPath(r.URL.Path)
	if shiftID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "shift_id is required"})
		return
	}

	// Fetch first to check ownership
	getResp, err := h.scheduleClient.GetShift(r.Context(), &schedulepb.GetShiftRequest{ShiftId: shiftID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if getResp.Shift.UserId != caller.UserID && !hasAdminOrManager(caller.Roles) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "access denied"})
		return
	}

	resp, err := h.scheduleClient.DeleteShift(r.Context(), &schedulepb.DeleteShiftRequest{ShiftId: shiftID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func hasAdminOrManager(roles []string) bool {
	for _, r := range roles {
		if r == "ADMIN" || r == "MANAGER" {
			return true
		}
	}
	return false
}

func extractShiftIDFromPath(path string) string {
	// /schedule/shifts/{id}
	trimmed := strings.TrimPrefix(path, "/schedule/shifts/")
	trimmed = strings.Trim(trimmed, "/")
	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 || segments[0] == "" {
		return ""
	}
	return segments[0]
}
