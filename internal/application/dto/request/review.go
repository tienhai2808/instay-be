package request

import "time"

type ReviewPaginationQuery struct {
	Page     uint32     `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32     `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort     string     `form:"sort" json:"sort"`
	Order    string     `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02" json:"from_date" binding:"omitempty"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02" json:"to_date" binding:"omitempty"`
	Star     uint8      `form:"star" binding:"omitempty,min=1,max=5" json:"star"`
}

type CreateReviewRequest struct {
	Star    uint8  `json:"star" binding:"required,min=1,max=5"`
	Content string `json:"content" binding:"required"`
}

type UpdateReviewRequest struct {
	Star    uint8  `json:"star" binding:"required,min=1,max=5"`
	Content string `json:"content" binding:"required"`
}
