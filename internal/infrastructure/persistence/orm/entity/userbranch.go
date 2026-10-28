package entity

const UserBranchTableName = "user_branches"

type UserBranchEntity struct {
	UserID   int64 `gorm:"type:bigint;primaryKey" json:"user_id"`
	BranchID int64 `gorm:"type:bigint;primaryKey" json:"branch_id"`

	User   *UserEntity   `gorm:"foreignKey:UserID;references:ID;constraint:fk_user_branches_user,OnUpdate:CASCADE,OnDelete:CASCADE" json:"user"`
	Branch *BranchEntity `gorm:"foreignKey:BranchID;references:ID;constraint:fk_user_branches_branch,OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch"`
}

func (UserBranchEntity) TableName() string {
	return UserBranchTableName
}
