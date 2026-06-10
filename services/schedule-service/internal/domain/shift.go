package domain

import (
	"strings"
	"time"
)

type Shift struct {
	ShiftID   string
	UserID    string
	Date      string // "YYYY-MM-DD"
	StartTime string // "HH:MM"
	EndTime   string // "HH:MM"
	Role      string // "CHEF" | "WAITER" | "MANAGER"
	Notes     string
	CreatedBy string
	CreatedAt time.Time
	UpdatedAt time.Time
}

var allowedRoles = map[string]bool{
	"CHEF": true, "WAITER": true, "MANAGER": true, "ADMIN": true,
}

func NewShift(userID, date, startTime, endTime, role, notes, createdBy string) (*Shift, error) {
	s := &Shift{
		UserID:    strings.TrimSpace(userID),
		Date:      strings.TrimSpace(date),
		StartTime: strings.TrimSpace(startTime),
		EndTime:   strings.TrimSpace(endTime),
		Role:      strings.TrimSpace(role),
		Notes:     strings.TrimSpace(notes),
		CreatedBy: strings.TrimSpace(createdBy),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Shift) Validate() error {
	if s.UserID == "" {
		return ErrUserIDRequired
	}
	if len(s.Date) != 10 {
		return ErrInvalidDate
	}
	if len(s.StartTime) < 4 || len(s.EndTime) < 4 {
		return ErrInvalidTime
	}
	if s.StartTime >= s.EndTime {
		return ErrEndTimeBeforeStart
	}
	if !allowedRoles[s.Role] {
		return ErrInvalidRole
	}
	return nil
}
