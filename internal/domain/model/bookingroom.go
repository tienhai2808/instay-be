package model

type BookingRoomStatus string

const (
	BookingRoomStatusReserved BookingRoomStatus = "reserved"
	BookingRoomStatusReleased BookingRoomStatus = "released"
)

type BookingRoom struct {
	ID        int64
	BookingID int64
	RoomID    int64
	Code      string
	Status    BookingRoomStatus

	Booking        *Booking
	Room           *Room
	Review         *Review
	Notifications  []*Notification
	OrderServices  []*OrderService
	SupportTickets []*SupportTicket
	Guests         []*Guest
	Customers      []*Customer
	Chats          []*Chat
}

func NewBookingRoom(
	id, bookingID, roomID int64,
	code string,
) *BookingRoom {
	return &BookingRoom{
		ID:        id,
		Code:      code,
		BookingID: bookingID,
		RoomID:    roomID,
		Status:    BookingRoomStatusReserved,
	}
}
