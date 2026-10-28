package port

import "context"

type VerificationSession struct {
	SessionID       string
	VerificationURL string
}

type KYCProvider interface {
	CreateVerificationSession(ctx context.Context, vendorData string, callbackURL string) (*VerificationSession, error)
}
