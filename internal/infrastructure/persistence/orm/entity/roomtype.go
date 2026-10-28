package entity

const RoomTypeTableName = "room_types"

type RoomTypeEntity struct {
	BaseEntity
	AuditEntity

	Name        string  `gorm:"type:varchar(150);not null;uniqueIndex:room_types_name_key" json:"name"`
	Price       float64 `gorm:"type:decimal(10,2);not null" json:"price"`
	SurfaceArea float32 `gorm:"type:decimal(5,2);not null" json:"surface_area"`
	MaxPeople   uint8   `gorm:"type:integer;not null" json:"max_people"`

	Rooms          []*RoomEntity              `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_rooms_room_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"rooms"`
	Bookings       []*BookingEntity           `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_bookings_room_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"bookings"`
	DailyRoomRates []*DailyRoomTypeRateEntity `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_daily_room_type_rates_room_type,OnUpdate:CASCADE,OnDelete:CASCADE" json:"daily_room_rates"`
	RoomTypePrices []*RoomTypePriceEntity     `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_room_type_prices_room_type,OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type_prices"`
	Branches       []*BranchEntity            `gorm:"many2many:branch_room_types;foreignKey:ID;joinForeignKey:RoomTypeID;references:ID;joinReferences:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branches"`
	RoomsCount     int64                      `gorm:"-" json:"rooms_count"`
}

func (RoomTypeEntity) TableName() string {
	return RoomTypeTableName
}
