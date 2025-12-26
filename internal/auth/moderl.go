package auth

type Session struct {
	ID        string
	UserID    string
	ExpiresAt int64
}
