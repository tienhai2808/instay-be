package model

import "time"

type PricingCalendar struct {
	Date         time.Time
	PricingTagID *int64

	PricingTag         *PricingTag
	DailyRoomTypeRates []*DailyRoomTypeRate
}
