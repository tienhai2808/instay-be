package model

import "time"

type DailyRoomTypeRate struct {
	RoomTypeID          int64
	PricingCalendarDate time.Time
	FinalPrice          float64

	RoomType        *RoomType
	PricingCalendar *PricingCalendar
}
