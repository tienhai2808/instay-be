package entity

const FloorTableName = "floors"

type FloorEntity struct {
	BaseEntity
	AuditEntity

	Name     string `gorm:"type:varchar(50);not null;uniqueIndex:floors_branch_id_name_key,priority:2" json:"name"`
	BranchID int64  `gorm:"type:bigint;uniqueIndex:floors_branch_id_name_key,priority:1" json:"branch_id"`

	Branch *BranchEntity `gorm:"foreignKey:BranchID;references:ID;constraint:fk_floors_branch,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"branch"`
	Rooms  []*RoomEntity `gorm:"foreignKey:FloorID;references:ID;constraint:fk_rooms_floor,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"rooms"`
}

func (FloorEntity) TableName() string {
	return FloorTableName
}
