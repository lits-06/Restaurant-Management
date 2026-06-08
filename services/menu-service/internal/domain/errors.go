package domain

import "errors"

// MenuItem errors
var (
	ErrMenuItemNotFound               = errors.New("menu item not found")
	ErrMenuItemNameRequired           = errors.New("menu item name is required")
	ErrMenuItemNameInvalid            = errors.New("menu item name must be between 2 and 100 characters")
	ErrMenuItemPriceInvalid           = errors.New("menu item price must be non-negative")
	ErrMenuItemPriceTooHigh           = errors.New("menu item price exceeds maximum allowed")
	ErrMenuItemCategoryRequired       = errors.New("menu item category is required")
	ErrMenuItemPreparationTimeInvalid = errors.New("preparation time must be non-negative")
	ErrMenuItemPreparationTimeTooLong = errors.New("preparation time exceeds maximum (480 minutes)")
	ErrMenuItemDescriptionTooLong     = errors.New("menu item description is too long (max 1000 characters)")
	ErrMenuItemImageURLTooLong        = errors.New("menu item image URL is too long (max 500 characters)")
	ErrMenuItemAlreadyExists          = errors.New("menu item with this name already exists")
)

// Category errors
var (
	ErrCategoryNotFound             = errors.New("category not found")
	ErrCategoryNameRequired         = errors.New("category name is required")
	ErrCategoryNameInvalid          = errors.New("category name must be between 2 and 100 characters")
	ErrCategoryDisplayOrderInvalid  = errors.New("display order must be non-negative")
	ErrCategoryDescriptionTooLong   = errors.New("category description is too long (max 500 characters)")
	ErrCategoryAlreadyExists        = errors.New("category with this name already exists")
	ErrCategoryHasMenuItems         = errors.New("category has menu items and cannot be deleted")
)
