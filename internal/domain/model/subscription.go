package model

import "time"

type Subscription struct {
	Base

	TenantID int64
	PlanID   int64
	StartAt  time.Time
	EndAt    time.Time
	IsUsing  bool

	Tenant *Tenant
	Plan   *Plan
}
