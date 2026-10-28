package request

type NotificationPaginationQuery struct {
	Page  uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
}
