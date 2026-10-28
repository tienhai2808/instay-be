package entity

import "time"

const NotificationTableName = "notifications"

type NotificationEntity struct {
	ID            int64     `gorm:"type:bigint;primaryKey" json:"id"`
	DepartmentID  *int64    `gorm:"type:bigint" json:"department_id"`
	BookingRoomID *int64    `gorm:"type:bigint" json:"booking_room_id"`
	ReceiverType  string    `gorm:"type:varchar(20);not null;check:receiver_type IN ('guest', 'staff')" json:"receiver_type"`
	ContentType   string    `gorm:"type:varchar(20);not null;check:content_type IN ('order_service', 'support_ticket', 'booking', 'event')" json:"content_type"`
	ContentKey    string    `gorm:"type:varchar(255);not null" json:"content_key"`
	ContentVI     string    `gorm:"type:text;not null" json:"content_vi"`
	ContentEN     string    `gorm:"type:text;not null" json:"content_en"`
	IsRead        bool      `gorm:"type:boolean;not null" json:"is_read"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`

	Department  *DepartmentEntity  `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_notifications_department,OnUpdate:CASCADE,OnDelete:CASCADE" json:"department"`
	BookingRoom *BookingRoomEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_notifications_booking_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking_room"`
	Users       []*UserEntity      `gorm:"many2many:user_notifications;joinForeignKey:NotificationID;joinReferences:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"users"`
}

func (NotificationEntity) TableName() string {
	return NotificationTableName
}
