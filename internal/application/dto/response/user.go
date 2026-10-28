package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type UserRoleResponse struct {
	Name          string `json:"name"`
	DisplayNameVI string `json:"display_name_vi"`
	DisplayNameEN string `json:"display_name_en"`
}

type BasicUserResponse struct {
	ID        int64  `json:"id,string"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type UserResponse struct {
	ID         int64                     `json:"id,string"`
	Username   string                    `json:"username"`
	Email      string                    `json:"email"`
	Phone      string                    `json:"phone"`
	Role       model.UserRole            `json:"role"`
	Avatar     *string                   `json:"avatar"`
	IsActive   bool                      `json:"is_active"`
	FirstName  string                    `json:"first_name"`
	LastName   string                    `json:"last_name"`
	CreatedAt  time.Time                 `json:"created_at"`
	Department *SimpleDepartmentResponse `json:"department"`
	Branches   []*BasicBranchResponse    `json:"branches"`
}

type SimpleUserResponse struct {
	ID         int64                    `json:"id,string"`
	FirstName  string                   `json:"first_name"`
	LastName   string                   `json:"last_name"`
	Role       model.UserRole           `json:"role"`
	IsActive   bool                     `json:"is_active"`
	CreatedAt  time.Time                `json:"created_at"`
	Department *BasicDepartmentResponse `json:"department"`
	Branches   []*BasicBranchResponse   `json:"branches"`
}

type UserDetailsResponse struct {
	ID         int64                    `json:"id,string"`
	Username   string                   `json:"username"`
	Email      string                   `json:"email"`
	Phone      string                   `json:"phone"`
	Role       model.UserRole           `json:"role"`
	Avatar     *string                  `json:"avatar"`
	IsActive   bool                     `json:"is_active"`
	FirstName  string                   `json:"first_name"`
	LastName   string                   `json:"last_name"`
	CreatedAt  time.Time                `json:"created_at"`
	UpdatedAt  time.Time                `json:"updated_at"`
	Department *BasicDepartmentResponse `json:"department"`
	Branches   []*BasicBranchResponse   `json:"branches"`
	CreatedBy  *BasicUserResponse       `json:"created_by"`
	UpdatedBy  *BasicUserResponse       `json:"updated_by"`
}
