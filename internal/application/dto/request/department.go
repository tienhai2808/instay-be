package request

type CreateDepartmentRequest struct {
	Name          string     `json:"name" binding:"required,min=2"`
	Phone         string     `json:"phone" binding:"required,max=20"`
	Description   string     `json:"description" binding:"required,min=1"`
	PermissionIDs Int64Slice `json:"permission_ids" binding:"required,min=1,dive,required"`
}

type UpdateDepartmentRequest struct {
	Name          string     `json:"name" binding:"required,min=2"`
	Phone         string     `json:"phone" binding:"required,max=20"`
	Description   string     `json:"description" binding:"required,min=1"`
	PermissionIDs Int64Slice `json:"permission_ids" binding:"required,min=1,dive,required"`
}

type DepartmentPaginationQuery struct {
	Page   uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit  uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search string `form:"search" json:"search"`
}
