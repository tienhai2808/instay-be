package model

const (
	SourceOffline = "Trực tiếp (Offline)"
)

type Source struct {
	Base
	Audit

	Name string

	Bookings []*Booking
}

type SourceQuery struct {
	Page   uint32
	Limit  uint32
	Search string
}

func NewSource(
	id int64,
	name string,
	createdByID, updatedByID *int64,
) *Source {
	return &Source{
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
