package entity

const SupportTicketTableName = "support_tickets"

type SupportTicketEntity struct {
	BaseEntity

	BookingRoomID     int64  `gorm:"type:bigint;not null" json:"booking_room_id"`
	Content           string `gorm:"type:text;not null" json:"content"`
	Status            string `gorm:"type:varchar(20);not null;check:status IN ('pending', 'accepted', 'cancelled', 'done')" json:"status"`
	Priority          string `gorm:"type:varchar(20);not null;check:priority IN ('low', 'medium', 'high', 'urgent')" json:"priority"`
	ResolutionSeconds *int64 `gorm:"type:bigint" json:"resolution_seconds"`
	SupportCategoryID int64  `gorm:"type:bigint;not null" json:"support_category_id"`
	UpdatedByID       *int64 `gorm:"type:bigint" json:"updated_by_id"`

	SupportCategory *SupportCategoryEntity `gorm:"foreignKey:SupportCategoryID;references:ID;constraint:fk_support_tickets_support_category,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"support_category"`
	BookingRoom     *BookingRoomEntity     `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_support_tickets_booking_room,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"booking_room"`
	UpdatedBy       *UserEntity            `gorm:"belongsTo;foreignKey:UpdatedByID;references:ID;constraint:-" json:"updated_by"`
}

func (SupportTicketEntity) TableName() string {
	return SupportTicketTableName
}
