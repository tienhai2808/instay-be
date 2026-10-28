package entity

import "time"

const DailyRoomTypeRateTableName = "daily_room_type_rates"

type DailyRoomTypeRateEntity struct {
	RoomTypeID          int64     `gorm:"type:bigint;primaryKey" json:"room_type_id"`
	PricingCalendarDate time.Time `gorm:"type:date;primaryKey" json:"pricing_calendar_date"`
	FinalPrice          float64   `gorm:"type:decimal(10,2);not null" json:"final_price"`

	PricingCalendar *PricingCalendarEntity `gorm:"foreignKey:PricingCalendarDate;references:Date;constraint:fk_daily_room_type_rates_pricing_calendar,OnUpdate:CASCADE,OnDelete:CASCADE" json:"pricing_calendar"`
	RoomType        *RoomTypeEntity        `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_daily_room_type_rates_room_type,OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type"`
}

func (DailyRoomTypeRateEntity) TableName() string {
	return DailyRoomTypeRateTableName
}
