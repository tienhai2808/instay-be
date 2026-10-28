package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToDepartmentEntity(dept *model.Department) *entity.DepartmentEntity {
	if dept == nil {
		return nil
	}

	return &entity.DepartmentEntity{
		BaseEntity: entity.BaseEntity{
			ID:        dept.ID,
			CreatedAt: dept.CreatedAt,
			UpdatedAt: dept.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: dept.CreatedByID,
			UpdatedByID: dept.UpdatedByID,
		},
		Name:        dept.Name,
		Phone:       dept.Phone,
		Description: dept.Description,
		Permissions: ToPermissionEntities(dept.Permissions),
	}
}

func ToPermissionEntities(perms []*model.Permission) []*entity.PermissionEntity {
	if len(perms) == 0 {
		return make([]*entity.PermissionEntity, 0)
	}

	permsRes := make([]*entity.PermissionEntity, 0, len(perms))
	for _, perm := range perms {
		permsRes = append(permsRes, ToPermissionEntity(perm))
	}

	return permsRes
}

func ToPermissionEntity(perm *model.Permission) *entity.PermissionEntity {
	if perm == nil {
		return nil
	}

	return &entity.PermissionEntity{
		ID:            perm.ID,
		Name:          string(perm.Name),
		DisplayNameVI: perm.DisplayNameVI,
		DisplayNameEN: perm.DisplayNameEN,
	}
}

func ToDepartmentPermissionEntity(dp *model.DepartmentPermission) *entity.DepartmentPermissionEntity {
	if dp == nil {
		return nil
	}

	return &entity.DepartmentPermissionEntity{
		DepartmentID: dp.DepartmentID,
		PermissionID: dp.PermissionID,
	}
}

func ToDepartmentPermissionEntities(dps []*model.DepartmentPermission) []*entity.DepartmentPermissionEntity {
	if len(dps) == 0 {
		return make([]*entity.DepartmentPermissionEntity, 0)
	}

	dpsRes := make([]*entity.DepartmentPermissionEntity, 0, len(dps))
	for _, dp := range dps {
		dpsRes = append(dpsRes, ToDepartmentPermissionEntity(dp))
	}

	return dpsRes
}

func ToDepartment(dept *entity.DepartmentEntity) *model.Department {
	if dept == nil {
		return nil
	}

	return &model.Department{
		Base: model.Base{
			ID:        dept.ID,
			CreatedAt: dept.CreatedAt,
			UpdatedAt: dept.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: dept.CreatedByID,
			UpdatedByID: dept.UpdatedByID,
			CreatedBy:   ToUser(dept.CreatedBy),
			UpdatedBy:   ToUser(dept.UpdatedBy),
		},
		Name:        dept.Name,
		Phone:       dept.Phone,
		Description: dept.Description,
		Permissions: ToPermissions(dept.Permissions),
		Users:       ToUsers(dept.Users),
	}
}

func ToPermission(perm *entity.PermissionEntity) *model.Permission {
	if perm == nil {
		return nil
	}

	return &model.Permission{
		ID:            perm.ID,
		Name:          model.PermissionName(perm.Name),
		DisplayNameVI: perm.DisplayNameVI,
		DisplayNameEN: perm.DisplayNameEN,
		Departments:   ToDepartments(perm.Departments),
	}
}

func ToPermissions(perms []*entity.PermissionEntity) []*model.Permission {
	if len(perms) == 0 {
		return make([]*model.Permission, 0)
	}

	permsRes := make([]*model.Permission, 0, len(perms))
	for _, perm := range perms {
		permsRes = append(permsRes, ToPermission(perm))
	}

	return permsRes
}

func ToDepartments(depts []*entity.DepartmentEntity) []*model.Department {
	if len(depts) == 0 {
		return make([]*model.Department, 0)
	}

	deptsRes := make([]*model.Department, 0, len(depts))
	for _, dept := range depts {
		deptsRes = append(deptsRes, ToDepartment(dept))
	}

	return deptsRes
}
