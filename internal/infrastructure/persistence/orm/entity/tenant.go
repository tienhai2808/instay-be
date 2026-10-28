package entity

import "gorm.io/datatypes"

const TenantTableName = "tenants"

type TenantEntity struct {
	BaseEntity

	Name       string         `gorm:"type:varchar(150);not null" json:"name"`
	Logo       *string        `gorm:"type:varchar(150);uniqueIndex:tenants_logo_key" json:"logo"`
	Phone      string         `gorm:"type:varchar(20);not null" json:"phone"`
	City       string         `gorm:"type:varchar(100);not null" json:"city"`
	Domain     string         `gorm:"type:varchar(255);not null;uniqueIndex:tenants_domain_key" json:"domain"`
	SchemaName string         `gorm:"type:varchar(50);not null;uniqueIndex:tenants_schema_name_key" json:"schema_name"`
	Settings   datatypes.JSON `gorm:"type:jsonb" json:"settings"`
	IsActive   bool           `gorm:"type:boolean;not null" json:"is_active"`
}

func (TenantEntity) TableName() string {
	return TenantTableName
}
