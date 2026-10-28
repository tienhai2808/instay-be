package entity

const ServiceTypeTableName = "service_types"

type ServiceTypeEntity struct {
	BaseEntity
	AuditEntity

	Name         string `gorm:"type:varchar(150);not null" json:"name"`
	Icon         string `gorm:"type:varchar(100);not null" json:"icon"`
	Slug         string `gorm:"type:varchar(150);not null;uniqueIndex:service_types_slug_key" json:"slug"`
	DepartmentID int64  `gorm:"type:bigint;not null" json:"department_id"`

	Department    *DepartmentEntity `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_service_types_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"department"`
	Services      []*ServiceEntity  `gorm:"foreignKey:ServiceTypeID;references:ID;constraint:fk_services_service_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"services"`
	ServicesCount int64             `gorm:"-" json:"services_count"`
}

func (ServiceTypeEntity) TableName() string {
	return ServiceTypeTableName
}
