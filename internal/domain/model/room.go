package model

import "time"

type Room struct {
	Base
	Audit

	Name        string
	FloorID     int64
	RoomTypeID  int64
	IsAvailable bool

	RoomType     *RoomType
	Floor        *Floor
	BookingRooms []*BookingRoom
}

type RoomQuery struct {
	Page       uint32
	Limit      uint32
	Sort       string
	Order      string
	BranchID   int64
	BranchIDs  []int64
	FloorID    int64
	RoomTypeID int64
	Search     string
}

type RoomStatusQuery struct {
	RoomTypeID int64
	BranchID   int64
	BranchIDs  []int64
	SnapshotAt time.Time
}

type RoomTimelineQuery struct {
	BranchID          int64
	BranchIDs         []int64
	RoomTypeID        int64
	FloorID           int64
	Search            string
	SnapshotAt        time.Time
	UpcomingThreshold time.Duration
}

type RoomStatusSummary struct {
	All       int64
	Available int64
	Occupied  int64
	Booked    int64
}

type RoomTimeline struct {
	ServerNow time.Time
	Status    *RoomTimelineStatusSummary
	Data      []*RoomTimelineRoomType
}

type RoomTimelineStatusSummary struct {
	All             int64
	Available       int64
	Occupied        int64
	UpcomingCheckIn int64
	OverdueCheckout int64
	OverdueCheckIn  int64
}

type RoomTimelineRoomType struct {
	RoomTypeID   int64
	RoomTypeName string
	RoomsCount   int64
	Rooms        []*RoomTimelineRoom
}

type RoomTimelineRoom struct {
	RoomID   int64
	RoomName string
	Price    float64
	Booking  *RoomTimelineBooking
}

type RoomTimelineBooking struct {
	BookingID     int64
	BookingNumber string
	Status        BookingStatus
	SourceName    string
	GuestName     string
	CheckIn       time.Time
	CheckOut      time.Time
}

func NewRoom(
	id int64,
	name string,
	floorID, roomTypeID int64,
	createdByID, updatedByID *int64,
) *Room {
	return &Room{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:       name,
		RoomTypeID: roomTypeID,
		FloorID:    floorID,
	}
}
