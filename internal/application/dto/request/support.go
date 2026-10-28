package request

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type SupportCategoryPaginationQuery struct {
	Page         uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit        uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search       string `form:"search" json:"search"`
	DepartmentID int64  `form:"department_id" binding:"omitempty" json:"department_id,string"`
}

type CreateSupportCategoryRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	DepartmentID int64  `json:"department_id,string" binding:"required"`
}

type UpdateSupportCategoryRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	DepartmentID int64  `json:"department_id,string" binding:"required"`
}

type CreateSupportTicketRequest struct {
	SupportCategoryID int64                       `json:"support_category_id,string" binding:"required"`
	Content           string                      `json:"content" binding:"required,min=1"`
	Priority          model.SupportTicketPriority `json:"priority" binding:"required,oneof=low medium high urgent"`
}

type UpdateSupportTicketStatusForGuestRequest struct {
	Status model.SupportTicketStatus `json:"status" binding:"oneof=cancelled"`
}

type UpdateSupportTicketStatusForAdminRequest struct {
	Status model.SupportTicketStatus `json:"status" binding:"oneof=accepted done"`
}

type SupportTicketPaginationQuery struct {
	Page     uint32                      `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32                      `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort     string                      `form:"sort" json:"sort"`
	Order    string                      `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	FromDate time.Time                   `form:"from_date" time_format:"2006-01-02"`
	ToDate   time.Time                   `form:"to_date" time_format:"2006-01-02"`
	Status   model.SupportTicketStatus   `form:"status" binding:"omitempty,oneof=pending accepted cancelled done"`
	Priority model.SupportTicketPriority `form:"priority" binding:"omitempty,oneof=low medium high urgent"`
}
