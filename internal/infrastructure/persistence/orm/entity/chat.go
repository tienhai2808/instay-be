package entity

import "time"

const ChatTableName = "chats"

type ChatEntity struct {
	ID                        int64      `gorm:"type:bigint;primaryKey" json:"id"`
	CreatedAt                 time.Time  `gorm:"autoCreateTime" json:"created_at"`
	BookingRoomID             int64      `gorm:"type:bigint;not null;uniqueIndex:chats_booking_room_id_department_id" json:"booking_room_id"`
	DepartmentID              int64      `gorm:"type:bigint;not null;uniqueIndex:chats_booking_room_id_department_id" json:"department_id"`
	ExpiredAt                 time.Time  `gorm:"not null" json:"expired_at"`
	LastMessageAt             *time.Time `json:"last_message_at"`
	LastMessageContent        *string    `gorm:"type:varchar(50)" json:"last_message_content"`
	LastMessageAttachmentType *string    `gorm:"type:varchar(50);check:last_message_attachment_type IN ('image', 'file')" json:"last_message_attachment_type"`
	LastMessageSenderType     *string    `gorm:"type:varchar(20);check:last_message_sender_type IN ('guest', 'staff')" json:"last_message_sender_type"`
	LastMessageSenderName     *string    `gorm:"type:varchar(50)" json:"last_message_sender_name"`

	Department  *DepartmentEntity  `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_chats_department,OnUpdate:CASCADE,OnDelete:CASCADE" json:"department"`
	BookingRoom *BookingRoomEntity `gorm:"foreignKey:BookingRoomID;references:ID;constraint:fk_chats_booking_room,OnUpdate:CASCADE,OnDelete:CASCADE" json:"booking_room"`
}

func (ChatEntity) TableName() string {
	return ChatTableName
}
