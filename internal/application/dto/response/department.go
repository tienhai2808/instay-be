package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type DepartmentResponse struct {
	ID          int64                 `json:"id,string"`
	Name        string                `json:"name"`
	Phone       string                `json:"phone"`
	Description string                `json:"description"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   time.Time             `json:"updated_at"`
	CreatedBy   *BasicUserResponse    `json:"created_by"`
	UpdatedBy   *BasicUserResponse    `json:"updated_by"`
	Permissions []*PermissionResponse `json:"permissions"`
	UsersCount  int64                 `json:"users_count"`
}

type BasicDepartmentResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type SimpleDepartmentResponse struct {
	ID          int64                 `json:"id,string"`
	Name        string                `json:"name"`
	Permissions []*PermissionResponse `json:"permissions"`
}

type PermissionResponse struct {
	ID            int64                `json:"id,string"`
	Name          model.PermissionName `json:"name"`
	DisplayNameVI string               `json:"display_name_vi"`
	DisplayNameEN string               `json:"display_name_en"`
}
