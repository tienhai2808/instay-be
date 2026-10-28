package entity

import "time"

const CustomerTableName = "customers"

type CustomerEntity struct {
	ID            int64     `gorm:"type:bigint;primaryKey" json:"id"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	BirthYear     uint32    `gorm:"type:integer;not null" json:"birth_year"`
	Phone         string    `gorm:"type:varchar(20);not null" json:"phone"`
	Gender        string    `gorm:"type:varchar(10);not null;check:gender IN ('male', 'female')" json:"gender"`
	BookingRoomID int64     `gorm:"type:bigint;not null" json:"booking_room_id"`

	BookingRoom *BookingRoomEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_customers_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking_room"`
}

func (CustomerEntity) TableName() string {
	return CustomerTableName
}
