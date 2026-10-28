package model

type RoomType struct {
	Base
	Audit

	Name        string
	Price       float64
	SurfaceArea float32
	MaxPeople   uint8

	Branches       []*Branch
	Rooms          []*Room
	Bookings       []*Booking
	DailyRoomRates []*DailyRoomTypeRate
	RoomTypePrices []*RoomTypePrice
	RoomsCount     int64
}

type RoomTypeQuery struct {
	Page      uint32
	Limit     uint32
	Search    string
	BranchID  int64
	BranchIDs []int64
}

func NewRoomType(
	id int64,
	name string,
	price float64,
	surfaceArea float32,
	maxPeople uint8,
	createdByID, updatedByID *int64,
) *RoomType {
	return &RoomType{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:        name,
		Price:       price,
		SurfaceArea: surfaceArea,
		MaxPeople:   maxPeople,
	}
}
