package model

import "time"

type RefreshToken struct {
	UserID    int64
	Token     string
	ExpiresAt time.Time

	User *User
}
