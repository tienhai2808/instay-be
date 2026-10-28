package entity

import "time"

type MessageUserEntity struct {
	UserID       int64     `bson:"user_id" json:"user_id"`
	UserFullName string    `bson:"user_full_name" json:"user_full_name"`
	ReadAt       time.Time `bson:"read_at" json:"read_at"`
}
