package entity

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

const PermissionTableName = "permissions"

type PermissionEntity struct {
	ID            int64  `gorm:"type:bigint;primaryKey" json:"id"`
	Name          string `gorm:"type:varchar(50);not null;uniqueIndex:permissions_name_key" json:"name"`
	DisplayNameVI string `gorm:"type:varchar(150);not null" json:"display_name_vi"`
	DisplayNameEN string `gorm:"type:varchar(150);not null" json:"display_name_en"`

	Departments []*DepartmentEntity `gorm:"many2many:department_permissions;joinForeignKey:PermissionID;joinReferences:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"departments"`
}

func (PermissionEntity) TableName() string {
	return PermissionTableName
}

func DefaultPermissions() []PermissionEntity {
	return []PermissionEntity{
		{
			Name:          string(model.PermBookingManagement),
			DisplayNameVI: string(model.PermBookingManagementDisplayNameVI),
			DisplayNameEN: string(model.PermBookingManagementDisplayNameEN),
		},
		{
			Name:          string(model.PermOrderServiceManagement),
			DisplayNameVI: string(model.PermOrderServiceManagementDisplayNameVI),
			DisplayNameEN: string(model.PermOrderServiceManagementDisplayNameEN),
		},
		{
			Name:          string(model.PermSupportTicketManagement),
			DisplayNameVI: string(model.PermSupportTicketManagementDisplayNameVI),
			DisplayNameEN: string(model.PermSupportTicketManagementDisplayNameEN),
		},
		{
			Name:          string(model.PermGuestManagement),
			DisplayNameVI: string(model.PermGuestManagementDisplayNameVI),
			DisplayNameEN: string(model.PermGuestManagementDisplayNameEN),
		},
		{
			Name:          string(model.PermEventManagement),
			DisplayNameVI: string(model.PermEventManagementDisplayNameVI),
			DisplayNameEN: string(model.PermEventManagementDisplayNameEN),
		},
	}
}
