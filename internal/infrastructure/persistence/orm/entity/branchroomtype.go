package entity

const BranchRoomTypeTableName = "branch_room_types"

type BranchRoomTypeEntity struct {
	BranchID   int64 `gorm:"type:bigint;primaryKey" json:"branch_id"`
	RoomTypeID int64 `gorm:"type:bigint;primaryKey" json:"room_type_id"`

	Branch   *BranchEntity   `gorm:"foreignKey:BranchID;references:ID;constraint:fk_branch_room_types_branch,OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch"`
	RoomType *RoomTypeEntity `gorm:"foreignKey:RoomTypeID;references:ID;constraint:fk_branch_room_types_room_type,OnUpdate:CASCADE,OnDelete:CASCADE" json:"room_type"`
}

func (BranchRoomTypeEntity) TableName() string {
	return BranchRoomTypeTableName
}
