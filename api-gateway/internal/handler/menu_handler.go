package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"restaurant-management/api-gateway/internal/grpcclient"
	menupb "restaurant-management/proto/menu"
)

// MenuHandler exposes HTTP endpoints backed by menu-service gRPC methods.
type MenuHandler struct {
	menuClient *grpcclient.MenuClient
	authClient *grpcclient.AuthClient
}

func NewMenuHandler(menuClient *grpcclient.MenuClient, authClient *grpcclient.AuthClient) *MenuHandler {
	return &MenuHandler{
		menuClient: menuClient,
		authClient: authClient,
	}
}

type createMenuItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Category    string  `json:"category"`
	CategoryID  string  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

type updateMenuItemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CategoryID  string  `json:"category_id"`
	ImageURL    string  `json:"image_url"`
}

type createCategoryRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayOrder int32  `json:"display_order"`
}

type updateCategoryRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayOrder int32  `json:"display_order"`
}

// CreateMenuItem handles POST /menu/items.
func (h *MenuHandler) CreateMenuItem(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	if requireAdminOrManager(w, r, h.authClient) == nil {
		return
	}

	var req createMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.menuClient.CreateMenuItem(r.Context(), &menupb.CreateMenuItemRequest{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    firstNonEmpty(req.CategoryID, req.Category),
		ImageUrl:    req.ImageURL,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListMenuItems handles GET /menu/items.
func (h *MenuHandler) ListMenuItems(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	resp, err := h.menuClient.ListMenuItems(r.Context(), &menupb.ListMenuItemsRequest{
		Page:       parseInt32Query(r, "page", 1),
		PageSize:   parseInt32Query(r, "page_size", 20),
		CategoryId: strings.TrimSpace(r.URL.Query().Get("category_id")),
		Keyword:    strings.TrimSpace(r.URL.Query().Get("keyword")),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetMenuItem handles GET /menu/items/{item_id}.
func (h *MenuHandler) GetMenuItem(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	itemID := extractMenuItemIDFromPath(r.URL.Path)
	if itemID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "item_id is required"})
		return
	}

	resp, err := h.menuClient.GetMenuItem(r.Context(), &menupb.GetMenuItemRequest{ItemId: itemID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateMenuItem handles PUT /menu/items/{item_id}.
func (h *MenuHandler) UpdateMenuItem(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}
	if requireAdminOrManager(w, r, h.authClient) == nil {
		return
	}

	itemID := extractMenuItemIDFromPath(r.URL.Path)
	if itemID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "item_id is required"})
		return
	}

	var req updateMenuItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.menuClient.UpdateMenuItem(r.Context(), &menupb.UpdateMenuItemRequest{
		ItemId:      itemID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		CategoryId:  req.CategoryID,
		ImageUrl:    req.ImageURL,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// DeleteMenuItem handles DELETE /menu/items/{item_id}.
func (h *MenuHandler) DeleteMenuItem(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	if requireAdminOrManager(w, r, h.authClient) == nil {
		return
	}

	itemID := extractMenuItemIDFromPath(r.URL.Path)
	if itemID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "item_id is required"})
		return
	}

	resp, err := h.menuClient.DeleteMenuItem(r.Context(), &menupb.DeleteMenuItemRequest{ItemId: itemID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// CreateCategory handles POST /menu/categories.
func (h *MenuHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	if !allowPost(w, r) {
		return
	}
	if requireAdminOrManager(w, r, h.authClient) == nil {
		return
	}

	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.menuClient.CreateCategory(r.Context(), &menupb.CreateCategoryRequest{
		Name:         req.Name,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// ListCategories handles GET /menu/categories.
func (h *MenuHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	resp, err := h.menuClient.ListCategories(r.Context(), &menupb.ListCategoriesRequest{
		Page:     parseInt32Query(r, "page", 1),
		PageSize: parseInt32Query(r, "page_size", 20),
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// GetAllCategories handles GET /menu/categories.
func (h *MenuHandler) GetAllCategories(w http.ResponseWriter, r *http.Request) {
	h.ListCategories(w, r)
}

// GetCategory handles GET /menu/categories/{category_id}.
func (h *MenuHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}

	categoryID := extractMenuCategoryIDFromPath(r.URL.Path)
	if categoryID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "category_id is required"})
		return
	}

	resp, err := h.menuClient.GetCategory(r.Context(), &menupb.GetCategoryRequest{CategoryId: categoryID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// UpdateCategory handles PUT /menu/categories/{category_id}.
func (h *MenuHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	if !allowPut(w, r) {
		return
	}
	if requireAdminOrManager(w, r, h.authClient) == nil {
		return
	}

	categoryID := extractMenuCategoryIDFromPath(r.URL.Path)
	if categoryID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "category_id is required"})
		return
	}

	var req updateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid request body"})
		return
	}

	resp, err := h.menuClient.UpdateCategory(r.Context(), &menupb.UpdateCategoryRequest{
		CategoryId:   categoryID,
		Name:         req.Name,
		Description:  req.Description,
		DisplayOrder: req.DisplayOrder,
	})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

// DeleteCategory handles DELETE /menu/categories/{category_id}.
func (h *MenuHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	if !allowDelete(w, r) {
		return
	}
	if requireAdminOrManager(w, r, h.authClient) == nil {
		return
	}

	categoryID := extractMenuCategoryIDFromPath(r.URL.Path)
	if categoryID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "category_id is required"})
		return
	}

	resp, err := h.menuClient.DeleteCategory(r.Context(), &menupb.DeleteCategoryRequest{CategoryId: categoryID})
	if err != nil {
		writeJSON(w, http.StatusBadGateway, map[string]any{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

type unauthorizedError string

func (e unauthorizedError) Error() string {
	return string(e)
}

func errUnauthorized(message string) error {
	return unauthorizedError(message)
}

func extractBearerToken(header string) string {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(header, bearerPrefix) {
		return ""
	}
	return strings.TrimSpace(header[len(bearerPrefix):])
}

func parseInt32Query(r *http.Request, key string, defaultValue int32) int32 {
	raw := strings.TrimSpace(r.URL.Query().Get(key))
	if raw == "" {
		return defaultValue
	}

	parsed, err := strconv.ParseInt(raw, 10, 32)
	if err != nil {
		return defaultValue
	}

	return int32(parsed)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func extractMenuItemIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/menu/items/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}
	return strings.TrimSpace(segments[0])
}

func extractMenuCategoryIDFromPath(path string) string {
	trimmed := strings.Trim(strings.TrimPrefix(path, "/menu/categories/"), " ")
	if trimmed == "" {
		return ""
	}

	segments := strings.Split(trimmed, "/")
	if len(segments) == 0 {
		return ""
	}
	return strings.TrimSpace(segments[0])
}
