package entity

import "time"

const ServiceTableName = "services"

type ServiceEntity struct {
	BaseEntity
	AuditEntity

	Name            string     `gorm:"type:varchar(255);not null" json:"name"`
	Slug            string     `gorm:"type:varchar(255);not null;uniqueIndex:services_slug_key" json:"slug"`
	Price           float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	IsActive        bool       `gorm:"type:boolean;not null" json:"is_active"`
	Description     string     `gorm:"type:text;not null" json:"description"`
	DiscountType    *string    `gorm:"type:varchar(20);check:discount_type IN ('percent', 'amount')" json:"discount_type"`
	DiscountValue   *float64   `gorm:"type:decimal(10,2)" json:"discount_value"`
	DiscountStartAt *time.Time `gorm:"type:date" json:"discount_start_at"`
	DiscountEndAt   *time.Time `gorm:"type:date" json:"discount_end_at"`
	ServiceTypeID   int64      `gorm:"type:bigint;not null" json:"service_type_id"`

	ServiceType   *ServiceTypeEntity    `gorm:"foreignKey:ServiceTypeID;references:ID;constraint:fk_services_service_type,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"service_type"`
	ServiceImages []*ServiceImageEntity `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_service_images_service,OnUpdate:CASCADE,OnDelete:CASCADE" json:"service_images"`
	Branches      []*BranchEntity       `gorm:"many2many:branch_services;foreignKey:ID;joinForeignKey:ServiceID;references:ID;joinReferences:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branches"`
}

func (ServiceEntity) TableName() string {
	return ServiceTableName
}
