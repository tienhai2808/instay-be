package entity

import "time"

type MessageReactionEntity struct {
	SenderType string    `bson:"sender_type" json:"sender_type" validate:"required,oneof=staff guest"`
	SenderID   int64     `bson:"sender_id" json:"sender_id"`
	SenderName string    `bson:"sender_name" json:"sender_name"`
	Emoji      string    `bson:"emoji" json:"emoji"`
	CreatedAt  time.Time `bson:"created_at" json:"created_at"`
}
