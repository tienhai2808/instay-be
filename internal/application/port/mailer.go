package port

type MailerProvider interface {
	AuthEmail(to, subject, otp string) error
}
