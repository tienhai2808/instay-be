package entity

const BranchServiceTableName = "branch_services"

type BranchServiceEntity struct {
	BranchID  int64 `gorm:"type:bigint;primaryKey" json:"branch_id"`
	ServiceID int64 `gorm:"type:bigint;primaryKey" json:"service_id"`

	Branch  *BranchEntity  `gorm:"foreignKey:BranchID;references:ID;constraint:fk_branch_services_branch,OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch"`
	Service *ServiceEntity `gorm:"foreignKey:ServiceID;references:ID;constraint:fk_branch_services_service,OnUpdate:CASCADE,OnDelete:CASCADE" json:"service"`
}

func (BranchServiceEntity) TableName() string {
	return BranchServiceTableName
}
