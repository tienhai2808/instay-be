package model

import "time"

type MessageUser struct {
	UserID       int64
	UserFullName string
	ReadAt       time.Time
}
