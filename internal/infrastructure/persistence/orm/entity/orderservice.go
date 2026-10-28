package entity

import "time"

const OrderServiceTableName = "order_services"

type OrderServiceEntity struct {
	BaseEntity
	AuditEntity

	ServiceID     int64     `gorm:"type:bigint;not null" json:"service_id"`
	BookingRoomID int64     `gorm:"type:bigint;not null" json:"booking_room_id"`
	Quantity      uint32    `gorm:"type:integer;not null" json:"quantity"`
	UnitPrice     float64   `gorm:"type:decimal(10,2);not null" json:"unit_price"`
	TotalPrice    float64   `gorm:"type:decimal(10,2);not null" json:"total_price"`
	Status        string    `gorm:"type:varchar(20);not null;check:status IN ('pending', 'confirmed', 'cancelled')" json:"status"`
	GuestNote     *string   `gorm:"type:text" json:"guest_note"`
	StaffNote     *string   `gorm:"type:text" json:"staff_note"`
	CancelReason  *string   `gorm:"type:text" json:"cancel_reason"`
	EstimatedAt   time.Time `json:"estimated_at"`

	Service     *ServiceEntity     `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_order_services_service,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service"`
	BookingRoom *BookingRoomEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_order_services_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking_room"`
}

func (OrderServiceEntity) TableName() string {
	return OrderServiceTableName
}
