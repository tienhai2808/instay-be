package entity

const DepartmentPermissionTableName = "department_permissions"

type DepartmentPermissionEntity struct {
	DepartmentID int64 `gorm:"type:bigint;not null" json:"department_id"`
	PermissionID int64 `gorm:"type:bigint;not null" json:"permission_id"`
}

func (DepartmentPermissionEntity) TableName() string {
	return DepartmentPermissionTableName
}
