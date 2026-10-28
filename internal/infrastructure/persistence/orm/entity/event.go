package entity

import "time"

const EventTableName = "events"

type EventEntity struct {
	BaseEntity
	AuditEntity

	Title              string     `gorm:"type:varchar(255);not null" json:"title"`
	Slug               string     `gorm:"type:varchar(255);not null;uniqueIndex:events_slug_key" json:"slug"`
	IsActive           bool       `gorm:"type:boolean;not null" json:"is_active"`
	Thumbnail          *string    `gorm:"type:varchar(150);uniqueIndex:events_thumbnail_key" json:"thumbnail"`
	Content            string     `gorm:"type:text;not null" json:"content"`
	Excerpt            string     `gorm:"type:text;not null" json:"excerpt"`
	PublishedAt        *time.Time `json:"published_at"`
	IsNotificationSent bool       `gorm:"type:boolean;not null" json:"is_notification_sent"`

	Branches []*BranchEntity `gorm:"many2many:branch_events;foreignKey:ID;joinForeignKey:EventID;references:ID;joinReferences:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branches"`
}

func (EventEntity) TableName() string {
	return EventTableName
}
