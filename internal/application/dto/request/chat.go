package request

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

type ChatPaginationQuery struct {
	Page   uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit  uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search string `form:"search" json:"search"`
}

type MessagePaginationQuery struct {
	Page  uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
}

type ChatAvailabilityForGuestQuery struct {
	ServiceID         int64 `form:"service_id" binding:"omitempty,min=1" json:"service_id,string"`
	SupportCategoryID int64 `form:"support_category_id" binding:"omitempty,min=1" json:"support_category_id,string"`
}

type ChatAvailabilityForAdminQuery struct {
	BookingRoomID int64 `form:"booking_room_id" binding:"required,min=1" json:"booking_room_id,string"`
}

type CreateMessageRequest struct {
	AttachmentType *model.MessageAttachmentType `json:"attachment_type" validate:"omitempty"`
	Content        *string                      `json:"content" validate:"omitempty"`
	AttachmentKey  *string                      `json:"attachment_key" validate:"omitempty"`
	ChatID         *int64                       `json:"chat_id,string" validate:"omitempty"`
	DepartmentID   int64                        `json:"department_id,string" validate:"required"`
	BookingRoomID  int64                        `json:"booking_room_id,string" validate:"required"`
	MessageID      *string                      `json:"message_id" validate:"omitempty"`
	SenderID       int64                        `json:"sender_id,string" validate:"required"`
	SenderType     model.MessageActionType      `json:"sender_type" validate:"required,oneof=guest staff"`
	SenderName     string                       `json:"sender_name" validate:"required"`
}

type UpdateMessageReadRequest struct {
	ChatID        int64                   `json:"chat_id,string" validate:"required"`
	DepartmentID  int64                   `json:"department_id,string" validate:"required"`
	BookingRoomID int64                   `json:"booking_room_id,string" validate:"required"`
	ReaderName    string                  `json:"reader_name" validate:"required"`
	ReaderType    model.MessageActionType `json:"reader_type" validate:"required,oneof=guest staff"`
	ReaderID      int64                   `json:"reader_id,string" validate:"required"`
}

type UpdateMessageReactionRequest struct {
	MessageID     string                         `json:"message_id" validate:"required"`
	ChatID        int64                          `json:"chat_id,string" validate:"required"`
	DepartmentID  int64                          `json:"department_id,string" validate:"required"`
	BookingRoomID int64                          `json:"booking_room_id,string" validate:"required"`
	SenderName    string                         `json:"sender_name" validate:"required"`
	SenderType    model.MessageActionType        `json:"sender_type" validate:"required,oneof=guest staff"`
	SenderID      int64                          `json:"sender_id,string" validate:"required"`
	Emoji         model.MessageReactionEmojiType `json:"emoji" validate:"required,oneof=heart like sad fun angry"`
}

type TypingMessageRequest struct {
	ChatID        int64                   `json:"chat_id,string" validate:"required"`
	DepartmentID  int64                   `json:"department_id,string" validate:"required"`
	BookingRoomID int64                   `json:"booking_room_id,string" validate:"required"`
	SenderType    model.MessageActionType `json:"sender_type" validate:"required,oneof=guest staff"`
}
