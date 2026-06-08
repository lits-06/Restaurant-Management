package domain

// AuthSession represents an active authentication session
type AuthSession struct {
	UserID       string
	Email        string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64 // seconds
}

// NewAuthSession creates a new authentication session
func NewAuthSession(userID, email, accessToken, refreshToken string, expiresIn int64) *AuthSession {
	return &AuthSession{
		UserID:       userID,
		Email:        email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
	}
}
