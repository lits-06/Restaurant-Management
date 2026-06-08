package usecase

import (
	"context"
	"fmt"

	"restaurant-management/services/menu-service/internal/domain"
	"restaurant-management/services/menu-service/internal/repository"
)

// MenuUseCase handles menu business logic
type MenuUseCase struct {
	menuItemRepo repository.MenuItemRepository
	categoryRepo repository.CategoryRepository
}

// NewMenuUseCase creates a new menu use case
func NewMenuUseCase(menuItemRepo repository.MenuItemRepository, categoryRepo repository.CategoryRepository) *MenuUseCase {
	return &MenuUseCase{
		menuItemRepo: menuItemRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateMenuItem creates a new menu item
func (uc *MenuUseCase) CreateMenuItem(ctx context.Context, name, description string, price float64, category, imageURL string) (*domain.MenuItem, error) {
	// Verify category exists
	c, err := uc.categoryRepo.GetByName(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("category not found: %w", err)
	}

	// Create menu item
	item, err := domain.NewMenuItem(name, description, price, c.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to create menu item: %w", err)
	}

	item.ImageURL = imageURL

	// Validate again after setting optional fields
	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("menu item validation failed: %w", err)
	}

	// Save to repository
	if err := uc.menuItemRepo.Create(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to save menu item: %w", err)
	}

	return item, nil
}

// GetMenuItem retrieves a menu item by ID
func (uc *MenuUseCase) GetMenuItem(ctx context.Context, itemID string) (*domain.MenuItem, error) {
	item, err := uc.menuItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}
	return item, nil
}

// GetMenuItemByName retrieves a menu item by name
func (uc *MenuUseCase) GetMenuItemByName(ctx context.Context, name string) (*domain.MenuItem, error) {
	item, err := uc.menuItemRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item by name: %w", err)
	}
	return item, nil
}

// UpdateMenuItem updates a menu item
func (uc *MenuUseCase) UpdateMenuItem(ctx context.Context, itemID, name, description string, price float64, categoryID, imageURL string) (*domain.MenuItem, error) {
	// Get existing item
	item, err := uc.menuItemRepo.GetByID(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get menu item: %w", err)
	}

	// Verify category exists if changed
	if categoryID != "" && categoryID != item.CategoryID {
		_, err := uc.categoryRepo.GetByID(ctx, categoryID)
		if err != nil {
			return nil, fmt.Errorf("category not found: %w", err)
		}
		item.CategoryID = categoryID
	}

	// Update fields
	if name != "" {
		item.Name = name
	}
	if description != "" {
		item.Description = description
	}
	if price >= 0 {
		if err := item.UpdatePrice(price); err != nil {
			return nil, fmt.Errorf("failed to update price: %w", err)
		}
	}
	if imageURL != "" {
		item.ImageURL = imageURL
	}

	// Validate
	if err := item.Validate(); err != nil {
		return nil, fmt.Errorf("menu item validation failed: %w", err)
	}

	// Save
	if err := uc.menuItemRepo.Update(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to update menu item: %w", err)
	}

	return item, nil
}

// DeleteMenuItem deletes a menu item
func (uc *MenuUseCase) DeleteMenuItem(ctx context.Context, itemID string) error {
	if err := uc.menuItemRepo.Delete(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete menu item: %w", err)
	}
	return nil
}

// ListMenuItems retrieves menu items with pagination and filters
func (uc *MenuUseCase) ListMenuItems(ctx context.Context, page, pageSize int, categoryID string, keyword string) ([]*domain.MenuItem, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	items, total, err := uc.menuItemRepo.List(ctx, page, pageSize, categoryID, keyword)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list menu items: %w", err)
	}

	return items, total, nil
}

// CreateCategory creates a new category
func (uc *MenuUseCase) CreateCategory(ctx context.Context, name string) (*domain.Category, error) {
	category, err := domain.NewCategory(name)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	if err := uc.categoryRepo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	return category, nil
}

// GetCategory retrieves a category by ID
func (uc *MenuUseCase) GetCategory(ctx context.Context, categoryID string) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}
	return category, nil
}

// GetCategoryByName retrieves a category by name
func (uc *MenuUseCase) GetCategoryByName(ctx context.Context, name string) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get category by name: %w", err)
	}
	return category, nil
}

// UpdateCategory updates a category
func (uc *MenuUseCase) UpdateCategory(ctx context.Context, categoryID, name string) (*domain.Category, error) {
	category, err := uc.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	if name != "" {
		category.Name = name
	}

	if err := category.Validate(); err != nil {
		return nil, fmt.Errorf("category validation failed: %w", err)
	}

	if err := uc.categoryRepo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category
func (uc *MenuUseCase) DeleteCategory(ctx context.Context, categoryID string) error {
	if err := uc.categoryRepo.Delete(ctx, categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// // ListCategories retrieves categories with pagination
// func (uc *MenuUseCase) ListCategories(ctx context.Context, page, pageSize int) ([]*domain.Category, int, error) {
// 	if page < 1 {
// 		page = 1
// 	}
// 	if pageSize < 1 || pageSize > 100 {
// 		pageSize = 20
// 	}

// 	categories, total, err := uc.categoryRepo.List(ctx, page, pageSize)
// 	if err != nil {
// 		return nil, 0, fmt.Errorf("failed to list categories: %w", err)
// 	}

// 	return categories, total, nil
// }

// GetAllCategories retrieves all categories
func (uc *MenuUseCase) GetAllCategories(ctx context.Context) ([]*domain.Category, error) {
	categories, err := uc.categoryRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	return categories, nil
}
