package entity

const UserNotificationTableName = "user_notifications"

type UserNotificationEntity struct {
	NotificationID int64 `gorm:"type:bigint;not null" json:"notification_id"`
	UserID         int64 `gorm:"type:bigint;not null" json:"user_id"`
}

func (UserNotificationEntity) TableName() string {
	return UserNotificationTableName
}
