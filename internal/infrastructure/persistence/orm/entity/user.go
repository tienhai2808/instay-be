package entity

const UserTableName = "users"

type UserEntity struct {
	BaseEntity
	AuditEntity

	Username     string  `gorm:"type:varchar(50);not null;uniqueIndex:users_username_key" json:"username"`
	Email        string  `gorm:"type:varchar(150);not null;uniqueIndex:users_email_key" json:"email"`
	Role         string  `gorm:"type:varchar(20);not null;check:role IN ('staff', 'admin')" json:"role"`
	FirstName    string  `gorm:"type:varchar(150);not null" json:"first_name"`
	LastName     string  `gorm:"type:varchar(150);not null" json:"last_name"`
	Phone        string  `gorm:"type:varchar(10);not null;uniqueIndex:users_phone_key" json:"phone"`
	Password     string  `gorm:"type:varchar(255);not null" json:"password"`
	Avatar       *string `gorm:"type:varchar(150);uniqueIndex:users_avatar_key" json:"avatar"`
	IsActive     bool    `gorm:"type:boolean;not null" json:"is_active"`
	DepartmentID *int64  `gorm:"type:bigint" json:"department_id"`

	Department    *DepartmentEntity     `gorm:"foreignKey:DepartmentID;references:ID;constraint:fk_users_department,OnUpdate:CASCADE,OnDelete:RESTRICT" json:"department"`
	RefreshTokens []*RefreshTokenEntity `gorm:"foreignKey:UserID;references:ID;constraint:fk_refresh_tokens_user,OnUpdate:CASCADE,OnDelete:CASCADE" json:"refresh_tokens"`
	Notifications []*NotificationEntity `gorm:"many2many:user_notifications;joinForeignKey:UserID;joinReferences:NotificationID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"notifications"`
	Branches      []*BranchEntity       `gorm:"many2many:user_branches;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:BranchID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"branches"`
}

func (UserEntity) TableName() string {
	return UserTableName
}
