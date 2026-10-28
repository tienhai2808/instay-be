package entity

import "time"

const GuestTableName = "guests"

type GuestEntity struct {
	BaseEntity
	AuditEntity

	FullName       string    `gorm:"type:varchar(150);not null" json:"full_name"`
	DOB            time.Time `gorm:"type:date;not null" json:"dob"`
	National       string    `gorm:"type:varchar(10);not null" json:"national"`
	PassportNumber string    `gorm:"type:varchar(20);not null;uniqueIndex:guests_passport_number_booking_room_id_key" json:"passport_number"`
	CheckIn        time.Time `gorm:"type:date;not null" json:"check_in"`
	CheckOut       time.Time `gorm:"type:date;not null" json:"check_out"`
	BookingRoomID  int64     `gorm:"type:bigint;not null;uniqueIndex:guests_passport_number_booking_room_id_key" json:"booking_room_id"`

	BookingRoom *BookingRoomEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_guests_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking_room"`
}

func (GuestEntity) TableName() string {
	return GuestTableName
}
