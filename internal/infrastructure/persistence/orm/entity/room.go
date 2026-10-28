package entity

const RoomTableName = "rooms"

type RoomEntity struct {
	BaseEntity
	AuditEntity

	Name        string `gorm:"type:varchar(150);not null;uniqueIndex:rooms_name_key" json:"name"`
	FloorID     int64  `gorm:"type:bigint;not null" json:"floor_id"`
	RoomTypeID  int64  `gorm:"type:bigint;not null" json:"room_type_id"`
	IsAvailable bool   `gorm:"column:is_available;<-:false" json:"is_available"`

	RoomType     *RoomTypeEntity      `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_rooms_room_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"room_type"`
	Floor        *FloorEntity         `gorm:"foreignKey:FloorID;references:ID;constraint:fk_rooms_floor,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"floor"`
	BookingRooms []*BookingRoomEntity `gorm:"foreignKey:RoomID;references:ID;constraint:fk_booking_rooms_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking_rooms"`
}

func (RoomEntity) TableName() string {
	return RoomTableName
}
