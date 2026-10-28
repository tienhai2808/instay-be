package entity

const RoomTypePriceTableName = "room_type_prices"

type RoomTypePriceEntity struct {
	ID           int64   `gorm:"type:bigint;primaryKey" json:"id"`
	RoomTypeID   int64   `gorm:"type:bigint;not null;uniqueIndex:room_type_prices_room_type_id_pricing_tag_id_key" json:"room_type_id"`
	PricingTagID int64   `gorm:"type:bigint;not null;uniqueIndex:room_type_prices_room_type_id_pricing_tag_id_key" json:"pricing_tag_id"`
	PricingType  string  `gorm:"type:varchar(20);not null;check:pricing_type IN ('percent', 'amount')" json:"pricing_type"`
	PricingValue float64 `gorm:"type:decimal(10,2);not null" json:"pricing_value"`

	RoomType   *RoomTypeEntity   `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_room_type_prices_room_type,OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type"`
	PricingTag *PricingTagEntity `gorm:"foreignKey:PricingTagID;references:ID;constraint:fk_room_type_prices_pricing_tag,OnUpdate:CASCADE,OnDelete:CASCADE" json:"pricing_tag"`
}

func (RoomTypePriceEntity) TableName() string {
	return RoomTypePriceTableName
}
