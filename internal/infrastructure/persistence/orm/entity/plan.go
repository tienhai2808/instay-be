package entity

import "gorm.io/datatypes"

const PlanTableName = "plans"

type PlanEntity struct {
	ID           int64          `gorm:"type:bigint;primaryKey" json:"id"`
	Name         string         `gorm:"type:varchar(150);not null;uniqueIndex:plans_name_key" json:"name"`
	Price        float64        `gorm:"type:decimal(10,2);not null" json:"price"`
	Currency     string         `gorm:"type:varchar(20);not null;check:currency IN ('usd', 'vnd')" json:"currency"`
	Features     datatypes.JSON `gorm:"type:jsonb" json:"features"`
	Description  string         `gorm:"type:text;not null" json:"description"`
	Duration     string         `gorm:"type:varchar(50);not null" json:"duration"`
	DurationDays int64          `gorm:"type:bigint;not null" json:"duration_days"`
}

func (PlanEntity) TableName() string {
	return PlanTableName
}
