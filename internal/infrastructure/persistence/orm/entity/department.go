package entity

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

const DepartmentTableName = "departments"

type DepartmentEntity struct {
	BaseEntity
	AuditEntity

	Name        string `gorm:"type:varchar(150);not null;uniqueIndex:departments_name_key" json:"name"`
	Phone       string `gorm:"type:varchar(20);not null;uniqueIndex:departments_phone_key" json:"phone"`
	Description string `gorm:"type:text;not null" json:"description"`

	Users         []*UserEntity         `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_users_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"users"`
	ServiceTypes  []*ServiceTypeEntity  `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_service_types_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service_types"`
	Notifications []*NotificationEntity `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_notifications_department,OnUpdate:CASCADE,OnDelete:CASCADE" json:"notifications"`
	Chats         []*ChatEntity         `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_chats_department,OnUpdate:CASCADE,OnDelete:CASCADE" json:"chats"`
	Permissions   []*PermissionEntity   `gorm:"many2many:department_permissions;joinForeignKey:DepartmentID;joinReferences:PermissionID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"permissions"`
	UsersCount    int64                 `gorm:"-" json:"users_count"`
}

func (DepartmentEntity) TableName() string {
	return DepartmentTableName
}

func DefaultDepartments() []DepartmentEntity {
	return []DepartmentEntity{
		{
			Name:        model.DepartmentCustomerService,
			Phone:       "012345678",
			Description: "Bộ phận chăm sóc khách hàng",
		},
	}
}
