package entity

import "gorm.io/datatypes"

const ReviewTableName = "reviews"

type ReviewEntity struct {
	BaseEntity

	BookingRoomID int64          `gorm:"type:bigint;not null;uniqueIndex:reviews_booking_room_id" json:"booking_room_id"`
	Star          uint8          `gorm:"type:integer;not null" json:"star"`
	Content       string         `gorm:"type:text;not null" json:"content"`
	Aspects       datatypes.JSON `gorm:"type:jsonb" json:"aspects"`

	BookingRoom *BookingRoomEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_reviews_booking_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking_room"`
}

func (ReviewEntity) TableName() string {
	return ReviewTableName
}
