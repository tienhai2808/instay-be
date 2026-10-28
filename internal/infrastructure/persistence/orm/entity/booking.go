package entity

import (
	"time"

	"gorm.io/datatypes"
)

const BookingTableName = "bookings"

type BookingEntity struct {
	BaseEntity
	AuditEntity

	BookingNumber      string         `gorm:"type:varchar(50);not null;uniqueIndex:bookings_booking_number_key" json:"booking_number"`
	GuestFullName      string         `gorm:"type:varchar(150);not null" json:"guest_full_name"`
	GuestEmail         *string        `gorm:"type:varchar(150)" json:"guest_email"`
	GuestPhone         *string        `gorm:"type:varchar(20)" json:"guest_phone"`
	CheckIn            time.Time      `gorm:"not null" json:"check_in"`
	CheckOut           time.Time      `gorm:"not null" json:"check_out"`
	RoomNumber         uint32         `gorm:"type:integer;not null" json:"room_number"`
	GuestNumber        string         `gorm:"type:varchar(50);not null" json:"guest_number"`
	BookedOn           time.Time      `gorm:"type:date;not null" json:"booked_on"`
	TotalNetPrice      float64        `gorm:"type:decimal(10,2);not null" json:"total_net_price"`
	TotalSellPrice     float64        `gorm:"type:decimal(10,2);not null" json:"total_sell_price"`
	PromotionName      *string        `gorm:"type:varchar(150)" json:"promotion_name"`
	MealPlan           *string        `gorm:"type:varchar(150)" json:"meal_plan"`
	BookingPreferences *string        `gorm:"type:varchar(255)" json:"booking_preferences"`
	BookingConditions  *string        `gorm:"type:varchar(255)" json:"booking_conditions"`
	Note               *string        `gorm:"type:text" json:"note"`
	Status             string         `gorm:"type:varchar(50);not null;check:status IN ('confirmed', 'checked-in', 'checked-out', 'cancelled')" json:"status"`
	BranchID           *int64         `gorm:"type:bigint" json:"branch_id"`
	SourceID           int64          `gorm:"type:bigint;not null" json:"source_id"`
	RoomTypeID         int64          `gorm:"type:bigint;not null" json:"room_type_id"`
	Predict            datatypes.JSON `gorm:"type:jsonb" json:"predict"`

	Branch       *BranchEntity        `gorm:"foreignKey:BranchID;references:ID;constraint:fk_bookings_branch,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"branch"`
	Source       *SourceEntity        `gorm:"foreignKey:SourceID;references:ID;constraint:fk_bookings_source,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"source"`
	RoomType     *RoomTypeEntity      `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_bookings_room_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room_type"`
	BookingRooms []*BookingRoomEntity `gorm:"foreignKey:BookingID;references:ID;constraint:fk_booking_rooms_booking,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking_rooms"`
}

func (BookingEntity) TableName() string {
	return BookingTableName
}
