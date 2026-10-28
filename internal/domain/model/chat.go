package model

import "time"

type Chat struct {
	ID                        int64
	CreatedAt                 time.Time
	BookingRoomID             int64
	DepartmentID              int64
	ExpiredAt                 time.Time
	LastMessageAt             *time.Time
	LastMessageContent        *string
	LastMessageAttachmentType *MessageAttachmentType
	LastMessageSenderType     *MessageActionType
	LastMessageSenderName     *string

	Department  *Department
	BookingRoom *BookingRoom
}

type ChatQuery struct {
	Page      uint32
	Limit     uint32
	Search    string
	BranchIDs []int64
}

func NewChat(
	id, bookingRoomID, departmentID int64,
	expiredAt time.Time,
) *Chat {
	return &Chat{
		ID:            id,
		BookingRoomID: bookingRoomID,
		DepartmentID:  departmentID,
		ExpiredAt:     expiredAt,
	}
}
