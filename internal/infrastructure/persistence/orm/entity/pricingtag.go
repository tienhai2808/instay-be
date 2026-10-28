package entity

const PricingTagTableName = "pricing_tags"

type PricingTagEntity struct {
	BaseEntity
	AuditEntity

	Name string `gorm:"type:varchar(150);not null;uniqueIndex:pricing_tags_name_key" json:"name"`

	PricingCalendars []*PricingCalendarEntity `gorm:"foreignKey:PricingTagID;references:ID;constraint:fk_pricing_calendars_pricing_tag,OnUpdate:CASCADE,OnDelete:CASCADE" json:"pricing_calendars"`
	RoomTypePrices   []*RoomTypePriceEntity   `gorm:"foreignKey:PricingTagID;references:ID;constraint:fk_room_type_prices_pricing_tag,OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type_prices"`
}

func (PricingTagEntity) TableName() string {
	return PricingTagTableName
}
