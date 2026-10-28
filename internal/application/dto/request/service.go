package request

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type ServicePaginationQuery struct {
	Page          uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit         uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort          string `form:"sort" json:"sort"`
	Order         string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	ServiceTypeID int64  `form:"service_type_id" binding:"omitempty" json:"service_type_id,string"`
	BranchID      int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
	IsActive      *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
	HasDiscount   *bool  `form:"has_discount" binding:"omitempty" json:"has_discount"`
	Search        string `form:"search" json:"search"`
}

type OrderServicePaginationQuery struct {
	Page     uint32                   `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32                   `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort     string                   `form:"sort" json:"sort"`
	Order    string                   `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	FromDate time.Time                `form:"from_date" time_format:"2006-01-02"`
	ToDate   time.Time                `form:"to_date" time_format:"2006-01-02"`
	Status   model.OrderServiceStatus `form:"status" binding:"omitempty,oneof=confirmed pending cancelled"`
}

type CreateOrderServiceForGuestRequest struct {
	ServiceID   int64     `json:"service_id,string" binding:"required"`
	Quantity    uint32    `json:"quantity" binding:"required,min=1"`
	EstimatedAt time.Time `json:"estimated_at" binding:"required"`
	GuestNote   *string   `json:"guest_note" binding:"omitempty,min=1"`
}

type ServiceTypeWithActiveServicesQuery struct {
	ServiceTypeID int64  `form:"service_type_id" binding:"omitempty" json:"service_type_id,string"`
	BranchID      int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
	Search        string `form:"search" json:"search"`
}

type ServiceTypePaginationQuery struct {
	Page         uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit        uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search       string `form:"search" json:"search"`
	DepartmentID int64  `form:"department_id" binding:"omitempty" json:"department_id,string"`
}

type CreateServiceTypeRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	Icon         string `json:"icon" binding:"required,min=1,max=100"`
	DepartmentID int64  `json:"department_id,string" binding:"required"`
}

type UpdateServiceTypeRequest struct {
	Name         string `json:"name" binding:"required,min=2"`
	Icon         string `json:"icon" binding:"required,min=1,max=100"`
	DepartmentID int64  `json:"department_id,string" binding:"required"`
}

type CreateOrderServiceForAdminRequest struct {
	BookingRoomID int64     `json:"booking_room_id,string" binding:"required"`
	ServiceID     int64     `json:"service_id,string" binding:"required"`
	Quantity      uint32    `json:"quantity" binding:"required,min=1"`
	EstimatedAt   time.Time `json:"estimated_at" binding:"required"`
	GuestNote     *string   `json:"guest_note" binding:"omitempty,min=1"`
}

type UpdateOrderServiceStatusForGuestRequest struct {
	Status       model.OrderServiceStatus `json:"status" binding:"oneof=cancelled"`
	CancelReason *string                  `json:"cancel_reason" binding:"omitempty,min=1"`
}

type UpdateOrderServiceStatusForAdminRequest struct {
	Status    model.OrderServiceStatus `json:"status" binding:"oneof=confirmed"`
	StaffNote *string                  `json:"staff_note" binding:"omitempty,min=1"`
}

type CreateServiceRequest struct {
	Name            string                      `json:"name" binding:"required,min=2"`
	ServiceTypeID   int64                       `json:"service_type_id,string" binding:"required"`
	Price           float64                     `json:"price" binding:"required,gt=0"`
	IsActive        *bool                       `json:"is_active" binding:"required"`
	Description     string                      `json:"description" binding:"required"`
	DiscountType    *model.ServiceDiscountType  `json:"discount_type" binding:"omitempty,oneof=percent amount"`
	DiscountValue   *float64                    `json:"discount_value" binding:"omitempty,gt=0"`
	DiscountStartAt *time.Time                  `json:"discount_start_at" binding:"omitempty"`
	DiscountEndAt   *time.Time                  `json:"discount_end_at" binding:"omitempty"`
	BranchIDs       Int64Slice                  `json:"branch_ids" binding:"required,min=1,dive,required"`
	Images          []CreateServiceImageRequest `json:"images" binding:"required,min=1,dive,required"`
}

type CreateServiceImageRequest struct {
	Key         string `json:"key" binding:"required,min=2"`
	IsThumbnail *bool  `json:"is_thumbnail" binding:"required"`
	SortOrder   uint32 `json:"sort_order" binding:"required,gt=0"`
}

type UpdateServiceRequest struct {
	Name            string                      `json:"name" binding:"required,min=2"`
	ServiceTypeID   int64                       `json:"service_type_id,string" binding:"required"`
	Price           float64                     `json:"price" binding:"required,gt=0"`
	IsActive        *bool                       `json:"is_active" binding:"required"`
	Description     string                      `json:"description" binding:"required"`
	DiscountType    *model.ServiceDiscountType  `json:"discount_type" binding:"omitempty,oneof=percent amount"`
	DiscountValue   *float64                    `json:"discount_value" binding:"omitempty,gt=0"`
	DiscountStartAt *time.Time                  `json:"discount_start_at" binding:"omitempty"`
	DiscountEndAt   *time.Time                  `json:"discount_end_at" binding:"omitempty"`
	BranchIDs       Int64Slice                  `json:"branch_ids" binding:"required,min=1,dive,required"`
	Images          []UpdateServiceImageRequest `json:"images" binding:"required,min=1,dive,required"`
}

type UpdateServiceImageRequest struct {
	Key         string `json:"key" binding:"required,min=2"`
	IsThumbnail *bool  `json:"is_thumbnail" binding:"required"`
	SortOrder   uint32 `json:"sort_order" binding:"required,gt=0"`
}
