package entity

const BranchEventTableName = "branch_events"

type BranchEventEntity struct {
	BranchID int64 `gorm:"type:bigint;primaryKey" json:"branch_id"`
	EventID  int64 `gorm:"type:bigint;primaryKey" json:"event_id"`

	Branch *BranchEntity `gorm:"foreignKey:BranchID;references:ID;constraint:fk_branch_events_branch,OnUpdate:CASCADE,OnDelete:CASCADE" json:"branch"`
	Event  *EventEntity  `gorm:"foreignKey:EventID;references:ID;constraint:fk_branch_events_event,OnUpdate:CASCADE,OnDelete:CASCADE" json:"event"`
}

func (BranchEventEntity) TableName() string {
	return BranchEventTableName
}
