package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	tablepb "restaurant-management/proto/table"
)

// TableHandler exposes HTTP endpoints backed by table-service gRPC methods.
type TableHandler struct {
	tableClient *grpcclient.TableClient
	authClient  *grpcclient.AuthClient
}

func NewTableHandler(tableClient *grpcclient.TableClient, authClient *grpcclient.AuthClient) *TableHandler {
	return &TableHandler{tableClient: tableClient, authClient: authClient}
}

// --- Request structs ---

type createTableRequest struct {
	TableNumber int32 `json:"table_number"`
	Capacity    int32 `json:"capacity"`
}

type updateTableRequest struct {
	TableNumber int32 `json:"table_number"`
	Capacity    int32 `json:"capacity"`
}

type updateTableStatusRequest struct {
	Status string `json:"status"`
}

// --- Handlers ---

// CreateTable handles POST /tables.
func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	var req createTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.tableClient.CreateTable(r.Context(), &tablepb.CreateTableRequest{
		TableNumber: req.TableNumber,
		Capacity:    req.Capacity,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, resp)
}

// ListTables handles GET /tables.
func (h *TableHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	resp, err := h.tableClient.ListTables(r.Context(), &tablepb.ListTablesRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 20),
		Status:   parseTableStatus(r.URL.Query().Get("status")),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetTable handles GET /tables/{id}.
func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	tableID := extractIDFromPath(r.URL.Path, "/tables/")
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}
	resp, err := h.tableClient.GetTable(r.Context(), &tablepb.GetTableRequest{TableId: tableID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateTable handles PUT /tables/{id}.
func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}
	tableID := extractIDFromPath(r.URL.Path, "/tables/")
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}
	var req updateTableRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.tableClient.UpdateTable(r.Context(), &tablepb.UpdateTableRequest{
		TableId:     tableID,
		TableNumber: req.TableNumber,
		Capacity:    req.Capacity,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// DeleteTable handles DELETE /tables/{id}.
func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	tableID := extractIDFromPath(r.URL.Path, "/tables/")
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}
	resp, err := h.tableClient.DeleteTable(r.Context(), &tablepb.DeleteTableRequest{TableId: tableID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// UpdateTableStatus handles PATCH /tables/{id}/status.
func (h *TableHandler) UpdateTableStatus(w http.ResponseWriter, r *http.Request) {
	if !allowPatch(w, r) {
		return
	}
	path := strings.TrimSuffix(r.URL.Path, "/status")
	tableID := extractIDFromPath(path, "/tables/")
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}
	var req updateTableStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.tableClient.UpdateTableStatus(r.Context(), &tablepb.UpdateTableStatusRequest{
		TableId: tableID,
		Status:  parseTableStatus(req.Status),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// GetAvailableTables handles GET /tables/available.
func (h *TableHandler) GetAvailableTables(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	resp, err := h.tableClient.GetAvailableTables(r.Context(), &tablepb.GetAvailableTablesRequest{
		MinCapacity: parseInt32Query(r, "min_capacity", 0),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

// --- Helpers ---

func extractIDFromPath(path, prefix string) string {
	trimmed := strings.TrimPrefix(path, prefix)
	trimmed = strings.Trim(trimmed, "/")
	if trimmed == "" {
		return ""
	}
	return strings.Split(trimmed, "/")[0]
}

func parseTableStatus(s string) tablepb.TableStatus {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case "AVAILABLE":
		return tablepb.TableStatus_STATUS_AVAILABLE
	case "CLEANING":
		return tablepb.TableStatus_STATUS_CLEANING
	case "OUT_OF_SERVICE":
		return tablepb.TableStatus_STATUS_OUT_OF_SERVICE
	default:
		return tablepb.TableStatus_STATUS_UNKNOWN
	}
}
