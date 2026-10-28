package model

import "time"

type Event struct {
	Base
	Audit

	Thumbnail          *string
	Title              string
	Slug               string
	Content            string
	Excerpt            string
	IsActive           bool
	PublishedAt        *time.Time
	IsNotificationSent bool

	Branches []*Branch
}

type EventQuery struct {
	Page          uint32
	Limit         uint32
	Sort          string
	Order         string
	BranchID      int64
	BranchIDs     []int64
	IncludeGlobal bool
	IsActive      *bool
	Search        string
}

func NewEvent(
	id int64,
	thumbnail *string,
	title, slug, content, excerpt string,
	isActive bool,
	publishedAt *time.Time,
	isNotificationSent bool,
	createdByID, updatedByID *int64,
) *Event {
	return &Event{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Thumbnail:          thumbnail,
		Title:              title,
		Slug:               slug,
		Content:            content,
		Excerpt:            excerpt,
		IsActive:           isActive,
		PublishedAt:        publishedAt,
		IsNotificationSent: isNotificationSent,
	}
}
