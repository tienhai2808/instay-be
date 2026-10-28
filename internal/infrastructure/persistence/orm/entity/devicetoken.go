package entity

import "time"

const DeviceTokenTableName = "device_tokens"

type DeviceTokenEntity struct {
	ClientID   int64      `gorm:"type:bigint;not null;index:device_tokens_client_id_idx" json:"client_id"`
	ClientType string     `gorm:"type:varchar(20);not null;check:client_type IN ('user', 'booking_room')" json:"client_type"`
	Token      string     `gorm:"type:varchar(255);not null;primaryKey" json:"token"`
	LastUsedAt *time.Time `json:"last_used_at"`
}

func (DeviceTokenEntity) TableName() string {
	return DeviceTokenTableName
}
