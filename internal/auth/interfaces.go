package auth

type IAuthRepository interface {
	SaveSession(session *Session) error
}
