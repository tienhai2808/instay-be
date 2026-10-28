package model

type RoomTypePricePricingType string

const (
	RoomTypePricePricingTypePercent RoomTypePricePricingType = "percent"
	RoomTypePricePricingTypeAmount  RoomTypePricePricingType = "amount"
)

type RoomTypePrice struct {
	ID           int64
	RoomTypeID   int64
	PricingTagID int64
	PricingType  RoomTypePricePricingType
	PricingValue float64

	RoomType   *RoomType
	PricingTag *PricingTag
}
