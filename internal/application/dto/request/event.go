package request

type CreateEventRequest struct {
	Title            string     `json:"title" binding:"required,min=2"`
	Excerpt          string     `json:"excerpt" binding:"required"`
	Content          string     `json:"content" binding:"required"`
	Thumbnail        *string    `json:"thumbnail" binding:"required,min=2"`
	IsActive         *bool      `json:"is_active" binding:"required"`
	SendNotification *bool      `json:"send_notification" binding:"required"`
	BranchIDs        Int64Slice `json:"branch_ids" binding:"omitempty,dive,required"`
}

type UpdateEventRequest struct {
	Title            string     `json:"title" binding:"required,min=2"`
	Excerpt          string     `json:"excerpt" binding:"required"`
	Content          string     `json:"content" binding:"required"`
	Thumbnail        *string    `json:"thumbnail" binding:"required,min=2"`
	IsActive         *bool      `json:"is_active" binding:"required"`
	SendNotification *bool      `json:"send_notification" binding:"omitempty"`
	BranchIDs        Int64Slice `json:"branch_ids" binding:"omitempty,dive,required"`
}

type EventPaginationQuery struct {
	Page     uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort     string `form:"sort" json:"sort"`
	Order    string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	BranchID int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
	IsActive *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
	Search   string `form:"search" json:"search"`
}
