package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

const MessageCollection = "messages"

type MessageEntity struct {
	ID                bson.ObjectID            `bson:"_id,omitempty" json:"id"`
	Schema            string                   `bson:"schema" json:"schema"`
	ChatID            int64                    `bson:"chat_id" json:"chat_id"`
	SenderType        string                   `bson:"sender_type" json:"sender_type" validate:"required,oneof=staff guest"`
	SenderID          int64                    `bson:"sender_id" json:"sender_id"`
	SenderName        string                   `bson:"sender_name" json:"sender_name"`
	AttachmentKey     *string                  `bson:"attachment_key,omitempty" json:"attachment_key"`
	AttachmentType    *string                  `bson:"attachment_type,omitempty" json:"attachment_type" validate:"omitempty,oneof=file image"`
	Content           *string                  `bson:"content,omitempty" json:"content" validate:"required,min=2"`
	MessageReply      *MessageReplyEntity      `bson:"message_reply,omitempty" json:"message_reply"`
	BookingRoomReadAt *time.Time               `bson:"booking_room_read_at,omitempty" json:"booking_room_read_at"`
	MessageUsers      []*MessageUserEntity     `bson:"message_users,omitempty" json:"message_users"`
	IsDeleted         bool                     `bson:"is_deleted" json:"is_deleted"`
	MessageReactions  []*MessageReactionEntity `bson:"message_reactions,omitempty" json:"message_reactions"`
	CreatedAt         time.Time                `bson:"created_at" json:"created_at"`
}
