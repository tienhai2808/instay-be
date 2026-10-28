package entity

import "time"

type BaseEntity struct {
	ID        int64     `gorm:"type:bigint;primaryKey" json:"id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
