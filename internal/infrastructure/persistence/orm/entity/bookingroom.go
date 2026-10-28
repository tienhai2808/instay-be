package entity

const BookingRoomTableName = "booking_rooms"

type BookingRoomEntity struct {
	ID        int64  `gorm:"type:bigint;primaryKey" json:"id"`
	BookingID int64  `gorm:"type:bigint;not null;uniqueIndex:booking_rooms_booking_id_room_id_key" json:"booking_id"`
	RoomID    int64  `gorm:"type:bigint;not null;uniqueIndex:booking_rooms_booking_id_room_id_key" json:"room_id"`
	Code      string `gorm:"type:varchar(50);not null;uniqueIndex:booking_rooms_code_key" json:"code"`
	Status    string `gorm:"type:varchar(20);not null;check:status IN ('reserved', 'released')"`

	Booking        *BookingEntity         `gorm:"foreignKey:BookingID;references:ID;constraint:fk_booking_rooms_booking,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking"`
	Room           *RoomEntity            `gorm:"foreignKey:RoomID;references:ID;constraint:fk_booking_rooms_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room"`
	Review         *ReviewEntity          `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_reviews_booking_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"review"`
	OrderServices  []*OrderServiceEntity  `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_order_services_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"order_services"`
	SupportTickets []*SupportTicketEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_support_tickets_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"support_tickets"`
	Guests         []*GuestEntity         `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_guests_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"guests"`
	Customers      []*CustomerEntity      `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_customers_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"customers"`
	Chats          []*ChatEntity          `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_chats_booking_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"chats"`
}

func (BookingRoomEntity) TableName() string {
	return BookingRoomTableName
}
