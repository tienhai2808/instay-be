package entity

import "time"

const PricingCalendarTableName = "pricing_calendars"

type PricingCalendarEntity struct {
	Date         time.Time `gorm:"type:date;primaryKey" json:"date"`
	PricingTagID *int64    `gorm:"type:bigint" json:"pricing_tag_id"`

	PricingTag         *PricingTagEntity          `gorm:"foreignKey:PricingTagID;references:ID;constraint:fk_pricing_calendars_pricing_tag,OnUpdate:CASCADE,OnDelete:CASCADE" json:"pricing_tag"`
	DailyRoomTypeRates []*DailyRoomTypeRateEntity `gorm:"foreignKey:PricingCalendarDate;references:Date;constraint:fk_daily_room_type_rates_pricing_calendar,OnUpdate:CASCADE,OnDelete:CASCADE" json:"daily_room_type_rates"`
}

func (PricingCalendarEntity) TableName() string {
	return PricingCalendarTableName
}
