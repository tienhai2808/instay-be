package response

import "time"

type EventDetailsResponse struct {
	ID          int64                  `json:"id,string"`
	Title       string                 `json:"title"`
	Excerpt     string                 `json:"excerpt"`
	Content     string                 `json:"content"`
	Thumbnail   *string                `json:"thumbnail"`
	IsActive    bool                   `json:"is_active"`
	PublishedAt *time.Time             `json:"published_at"`
	Branches    []*BasicBranchResponse `json:"branches"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
	CreatedBy   *BasicUserResponse     `json:"created_by"`
	UpdatedBy   *BasicUserResponse     `json:"updated_by"`
}

type EventResponse struct {
	ID          int64                  `json:"id,string"`
	Title       string                 `json:"title"`
	Content     string                 `json:"content"`
	Thumbnail   *string                `json:"thumbnail"`
	PublishedAt *time.Time             `json:"published_at"`
	Branches    []*BasicBranchResponse `json:"branches"`
}

type BasicEventResponse struct {
	ID          int64      `json:"id,string"`
	Slug        string     `json:"slug"`
	Title       string     `json:"title"`
	Excerpt     string     `json:"excerpt"`
	Thumbnail   *string    `json:"thumbnail"`
	PublishedAt *time.Time `json:"published_at"`
}

type SimpleEventResponse struct {
	ID          int64                  `json:"id,string"`
	Title       string                 `json:"title"`
	Thumbnail   *string                `json:"thumbnail"`
	IsActive    bool                   `json:"is_active"`
	PublishedAt *time.Time             `json:"published_at"`
	Branches    []*BasicBranchResponse `json:"branches"`
	CreatedAt   time.Time              `json:"created_at"`
}
