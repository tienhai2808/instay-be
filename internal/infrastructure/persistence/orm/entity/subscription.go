package entity

import "time"

const SubscriptionTableName = "subscriptions"

type SubscriptionEntity struct {
	BaseEntity

	TenantID int64     `gorm:"type:bigint;not null" json:"tenant_id"`
	PlanID   int64     `gorm:"type:bigint;not null" json:"plan_id"`
	StartAt  time.Time `gorm:"not null" json:"start_at"`
	EndAt    time.Time `gorm:"not null" json:"end_at"`
	IsUsing  bool      `gorm:"type:boolean;not null" json:"is_using"`

	Tenant *TenantEntity `gorm:"foreignKey:TenantID;references:ID;constraint:fk_subscriptions_tenant,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"tenant"`
	Plan   *PlanEntity   `gorm:"foreignKey:PlanID;references:ID;constraint:fk_subscriptions_plan,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"plan"`
}

func (SubscriptionEntity) TableName() string {
	return SubscriptionTableName
}
