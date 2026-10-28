package model

import "time"

type DeviceTokenClientType string

const (
	DeviceTokenClientTypeUser        DeviceTokenClientType = "user"
	DeviceTokenClientTypeBookingRoom DeviceTokenClientType = "booking_room"
)

type DeviceToken struct {
	ClientID   int64
	ClientType DeviceTokenClientType
	Token      string
	LastUsedAt *time.Time
}
