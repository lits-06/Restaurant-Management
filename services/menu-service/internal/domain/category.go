package domain

// Category represents a menu category
type Category struct {
	CategoryID string
	Name       string
}

// Validate validates the category
func (c *Category) Validate() error {
	if c.Name == "" {
		return ErrCategoryNameRequired
	}

	if len(c.Name) < 2 || len(c.Name) > 100 {
		return ErrCategoryNameInvalid
	}

	return nil
}

// NewCategory creates a new category
func NewCategory(name string) (*Category, error) {
	category := &Category{
		Name: name,
	}

	if err := category.Validate(); err != nil {
		return nil, err
	}

	return category, nil
}
