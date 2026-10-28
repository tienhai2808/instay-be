package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBasicDepartmentResponse(dept *model.Department) *response.BasicDepartmentResponse {
	if dept == nil {
		return nil
	}

	return &response.BasicDepartmentResponse{
		ID:   dept.ID,
		Name: dept.Name,
	}
}

func ToDepartmentResponse(dept *model.Department) *response.DepartmentResponse {
	if dept == nil {
		return nil
	}

	return &response.DepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Phone:       dept.Phone,
		Description: dept.Description,
		CreatedAt:   dept.CreatedAt,
		UpdatedAt:   dept.UpdatedAt,
		CreatedBy:   ToBasicUserResponse(dept.CreatedBy),
		UpdatedBy:   ToBasicUserResponse(dept.UpdatedBy),
		Permissions: ToPermissionsResponse(dept.Permissions),
		UsersCount:  dept.UsersCount,
	}
}

func ToDepartmentsResponse(depts []*model.Department) []*response.DepartmentResponse {
	if len(depts) == 0 {
		return make([]*response.DepartmentResponse, 0)
	}

	deptsRes := make([]*response.DepartmentResponse, 0, len(depts))
	for _, dept := range depts {
		deptsRes = append(deptsRes, ToDepartmentResponse(dept))
	}

	return deptsRes
}

func ToBasicDepartmentsResponse(depts []*model.Department) []*response.BasicDepartmentResponse {
	if len(depts) == 0 {
		return make([]*response.BasicDepartmentResponse, 0)
	}

	deptsRes := make([]*response.BasicDepartmentResponse, 0, len(depts))
	for _, dept := range depts {
		deptsRes = append(deptsRes, ToBasicDepartmentResponse(dept))
	}

	return deptsRes
}

func ToSimpleDepartmentResponse(dept *model.Department) *response.SimpleDepartmentResponse {
	if dept == nil {
		return nil
	}

	return &response.SimpleDepartmentResponse{
		ID:          dept.ID,
		Name:        dept.Name,
		Permissions: ToPermissionsResponse(dept.Permissions),
	}
}

func ToPermissionResponse(perm *model.Permission) *response.PermissionResponse {
	if perm == nil {
		return nil
	}

	return &response.PermissionResponse{
		ID:            perm.ID,
		Name:          perm.Name,
		DisplayNameVI: perm.DisplayNameVI,
		DisplayNameEN: perm.DisplayNameEN,
	}
}

func ToPermissionsResponse(perms []*model.Permission) []*response.PermissionResponse {
	if len(perms) == 0 {
		return make([]*response.PermissionResponse, 0)
	}

	permsRes := make([]*response.PermissionResponse, 0, len(perms))
	for _, perm := range perms {
		permsRes = append(permsRes, ToPermissionResponse(perm))
	}

	return permsRes
}

func ToDepartmentQuery(query request.DepartmentPaginationQuery) model.DepartmentQuery {
	return model.DepartmentQuery{
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
	}
}
