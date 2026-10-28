package model

import "time"

type (
	MessageActionType     string
	MessageAttachmentType string
)

const (
	MessageActionTypeGuest     MessageActionType     = "guest"
	MessageActionTypeStaff     MessageActionType     = "staff"
	MessageAttachmentTypeFile  MessageAttachmentType = "file"
	MessageAttachmentTypeImage MessageAttachmentType = "image"
)

type Message struct {
	ID                string
	Schema            string
	ChatID            int64
	SenderType        MessageActionType
	SenderID          int64
	SenderName        string
	AttachmentKey     *string
	AttachmentType    *MessageAttachmentType
	Content           *string
	BookingRoomReadAt *time.Time
	MessageUsers      []*MessageUser
	IsDeleted         bool
	MessageReactions  []*MessageReaction
	CreatedAt         time.Time
	MessageReply      *MessageReply
}

func NewMessage(
	id, schema string,
	chatID int64,
	senderType MessageActionType,
	senderID int64,
	senderName string,
	attachmentKey *string,
	attachmentType *MessageAttachmentType,
	content *string,
	bookingRoomReadAt *time.Time,
	isDeleted bool,
	messageReply *MessageReply,
	createdAt time.Time,
) *Message {
	return &Message{
		ID:                id,
		Schema:            schema,
		ChatID:            chatID,
		SenderType:        senderType,
		SenderID:          senderID,
		SenderName:        senderName,
		AttachmentKey:     attachmentKey,
		AttachmentType:    attachmentType,
		Content:           content,
		BookingRoomReadAt: bookingRoomReadAt,
		IsDeleted:         isDeleted,
		MessageReply:      messageReply,
		CreatedAt:         createdAt,
	}
}
