package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
	orderpb "restaurant-management/proto/order"
)

type OrderHandler struct {
	orderClient *grpcclient.OrderClient
	authClient  *grpcclient.AuthClient
}

func NewOrderHandler(orderClient *grpcclient.OrderClient, authClient *grpcclient.AuthClient) *OrderHandler {
	return &OrderHandler{orderClient: orderClient, authClient: authClient}
}

type orderItemRequest struct {
	ItemID   string  `json:"item_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int32   `json:"quantity"`
}

type createOrderRequest struct {
	TableID   string             `json:"table_id"`
	Name      string             `json:"name"`
	Phone     string             `json:"phone"`
	Notes     string             `json:"notes"`
	Time      string             `json:"time"`
	EndTime   string             `json:"end_time"`
	Date      string             `json:"date"`
	PartySize int32              `json:"party_size"`
	Status    string             `json:"status"`
	Items     []orderItemRequest `json:"items"`
}

type updateOrderRequest struct {
	TableID   string             `json:"table_id"`
	Name      string             `json:"name"`
	Phone     string             `json:"phone"`
	Notes     string             `json:"notes"`
	Time      string             `json:"time"`
	EndTime   string             `json:"end_time"`
	Date      string             `json:"date"`
	PartySize int32              `json:"party_size"`
	Status    string             `json:"status"`
	Items     []orderItemRequest `json:"items"`
}

type cancelOrderRequest struct {
	Reason string `json:"reason"`
}

type updateOrderStatusRequest struct {
	Status string `json:"status"`
}

type addOrderItemRequest struct {
	Item orderItemRequest `json:"item"`
}

type updateOrderItemStatusRequest struct {
	ItemStatus string `json:"item_status"`
}

type callerInfo struct {
	UserID string
	Roles  []string
}

func (h *OrderHandler) verifyCaller(r *http.Request) (*callerInfo, error) {
	authHeader := strings.TrimSpace(r.Header.Get("Authorization"))
	if authHeader == "" {
		return nil, nil
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

func hasStaffRole(roles []string) bool {
	for _, role := range roles {
		switch role {
		case "ADMIN", "MANAGER", "CHEF", "WAITER":
			return true
		}
	}
	return false
}

// canMarkItemStatus enforces role-based transition rules:
// COOKING/READY → CHEF, ADMIN, MANAGER
// SERVED        → WAITER, ADMIN, MANAGER
func canMarkItemStatus(roles []string, targetStatus string) bool {
	switch targetStatus {
	case "COOKING", "READY":
		for _, r := range roles {
			if r == "ADMIN" || r == "MANAGER" || r == "CHEF" {
				return true
			}
		}
	case "SERVED":
		for _, r := range roles {
			if r == "ADMIN" || r == "MANAGER" || r == "WAITER" {
				return true
			}
		}
	}
	return false
}

func (h *OrderHandler) checkOrderAccess(r *http.Request, orderUserID string) (int, error) {
	if orderUserID == "" {
		return 0, nil
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	if caller == nil {
		return http.StatusUnauthorized, errUnauthorized("authentication required to access this order")
	}
	if caller.UserID != orderUserID && !hasStaffRole(caller.Roles) {
		return http.StatusForbidden, fmt.Errorf("access denied")
	}
	return 0, nil
}

func (h *OrderHandler) checkUserIDAccess(r *http.Request, targetUserID string) (int, error) {
	caller, err := h.verifyCaller(r)
	if err != nil {
		return http.StatusUnauthorized, err
	}
	if caller == nil {
		return http.StatusUnauthorized, errUnauthorized("authentication required")
	}
	if caller.UserID != targetUserID && !hasStaffRole(caller.Roles) {
		return http.StatusForbidden, fmt.Errorf("access denied")
	}
	return 0, nil
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}
	userID := ""
	if caller != nil {
		userID = caller.UserID
	}
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.orderClient.CreateOrder(r.Context(), &orderpb.CreateOrderRequest{
		TableId:   req.TableID,
		UserId:    userID,
		Name:      req.Name,
		Phone:     req.Phone,
		Notes:     req.Notes,
		Time:      req.Time,
		EndTime:   req.EndTime,
		Date:      req.Date,
		PartySize: req.PartySize,
		Status:    req.Status,
		Items:     toOrderItemProto(req.Items),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	if userID != "" {
		if status, err := h.checkUserIDAccess(r, userID); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	resp, err := h.orderClient.ListOrders(r.Context(), &orderpb.ListOrdersRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 20),
		Status:   strings.TrimSpace(r.URL.Query().Get("status")),
		Keyword:  strings.TrimSpace(r.URL.Query().Get("keyword")),
		UserId:   userID,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	orderID := extractOrderIDFromPath(r.URL.Path)
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}
	resp, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if resp.Order != nil {
		if status, err := h.checkOrderAccess(r, resp.Order.UserId); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}
	orderID := extractOrderIDFromPath(r.URL.Path)
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}
	existing, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if existing.Order != nil {
		if status, err := h.checkOrderAccess(r, existing.Order.UserId); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	var req updateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.orderClient.UpdateOrder(r.Context(), &orderpb.UpdateOrderRequest{
		OrderId:   orderID,
		TableId:   req.TableID,
		Name:      req.Name,
		Phone:     req.Phone,
		Notes:     req.Notes,
		Time:      req.Time,
		EndTime:   req.EndTime,
		Date:      req.Date,
		PartySize: req.PartySize,
		Status:    req.Status,
		Items:     toOrderItemProto(req.Items),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) DeleteOrder(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	orderID := extractOrderIDFromPath(r.URL.Path)
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}
	existing, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if existing.Order != nil {
		if status, err := h.checkOrderAccess(r, existing.Order.UserId); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	resp, err := h.orderClient.DeleteOrder(r.Context(), &orderpb.DeleteOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	orderID := extractOrderIDForActionPath(r.URL.Path, "cancel")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}
	existing, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if existing.Order != nil {
		if status, err := h.checkOrderAccess(r, existing.Order.UserId); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	var req cancelOrderRequest
	if r.Body != nil {
		_ = json.NewDecoder(r.Body).Decode(&req)
	}
	resp, err := h.orderClient.CancelOrder(r.Context(), &orderpb.CancelOrderRequest{OrderId: orderID, Reason: req.Reason})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	if !allowPatch(w, r) {
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}
	if caller == nil || !hasStaffRole(caller.Roles) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "staff access required"})
		return
	}
	orderID := extractOrderIDForActionPath(r.URL.Path, "status")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}
	var req updateOrderStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.orderClient.UpdateOrderStatus(r.Context(), &orderpb.UpdateOrderStatusRequest{OrderId: orderID, Status: req.Status})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) AddOrderItem(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	orderID := extractOrderIDForActionPath(r.URL.Path, "items")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}
	existing, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if existing.Order != nil {
		if status, err := h.checkOrderAccess(r, existing.Order.UserId); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	var req addOrderItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	resp, err := h.orderClient.AddOrderItem(r.Context(), &orderpb.AddOrderItemRequest{
		OrderId: orderID,
		Item:    toSingleOrderItemProto(req.Item),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) RemoveOrderItem(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	orderID, itemID := extractOrderAndItemIDForItemsPath(r.URL.Path)
	if orderID == "" || itemID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id and item_id are required"})
		return
	}
	existing, err := h.orderClient.GetOrder(r.Context(), &orderpb.GetOrderRequest{OrderId: orderID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	if existing.Order != nil {
		if status, err := h.checkOrderAccess(r, existing.Order.UserId); err != nil {
			writeJSON(w, status, map[string]any{"error": err.Error()})
			return
		}
	}
	resp, err := h.orderClient.RemoveOrderItem(r.Context(), &orderpb.RemoveOrderItemRequest{OrderId: orderID, ItemId: itemID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrderItemStatus(w http.ResponseWriter, r *http.Request) {
	if !allowPatch(w, r) {
		return
	}
	orderID, itemID := extractOrderItemStatusPath(r.URL.Path)
	if orderID == "" || itemID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id and item_id are required"})
		return
	}
	var req updateOrderItemStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}
	if req.ItemStatus == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "item_status is required"})
		return
	}
	caller, err := h.verifyCaller(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
		return
	}
	if caller == nil || !canMarkItemStatus(caller.Roles, req.ItemStatus) {
		writeJSON(w, http.StatusForbidden, map[string]any{"error": "insufficient role for this status transition"})
		return
	}
	resp, err := h.orderClient.UpdateOrderItemStatus(r.Context(), &orderpb.UpdateOrderItemStatusRequest{
		OrderId:    orderID,
		ItemId:     itemID,
		ItemStatus: req.ItemStatus,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, resp)
}

func toOrderItemProto(items []orderItemRequest) []*orderpb.OrderItemRequest {
	result := make([]*orderpb.OrderItemRequest, 0, len(items))
	for _, item := range items {
		result = append(result, toSingleOrderItemProto(item))
	}
	return result
}

func toSingleOrderItemProto(item orderItemRequest) *orderpb.OrderItemRequest {
	return &orderpb.OrderItemRequest{
		ItemId:   item.ItemID,
		Name:     item.Name,
		Price:    item.Price,
		Quantity: item.Quantity,
	}
}

func extractOrderIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/orders/"), " ")
	if trimmed == "" {
		return ""
	}
	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}
	return strings.TrimSpace(segments[0])
}

func extractOrderIDForActionPath(path, action string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/orders/"), " ")
	if trimmed == "" {
		return ""
	}
	segments := strings.Split(trimmed, "/")
	if len(segments) != 2 || strings.TrimSpace(segments[1]) != action {
		return ""
	}
	return strings.TrimSpace(segments[0])
}

func extractOrderAndItemIDForItemsPath(path string) (string, string) {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/orders/"), " ")
	if trimmed == "" {
		return "", ""
	}
	segments := strings.Split(trimmed, "/")
	if len(segments) != 3 || strings.TrimSpace(segments[1]) != "items" {
		return "", ""
	}
	return strings.TrimSpace(segments[0]), strings.TrimSpace(segments[2])
}

func extractOrderItemStatusPath(path string) (string, string) {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/orders/"), " ")
	if trimmed == "" {
		return "", ""
	}
	segments := strings.Split(trimmed, "/")
	if len(segments) != 4 ||
		strings.TrimSpace(segments[1]) != "items" ||
		strings.TrimSpace(segments[3]) != "status" {
		return "", ""
	}
	return strings.TrimSpace(segments[0]), strings.TrimSpace(segments[2])
}
