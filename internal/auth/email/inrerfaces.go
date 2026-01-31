package email

type IEmailService interface {
	Send(to, code string) error
}
