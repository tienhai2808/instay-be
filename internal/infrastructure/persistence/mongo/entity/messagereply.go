package entity

import "go.mongodb.org/mongo-driver/v2/bson"

type MessageReplyEntity struct {
	MessageID      bson.ObjectID `bson:"message_id" json:"message_id"`
	SenderName     string        `bson:"sender_name" json:"sender_name"`
	Content        *string       `bson:"content,omitempty" json:"content"`
	AttachmentType *string       `bson:"attachment_type,omitempty" json:"attachment_type" validate:"omitempty,oneof=file image"`
}
