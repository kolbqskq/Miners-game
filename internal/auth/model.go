package auth

type RegisterSession struct {
	Email          string
	Code           string
	Username       string
	HashedPassword string
	ExpiresAt      int64
}
