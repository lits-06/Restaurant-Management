package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	authpb "restaurant-management/proto/auth"
	orderpb "restaurant-management/proto/order"
)

// OrderHandler exposes HTTP endpoints backed by order-service gRPC methods.
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
	Name      string             `json:"name"`
	Phone     string             `json:"phone"`
	Time      string             `json:"time"`
	Date      string             `json:"date"`
	PartySize int32              `json:"party_size"`
	Status    string             `json:"status"`
	Items     []orderItemRequest `json:"items"`
}

type updateOrderRequest struct {
	Name      string             `json:"name"`
	Phone     string             `json:"phone"`
	Time      string             `json:"time"`
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

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	items := toOrderItemProto(req.Items)
	resp, err := h.orderClient.CreateOrder(r.Context(), &orderpb.CreateOrderRequest{
		Name:      req.Name,
		Phone:     req.Phone,
		Time:      req.Time,
		Date:      req.Date,
		PartySize: req.PartySize,
		Status:    req.Status,
		Items:     items,
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

	resp, err := h.orderClient.ListOrders(r.Context(), &orderpb.ListOrdersRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 20),
		Status:   strings.TrimSpace(r.URL.Query().Get("status")),
		Keyword:  strings.TrimSpace(r.URL.Query().Get("keyword")),
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

	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	orderID := extractOrderIDFromPath(r.URL.Path)
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
	}

	var req updateOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.orderClient.UpdateOrder(r.Context(), &orderpb.UpdateOrderRequest{
		OrderId:   orderID,
		Name:      req.Name,
		Phone:     req.Phone,
		Time:      req.Time,
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

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	orderID := extractOrderIDFromPath(r.URL.Path)
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
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

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	orderID := extractOrderIDForActionPath(r.URL.Path, "cancel")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
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

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

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

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	orderID := extractOrderIDForActionPath(r.URL.Path, "items")
	if orderID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id is required"})
		return
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

	// if _, err := h.verifyAuthorization(r); err != nil {
	// 	writeJSON(w, http.StatusUnauthorized, map[string]any{"error": err.Error()})
	// 	return
	// }

	orderID, itemID := extractOrderAndItemIDForItemsPath(r.URL.Path)
	if orderID == "" || itemID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "order_id and item_id are required"})
		return
	}

	resp, err := h.orderClient.RemoveOrderItem(r.Context(), &orderpb.RemoveOrderItemRequest{OrderId: orderID, ItemId: itemID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *OrderHandler) verifyAuthorization(r *http.Request) (*authpb.VerifyTokenResponse, error) {
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
