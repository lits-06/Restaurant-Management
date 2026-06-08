package domain

// MenuItem represents a menu item in the restaurant
type MenuItem struct {
	ItemID      string
	Name        string
	Description string
	Price       float64
	CategoryID  string
	ImageURL    string
}

// Validate validates the menu item
func (m *MenuItem) Validate() error {
	if m.Name == "" {
		return ErrMenuItemNameRequired
	}

	if len(m.Name) < 2 || len(m.Name) > 100 {
		return ErrMenuItemNameInvalid
	}

	if m.Price < 0 {
		return ErrMenuItemPriceInvalid
	}

	if m.Price > 100000000 { // 100 million max
		return ErrMenuItemPriceTooHigh
	}

	if m.CategoryID == "" {
		return ErrMenuItemCategoryRequired
	}

	if len(m.Description) > 1000 {
		return ErrMenuItemDescriptionTooLong
	}

	if len(m.ImageURL) > 500 {
		return ErrMenuItemImageURLTooLong
	}

	return nil
}

// UpdatePrice updates the price of the menu item
func (m *MenuItem) UpdatePrice(newPrice float64) error {
	if newPrice < 0 {
		return ErrMenuItemPriceInvalid
	}
	if newPrice > 100000000 {
		return ErrMenuItemPriceTooHigh
	}
	m.Price = newPrice
	return nil
}

// NewMenuItem creates a new menu item
func NewMenuItem(name, description string, price float64, categoryID string) (*MenuItem, error) {
	item := &MenuItem{
		Name:        name,
		Description: description,
		Price:       price,
		CategoryID:  categoryID,
	}

	if err := item.Validate(); err != nil {
		return nil, err
	}

	return item, nil
}
