package request

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

type CustomerPaginationQuery struct {
	Page   uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit  uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort   string `form:"sort" json:"sort"`
	Order  string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Search string `form:"search" json:"search"`
}

type CreateCustomerRequest struct {
	BirthYear uint32               `json:"birth_year" binding:"required,gte=1900,lte=2100"`
	Gender    model.CustomerGender `json:"gender" binding:"required,oneof=male female"`
	Phone     string               `json:"phone" binding:"required,min=10,max=20"`
}

type ExportListCustomersRequest struct {
	Sort  *string `binding:"omitempty,oneof=created_at birth_year gender phone" json:"sort"`
	Order *string `binding:"omitempty,oneof=asc desc" json:"order"`
}
