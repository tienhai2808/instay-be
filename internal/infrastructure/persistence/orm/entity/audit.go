package entity

type AuditEntity struct {
	CreatedByID *int64 `gorm:"type:bigint" json:"created_by_id"`
	UpdatedByID *int64 `gorm:"type:bigint" json:"updated_by_id"`

	CreatedBy *UserEntity `gorm:"belongsTo;foreignKey:CreatedByID;references:ID;constraint:-" json:"created_by"`
	UpdatedBy *UserEntity `gorm:"belongsTo;foreignKey:UpdatedByID;references:ID;constraint:-" json:"updated_by"`
}
