package entity

const SupportCategoryTableName = "support_categories"

type SupportCategoryEntity struct {
	BaseEntity
	AuditEntity

	Name         string `gorm:"type:varchar(150);not null;uniqueIndex:support_categories_name_key" json:"name"`
	DepartmentID int64  `gorm:"type:bigint;not null" json:"department_id"`

	Department          *DepartmentEntity      `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_support_categories_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"department"`
	SupportTickets      []*SupportTicketEntity `gorm:"foreignKey:SupportCategoryID;references:ID;constraint:fk_support_tickets_support_category,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"support_tickets"`
	SupportTicketsCount int64                  `gorm:"-" json:"support_tickets_count"`
}

func (SupportCategoryEntity) TableName() string {
	return SupportCategoryTableName
}
