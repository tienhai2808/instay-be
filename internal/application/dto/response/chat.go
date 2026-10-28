package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type BasicChatResponse struct {
	ID                        int64                        `json:"id,string"`
	Department                *BasicDepartmentResponse     `json:"department"`
	LastMessageAt             *time.Time                   `json:"last_message_at"`
	LastMessageContent        *string                      `json:"last_message_content"`
	LastMessageAttachmentType *model.MessageAttachmentType `json:"last_message_attachment_type"`
	LastMessageSenderType     *model.MessageActionType     `json:"last_message_sender_type"`
	LastMessageSenderName     *string                      `json:"last_message_sender_name"`
	ExpiredAt                 time.Time                    `json:"expired_at"`
	CreatedAt                 time.Time                    `json:"created_at"`
}

type SimpleChatResponse struct {
	ID                        int64                        `json:"id,string"`
	BookingRoom               *SimpleBookingRoomResponse   `json:"booking_room"`
	LastMessageAt             *time.Time                   `json:"last_message_at"`
	LastMessageContent        *string                      `json:"last_message_content"`
	LastMessageAttachmentType *model.MessageAttachmentType `json:"last_message_attachment_type"`
	LastMessageSenderType     *model.MessageActionType     `json:"last_message_sender_type"`
	LastMessageSenderName     *string                      `json:"last_message_sender_name"`
	ExpiredAt                 time.Time                    `json:"expired_at"`
	CreatedAt                 time.Time                    `json:"created_at"`
}

type MessageResponse struct {
	ID                string                       `json:"id"`
	ChatID            int64                        `json:"chat_id,string"`
	SenderType        model.MessageActionType      `json:"sender_type"`
	SenderID          int64                        `json:"sender_id,string"`
	SenderName        string                       `json:"sender_name"`
	AttachmentKey     *string                      `json:"attachment_key"`
	AttachmentType    *model.MessageAttachmentType `json:"attachment_type"`
	Content           *string                      `json:"content"`
	BookingRoomReadAt *time.Time                   `json:"booking_room_read_at"`
	MessageUsers      []*MessageUserResponse       `json:"message_users"`
	IsDeleted         bool                         `json:"is_deleted"`
	MessageReactions  []*MessageReactionResponse   `json:"message_reaction"`
	CreatedAt         time.Time                    `json:"created_at"`
	MessageReply      *MessageReplyResponse        `json:"message_reply"`
}

type UpdateMessageReadResponse struct {
	ChatID     int64                   `json:"chat_id,string"`
	ReaderName string                  `json:"reader_name"`
	ReaderType model.MessageActionType `json:"reader_type"`
	ReaderID   int64                   `json:"reader_id,string"`
}

type TypingMessageResponse struct {
	ChatID     int64                   `json:"chat_id,string"`
	SenderType model.MessageActionType `json:"sender_type"`
}

type UpdateMessageReactionResponse struct {
	MessageID        string                         `json:"message_id"`
	ChatID           int64                          `json:"chat_id,string"`
	SenderType       model.MessageActionType        `json:"sender_type"`
	SenderID         int64                          `json:"sender_id,string"`
	SenderName       string                         `json:"sender_name"`
	Emoji            model.MessageReactionEmojiType `json:"emoji"`
	IsReacted        bool                           `json:"is_reacted"`
	MessageReactions []*MessageReactionResponse     `json:"message_reactions"`
}

type MessageReplyResponse struct {
	MessageID      string                       `json:"message_id"`
	SenderName     string                       `json:"sender_name"`
	Content        *string                      `json:"content"`
	AttachmentType *model.MessageAttachmentType `json:"attachment_type"`
}

type MessageReactionResponse struct {
	SenderType string                         `json:"sender_type"`
	SenderID   int64                          `json:"sender_id,string"`
	SenderName string                         `json:"sender_name"`
	Emoji      model.MessageReactionEmojiType `json:"emoji"`
	CreatedAt  time.Time                      `json:"created_at"`
}

type MessageUserResponse struct {
	UserID       int64     `json:"user_id,string"`
	UserFullName string    `json:"user_full_name"`
	ReadAt       time.Time `json:"read_at"`
}

type ChatAvailabilityForGuestResponse struct {
	DepartmentID int64 `json:"department_id,string"`
	ChatID       int64 `json:"chat_id,string"`
	Available    bool  `json:"available"`
}

type ChatAvailabilityForAdminResponse struct {
	ChatID    int64 `json:"chat_id,string"`
	Available bool  `json:"available"`
}
