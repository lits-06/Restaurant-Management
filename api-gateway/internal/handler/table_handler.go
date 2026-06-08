package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
	tablepb "restaurant-management/proto/table"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// TableHandler exposes HTTP endpoints backed by table-service gRPC methods.
type TableHandler struct {
	tableClient *grpcclient.TableClient
	authClient  *grpcclient.AuthClient
}

func NewTableHandler(tableClient *grpcclient.TableClient, authClient *grpcclient.AuthClient) *TableHandler {
	return &TableHandler{
		tableClient: tableClient,
		authClient:  authClient,
	}
}

type createTableRequest struct {
	TableNumber string `json:"table_number"`
	Capacity    int32  `json:"capacity"`
	Location    string `json:"location"`
}

type updateTableRequest struct {
	TableNumber string `json:"table_number"`
	Capacity    int32  `json:"capacity"`
	Location    string `json:"location"`
}

type updateTableStatusRequest struct {
	Status  string `json:"status"`
	OrderID string `json:"order_id"`
}

type reservationItemRequest struct {
	MenuItemID string `json:"menu_item_id"`
	Quantity   int32  `json:"quantity"`
	Note       string `json:"note"`
}

type createReservationRequest struct {
	CustomerName  string                   `json:"customer_name"`
	CustomerPhone string                   `json:"customer_phone"`
	Notes         string                   `json:"notes"`
	StartTime     string                   `json:"start_time"`
	EndTime       string                   `json:"end_time"`
	Items         []reservationItemRequest `json:"items"`
}

func (h *TableHandler) CreateTable(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
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
		Location:    req.Location,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) ListTables(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	status, err := parseTableStatus(r.URL.Query().Get("status"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	resp, rpcErr := h.tableClient.ListTables(r.Context(), &tablepb.ListTablesRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 10),
		Status:   status,
		Location: strings.TrimSpace(r.URL.Query().Get("location")),
	})
	if rpcErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": rpcErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) GetTable(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	tableID := extractTableIDFromPath(r.URL.Path)
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

func (h *TableHandler) UpdateTable(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	tableID := extractTableIDFromPath(r.URL.Path)
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
		Location:    req.Location,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) DeleteTable(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	tableID := extractTableIDFromPath(r.URL.Path)
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

func (h *TableHandler) UpdateTableStatus(w http.ResponseWriter, r *http.Request) {
	if !allowPatch(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	tableID := extractTableIDForStatusPath(r.URL.Path)
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}

	var req updateTableStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	status, err := parseTableStatus(req.Status)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	resp, rpcErr := h.tableClient.UpdateTableStatus(r.Context(), &tablepb.UpdateTableStatusRequest{
		TableId: tableID,
		Status:  status,
		OrderId: req.OrderID,
	})
	if rpcErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": rpcErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) GetAvailableTables(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	minCapacity := int32(0)
	rawMinCapacity := strings.TrimSpace(r.URL.Query().Get("min_capacity"))
	if rawMinCapacity != "" {
		parsed, err := strconv.ParseInt(rawMinCapacity, 10, 32)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid min_capacity"})
			return
		}
		minCapacity = int32(parsed)
	}

	resp, err := h.tableClient.GetAvailableTables(r.Context(), &tablepb.GetAvailableTablesRequest{
		MinCapacity: minCapacity,
		Location:    strings.TrimSpace(r.URL.Query().Get("location")),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	tableID := extractTableIDForReservationPath(r.URL.Path)
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}

	var req createReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	startTime, err := parseRFC3339Time(req.StartTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid start_time"})
		return
	}
	endTime, err := parseRFC3339Time(req.EndTime)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid end_time"})
		return
	}

	items := make([]*tablepb.ReservationItem, 0, len(req.Items))
	for _, item := range req.Items {
		items = append(items, &tablepb.ReservationItem{
			MenuItemId: item.MenuItemID,
			Quantity:   item.Quantity,
			Note:       item.Note,
		})
	}

	resp, rpcErr := h.tableClient.CreateReservation(r.Context(), &tablepb.CreateReservationRequest{
		TableId:       tableID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		Notes:         req.Notes,
		StartTime:     timestamppb.New(startTime),
		EndTime:       timestamppb.New(endTime),
		Items:         items,
	})
	if rpcErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": rpcErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) ListReservations(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	tableID := extractTableIDForReservationPath(r.URL.Path)
	if tableID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "table_id is required"})
		return
	}

	status, err := parseReservationStatus(r.URL.Query().Get("status"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	fromTime, err := parseOptionalRFC3339Time(r.URL.Query().Get("from"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid from"})
		return
	}
	toTime, err := parseOptionalRFC3339Time(r.URL.Query().Get("to"))
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid to"})
		return
	}

	req := &tablepb.ListReservationsRequest{
		TableId:  tableID,
		Status:   status,
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 10),
	}
	if fromTime != nil {
		req.FromTime = timestamppb.New(*fromTime)
	}
	if toTime != nil {
		req.ToTime = timestamppb.New(*toTime)
	}

	resp, rpcErr := h.tableClient.ListReservations(r.Context(), req)
	if rpcErr != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": rpcErr.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	reservationID := extractReservationIDFromPath(r.URL.Path)
	if reservationID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "reservation_id is required"})
		return
	}

	resp, err := h.tableClient.GetReservation(r.Context(), &tablepb.GetReservationRequest{ReservationId: reservationID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) CancelReservation(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	if _, err := h.verifyAuthorization(r); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}

	reservationID := extractReservationIDForCancelPath(r.URL.Path)
	if reservationID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "reservation_id is required"})
		return
	}

	resp, err := h.tableClient.CancelReservation(r.Context(), &tablepb.CancelReservationRequest{ReservationId: reservationID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *TableHandler) verifyAuthorization(r *http.Request) (*authpb.VerifyTokenResponse, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return nil, errUnauthorized("missing Authorization header")
	}

	token := extractBearerToken(authHeader)
	if token == "" {
		return nil, errUnauthorized("invalid Authorization header, expected Bearer token")
	}

	resp, err := h.authClient.VerifyToken(r.Context(), &authpb.VerifyTokenRequest{AccessToken: token})
	if err != nil {
		return nil, errUnauthorized("failed to verify access token")
	}
	if !resp.Valid {
		return nil, errUnauthorized("invalid or expired access token")
	}

	return resp, nil
}

func parseTableStatus(raw string) (tablepb.TableStatus, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "", "0", "unknown", "status_unknown":
		return tablepb.TableStatus_STATUS_UNKNOWN, nil
	case "1", "available", "status_available":
		return tablepb.TableStatus_STATUS_AVAILABLE, nil
	case "2", "occupied", "status_occupied":
		return tablepb.TableStatus_STATUS_OCCUPIED, nil
	case "3", "reserved", "status_reserved":
		return tablepb.TableStatus_STATUS_RESERVED, nil
	case "4", "cleaning", "status_cleaning":
		return tablepb.TableStatus_STATUS_CLEANING, nil
	case "5", "out_of_service", "status_out_of_service":
		return tablepb.TableStatus_STATUS_OUT_OF_SERVICE, nil
	default:
		return tablepb.TableStatus_STATUS_UNKNOWN, errors.New("invalid table status")
	}
}

func parseReservationStatus(raw string) (tablepb.ReservationStatus, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	switch value {
	case "", "0", "unknown", "reservation_status_unknown":
		return tablepb.ReservationStatus_RESERVATION_STATUS_UNKNOWN, nil
	case "1", "reserved", "reservation_status_reserved":
		return tablepb.ReservationStatus_RESERVATION_STATUS_RESERVED, nil
	case "2", "cancelled", "canceled", "reservation_status_cancelled":
		return tablepb.ReservationStatus_RESERVATION_STATUS_CANCELLED, nil
	case "3", "completed", "reservation_status_completed":
		return tablepb.ReservationStatus_RESERVATION_STATUS_COMPLETED, nil
	default:
		return tablepb.ReservationStatus_RESERVATION_STATUS_UNKNOWN, errors.New("invalid reservation status")
	}
}

func extractTableIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/tables/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}
	return strings.TrimSpace(segments[0])
}

func extractTableIDForStatusPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/tables/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) != 2 || strings.TrimSpace(segments[1]) != "status" {
		return ""
	}

	return strings.TrimSpace(segments[0])
}

func extractTableIDForReservationPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/tables/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) != 2 || strings.TrimSpace(segments[1]) != "reservations" {
		return ""
	}

	return strings.TrimSpace(segments[0])
}

func extractReservationIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/reservations/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}

	return strings.TrimSpace(segments[0])
}

func extractReservationIDForCancelPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/reservations/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) != 2 || strings.TrimSpace(segments[1]) != "cancel" {
		return ""
	}

	return strings.TrimSpace(segments[0])
}

func parseRFC3339Time(value string) (time.Time, error) {
	return time.Parse(time.RFC3339, strings.TrimSpace(value))
}

func parseOptionalRFC3339Time(value string) (*time.Time, error) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil, nil
	}

	parsed, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, err
	}

	return &parsed, nil
}
