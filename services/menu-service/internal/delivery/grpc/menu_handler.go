package grpc

import (
	"context"
	"errors"

	"restaurant-management/proto/menu"
	"restaurant-management/services/menu-service/internal/domain"
	"restaurant-management/services/menu-service/internal/usecase"
	"restaurant-management/shared/pkg/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MenuHandler handles gRPC requests for menu service
type MenuHandler struct {
	menu.UnimplementedMenuServiceServer
	menuUseCase *usecase.MenuUseCase
}

// NewMenuHandler creates a new menu handler
func NewMenuHandler(menuUseCase *usecase.MenuUseCase) *MenuHandler {
	return &MenuHandler{
		menuUseCase: menuUseCase,
	}
}

// CreateMenuItem creates a new menu item
func (h *MenuHandler) CreateMenuItem(ctx context.Context, req *menu.CreateMenuItemRequest) (*menu.CreateMenuItemResponse, error) {
	if req.Name == "" {
		return &menu.CreateMenuItemResponse{
			Success: false,
			Message: "name is required",
		}, nil
	}

	if req.Price < 0 {
		return &menu.CreateMenuItemResponse{
			Success: false,
			Message: "price must be non-negative",
		}, nil
	}

	if req.Category == "" {
		return &menu.CreateMenuItemResponse{
			Success: false,
			Message: "category is required",
		}, nil
	}

	// Create menu item
	item, err := h.menuUseCase.CreateMenuItem(
		ctx,
		req.Name,
		req.Description,
		req.Price,
		req.Category,
		req.ImageUrl,
	)

	if err != nil {
		logger.Error("Failed to create menu item", zap.Error(err))
		return &menu.CreateMenuItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.CreateMenuItemResponse{
		Item:    convertMenuItemToProto(item),
		Success: true,
		Message: "menu item created successfully",
	}, nil
}

// GetMenuItem retrieves a menu item by ID
func (h *MenuHandler) GetMenuItem(ctx context.Context, req *menu.GetMenuItemRequest) (*menu.GetMenuItemResponse, error) {
	logger.Info("GetMenuItem request", zap.String("item_id", req.ItemId))

	if req.ItemId == "" {
		return &menu.GetMenuItemResponse{
			Success: false,
			Message: "item_id is required",
		}, nil
	}

	item, err := h.menuUseCase.GetMenuItem(ctx, req.ItemId)
	if err != nil {
		if errors.Is(err, domain.ErrMenuItemNotFound) {
			return &menu.GetMenuItemResponse{
				Success: false,
				Message: "menu item not found",
			}, nil
		}
		logger.Error("Failed to get menu item", zap.Error(err))
		return &menu.GetMenuItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.GetMenuItemResponse{
		Item:    convertMenuItemToProto(item),
		Success: true,
		Message: "menu item retrieved successfully",
	}, nil
}

// UpdateMenuItem updates a menu item
func (h *MenuHandler) UpdateMenuItem(ctx context.Context, req *menu.UpdateMenuItemRequest) (*menu.UpdateMenuItemResponse, error) {
	logger.Info("UpdateMenuItem request", zap.String("item_id", req.ItemId))

	if req.ItemId == "" {
		return &menu.UpdateMenuItemResponse{
			Success: false,
			Message: "item_id is required",
		}, nil
	}

	item, err := h.menuUseCase.UpdateMenuItem(
		ctx,
		req.ItemId,
		req.Name,
		req.Description,
		req.Price,
		req.CategoryId,
		req.ImageUrl,
	)

	if err != nil {
		if errors.Is(err, domain.ErrMenuItemNotFound) {
			return &menu.UpdateMenuItemResponse{
				Success: false,
				Message: "menu item not found",
			}, nil
		}
		logger.Error("Failed to update menu item", zap.Error(err))
		return &menu.UpdateMenuItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.UpdateMenuItemResponse{
		Item:    convertMenuItemToProto(item),
		Success: true,
		Message: "menu item updated successfully",
	}, nil
}

// DeleteMenuItem deletes a menu item
func (h *MenuHandler) DeleteMenuItem(ctx context.Context, req *menu.DeleteMenuItemRequest) (*menu.DeleteMenuItemResponse, error) {
	logger.Info("DeleteMenuItem request", zap.String("item_id", req.ItemId))

	if req.ItemId == "" {
		return &menu.DeleteMenuItemResponse{
			Success: false,
			Message: "item_id is required",
		}, nil
	}

	err := h.menuUseCase.DeleteMenuItem(ctx, req.ItemId)
	if err != nil {
		if errors.Is(err, domain.ErrMenuItemNotFound) {
			return &menu.DeleteMenuItemResponse{
				Success: false,
				Message: "menu item not found",
			}, nil
		}
		logger.Error("Failed to delete menu item", zap.Error(err))
		return &menu.DeleteMenuItemResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.DeleteMenuItemResponse{
		Success: true,
		Message: "menu item deleted successfully",
	}, nil
}

// ListMenuItems retrieves menu items with pagination and filters
func (h *MenuHandler) ListMenuItems(ctx context.Context, req *menu.ListMenuItemsRequest) (*menu.ListMenuItemsResponse, error) {
	logger.Info("ListMenuItems request", zap.Int32("page", req.Page), zap.Int32("page_size", req.PageSize))

	page := int(req.Page)
	pageSize := int(req.PageSize)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	items, total, err := h.menuUseCase.ListMenuItems(ctx, page, pageSize, req.CategoryId, req.Keyword)
	if err != nil {
		logger.Error("Failed to list menu items", zap.Error(err))
		return &menu.ListMenuItemsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	protoItems := make([]*menu.MenuItem, len(items))
	for i, item := range items {
		protoItems[i] = convertMenuItemToProto(item)
	}

	return &menu.ListMenuItemsResponse{
		Items:    protoItems,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
		Success:  true,
		Message:  "menu items retrieved successfully",
	}, nil
}

// CreateCategory creates a new category
func (h *MenuHandler) CreateCategory(ctx context.Context, req *menu.CreateCategoryRequest) (*menu.CreateCategoryResponse, error) {
	logger.Info("CreateCategory request", zap.String("name", req.Name))

	if req.Name == "" {
		return &menu.CreateCategoryResponse{
			Success: false,
			Message: "name is required",
		}, nil
	}

	category, err := h.menuUseCase.CreateCategory(ctx, req.Name)
	if err != nil {
		logger.Error("Failed to create category", zap.Error(err))
		return &menu.CreateCategoryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.CreateCategoryResponse{
		Category: convertCategoryToProto(category),
		Success:  true,
		Message:  "category created successfully",
	}, nil
}

// GetCategory retrieves a category by ID
func (h *MenuHandler) GetCategory(ctx context.Context, req *menu.GetCategoryRequest) (*menu.GetCategoryResponse, error) {
	logger.Info("GetCategory request", zap.String("category_id", req.CategoryId))

	if req.CategoryId == "" {
		return &menu.GetCategoryResponse{
			Success: false,
			Message: "category_id is required",
		}, nil
	}

	category, err := h.menuUseCase.GetCategory(ctx, req.CategoryId)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return &menu.GetCategoryResponse{
				Success: false,
				Message: "category not found",
			}, nil
		}
		logger.Error("Failed to get category", zap.Error(err))
		return &menu.GetCategoryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.GetCategoryResponse{
		Category: convertCategoryToProto(category),
		Success:  true,
		Message:  "category retrieved successfully",
	}, nil
}

// UpdateCategory updates a category
func (h *MenuHandler) UpdateCategory(ctx context.Context, req *menu.UpdateCategoryRequest) (*menu.UpdateCategoryResponse, error) {
	logger.Info("UpdateCategory request", zap.String("category_id", req.CategoryId))

	if req.CategoryId == "" {
		return &menu.UpdateCategoryResponse{
			Success: false,
			Message: "category_id is required",
		}, nil
	}

	category, err := h.menuUseCase.UpdateCategory(ctx, req.CategoryId, req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return &menu.UpdateCategoryResponse{
				Success: false,
				Message: "category not found",
			}, nil
		}
		logger.Error("Failed to update category", zap.Error(err))
		return &menu.UpdateCategoryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.UpdateCategoryResponse{
		Category: convertCategoryToProto(category),
		Success:  true,
		Message:  "category updated successfully",
	}, nil
}

// DeleteCategory deletes a category
func (h *MenuHandler) DeleteCategory(ctx context.Context, req *menu.DeleteCategoryRequest) (*menu.DeleteCategoryResponse, error) {
	logger.Info("DeleteCategory request", zap.String("category_id", req.CategoryId))

	if req.CategoryId == "" {
		return &menu.DeleteCategoryResponse{
			Success: false,
			Message: "category_id is required",
		}, nil
	}

	err := h.menuUseCase.DeleteCategory(ctx, req.CategoryId)
	if err != nil {
		if errors.Is(err, domain.ErrCategoryNotFound) {
			return &menu.DeleteCategoryResponse{
				Success: false,
				Message: "category not found",
			}, nil
		}
		if errors.Is(err, domain.ErrCategoryHasMenuItems) {
			return &menu.DeleteCategoryResponse{
				Success: false,
				Message: "category has menu items and cannot be deleted",
			}, nil
		}
		logger.Error("Failed to delete category", zap.Error(err))
		return &menu.DeleteCategoryResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &menu.DeleteCategoryResponse{
		Success: true,
		Message: "category deleted successfully",
	}, nil
}

// ListCategories retrieves categories with pagination
func (h *MenuHandler) ListCategories(ctx context.Context, req *menu.ListCategoriesRequest) (*menu.ListCategoriesResponse, error) {
	logger.Info("ListCategories request", zap.Int32("page", req.Page), zap.Int32("page_size", req.PageSize))

	categories, err := h.menuUseCase.GetAllCategories(ctx)
	if err != nil {
		logger.Error("Failed to get all categories", zap.Error(err))
		return &menu.ListCategoriesResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	protoCategories := make([]*menu.Category, len(categories))
	for i, category := range categories {
		protoCategories[i] = convertCategoryToProto(category)
	}

	return &menu.ListCategoriesResponse{
		Categories: protoCategories,
		Total:      int32(len(categories)),
		Success:    true,
		Message:    "categories retrieved successfully",
	}, nil
}

// Helper functions

func convertMenuItemToProto(item *domain.MenuItem) *menu.MenuItem {
	if item == nil {
		return nil
	}

	return &menu.MenuItem{
		ItemId:      item.ItemID,
		Name:        item.Name,
		Description: item.Description,
		Price:       item.Price,
		Category:    item.CategoryID,
		ImageUrl:    item.ImageURL,
	}
}

func convertCategoryToProto(category *domain.Category) *menu.Category {
	if category == nil {
		return nil
	}

	return &menu.Category{
		CategoryId: category.CategoryID,
		Name:       category.Name,
	}
}

func mapErrorToGRPCStatus(err error) error {
	switch {
	case errors.Is(err, domain.ErrMenuItemNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrMenuItemAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrCategoryNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, domain.ErrCategoryAlreadyExists):
		return status.Error(codes.AlreadyExists, err.Error())
	case errors.Is(err, domain.ErrCategoryHasMenuItems):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
}
