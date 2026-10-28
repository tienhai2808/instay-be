package request

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

type CreateUserRequest struct {
	Username     string         `json:"username" binding:"required,min=5"`
	Email        string         `json:"email" binding:"required,email"`
	Phone        string         `json:"phone" binding:"required,len=10"`
	Password     string         `json:"password" binding:"required,min=6"`
	Role         model.UserRole `json:"role" binding:"required,oneof=staff admin"`
	IsActive     *bool          `json:"is_active" binding:"required"`
	FirstName    string         `json:"first_name" binding:"required"`
	LastName     string         `json:"last_name" binding:"required"`
	DepartmentID *int64         `json:"department_id,string" binding:"omitempty"`
	BranchIDs    Int64Slice     `json:"branch_ids" binding:"omitempty,dive,required"`
	Avatar       *string        `json:"avatar" binding:"omitempty"`
}

type UserPaginationQuery struct {
	Page         uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit        uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort         string `form:"sort" json:"sort"`
	Order        string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	Role         string `form:"role" binding:"omitempty,oneof=admin staff" json:"role"`
	DepartmentID int64  `form:"department_id" binding:"omitempty" json:"department_id,string"`
	BranchID     int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
	IsActive     *bool  `form:"is_active" binding:"omitempty" json:"is_active"`
	Search       string `form:"search" json:"search"`
}

type UpdateUserRequest struct {
	Username     string         `json:"username" binding:"required,min=5"`
	Email        string         `json:"email" binding:"required,email"`
	Phone        string         `json:"phone" binding:"required,len=10"`
	FirstName    string         `json:"first_name" binding:"required"`
	LastName     string         `json:"last_name" binding:"required"`
	Avatar       *string        `json:"avatar" binding:"omitempty"`
	Role         model.UserRole `json:"role" binding:"required,oneof=staff admin"`
	IsActive     *bool          `json:"is_active" binding:"required"`
	DepartmentID *int64         `json:"department_id,string" binding:"omitempty"`
	BranchIDs    Int64Slice     `json:"branch_ids" binding:"omitempty,dive,required"`
}

type UpdateUserPasswordRequest struct {
	NewPassword string `json:"new_password" binding:"required,min=6"`
}
