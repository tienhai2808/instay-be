package entity

const BranchTableName = "branches"

type BranchEntity struct {
	BaseEntity
	AuditEntity

	Name     string `gorm:"type:varchar(150);not null;uniqueIndex:branches_name_key" json:"name"`
	Address  string `gorm:"type:varchar(255);not null" json:"address"`
	Ward     string `gorm:"type:varchar(150);not null" json:"ward"`
	City     string `gorm:"type:varchar(150);not null" json:"city"`
	IsActive bool   `gorm:"type:boolean;not null" json:"is_active"`

	RoomTypes []*RoomTypeEntity `gorm:"many2many:branch_room_types;foreignKey:ID;joinForeignKey:BranchID;references:ID;joinReferences:RoomTypeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_types"`
	Services  []*ServiceEntity  `gorm:"many2many:branch_services;foreignKey:ID;joinForeignKey:BranchID;references:ID;joinReferences:ServiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"services"`
	Events    []*EventEntity    `gorm:"many2many:branch_events;foreignKey:ID;joinForeignKey:BranchID;references:ID;joinReferences:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"events"`
	Users     []*UserEntity     `gorm:"many2many:user_branches;foreignKey:ID;joinForeignKey:BranchID;references:ID;joinReferences:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"users"`
}

func (BranchEntity) TableName() string {
	return BranchTableName
}
