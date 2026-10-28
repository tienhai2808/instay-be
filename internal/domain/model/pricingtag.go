package model

type PricingTag struct {
	Base
	Audit

	Name string

	PricingCalendars []*PricingCalendar
	RoomTypePrices   []*RoomTypePrice
}

type PricingTagQuery struct {
	Page   uint32
	Limit  uint32
	Search string
}

func NewPricingTag(id int64, name string, createdByID, updatedByID *int64) *PricingTag {
	return &PricingTag{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name: name,
	}
}
