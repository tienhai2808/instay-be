package model

type MessageReply struct {
	MessageID      string
	SenderName     string
	Content        *string
	AttachmentType *MessageAttachmentType
}
