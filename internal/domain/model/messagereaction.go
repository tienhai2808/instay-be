package model

import "time"

type MessageReactionEmojiType string

const (
	MessageReactionEmojiTypeHeart MessageReactionEmojiType = "heart"
	MessageReactionEmojiTypeLike  MessageReactionEmojiType = "like"
	MessageReactionEmojiTypeSad   MessageReactionEmojiType = "sad"
	MessageReactionEmojiTypeFun   MessageReactionEmojiType = "fun"
	MessageReactionEmojiTypeAngry MessageReactionEmojiType = "angry"
)

type MessageReaction struct {
	SenderType string
	SenderID   int64
	SenderName string
	Emoji      MessageReactionEmojiType
	CreatedAt  time.Time
}
