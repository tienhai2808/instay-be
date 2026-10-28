package data

import "time"

type ForgotPasswordData struct {
	Email          string `json:"email"`
	Otp            string `json:"otp"`
	VerifyAttempts int    `json:"verify_attempts"`
	ResendAttempts int    `json:"resend_attempts"`
}

type AuthEmailMessage struct {
	Schema  string `json:"schema"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Otp     string `json:"otp"`
}

type DeleteFileMessage struct {
	Schema string   `json:"schema"`
	Keys   []string `json:"keys"`
}

type PushNotificationData struct {
	Data map[string]string `json:"data"`
}

type PushNotificationMessage struct {
	Schema string            `json:"schema"`
	Tokens []string          `json:"tokens"`
	Data   map[string]string `json:"data"`
}

type BookingRoomData struct {
	ID        int64     `json:"id"`
	ExpiredAt time.Time `json:"expired_at"`
}
