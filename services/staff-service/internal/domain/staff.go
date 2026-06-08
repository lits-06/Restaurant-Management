package domain

// Staff represents a staff member in the restaurant.
type Staff struct {
	StaffID string
	Name    string
	Role    string
	Contact string
	Avatar  string
}

// NewStaff creates a new staff member with validation.
func NewStaff(name, role, contact, avatar string) (*Staff, error) {
	staff := &Staff{
		Name:    name,
		Role:    role,
		Contact: contact,
		Avatar:  avatar,
	}

	if err := staff.Validate(); err != nil {
		return nil, err
	}

	return staff, nil
}

// Validate validates the staff fields.
func (s *Staff) Validate() error {
	if s.Name == "" {
		return ErrInvalidStaffName
	}
	if len(s.Name) < 2 {
		return ErrStaffNameTooShort
	}
	if len(s.Name) > 120 {
		return ErrStaffNameTooLong
	}
	if s.Role == "" {
		return ErrInvalidStaffRole
	}
	if s.Contact == "" {
		return ErrInvalidStaffContact
	}
	if len(s.Contact) > 100 {
		return ErrStaffContactTooLong
	}
	return nil
}

// Update updates staff information.
func (s *Staff) Update(name, role, contact, avatar string) error {
	if name != "" {
		if len(name) < 2 {
			return ErrStaffNameTooShort
		}
		if len(name) > 120 {
			return ErrStaffNameTooLong
		}
		s.Name = name
	}

	if role != "" {
		s.Role = role
	}

	if contact != "" {
		contact = contact
		if len(contact) > 100 {
			return ErrStaffContactTooLong
		}
		s.Contact = contact
	}

	if avatar != "" {
		avatar = avatar
		if len(avatar) > 500 {
			return ErrStaffAvatarTooLong
		}
		s.Avatar = avatar
	}

	if s.Name == "" {
		return ErrInvalidStaffName
	}
	if s.Role == "" {
		return ErrInvalidStaffRole
	}
	if s.Contact == "" {
		return ErrInvalidStaffContact
	}
	return nil
}
