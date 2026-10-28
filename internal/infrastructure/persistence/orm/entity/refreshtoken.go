package entity

import "time"

const RefreshTokenTableName = "refresh_tokens"

type RefreshTokenEntity struct {
	UserID    int64     `gorm:"type:bigint;not null;index:refresh_tokens_user_id_idx" json:"user_id"`
	Token     string    `gorm:"type:varchar(255);not null;primaryKey" json:"token"`
	ExpiresAt time.Time `gorm:"not null;index:refresh_tokens_expires_at_idx" json:"expires_at"`

	User *UserEntity `gorm:"foreignKey:UserID;references:ID;constraint:fk_refresh_tokens_user,OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
}

func (RefreshTokenEntity) TableName() string {
	return RefreshTokenTableName
}
