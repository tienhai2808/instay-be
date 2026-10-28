package request

import "time"

type DashboardOverviewQuery struct {
	FromDate *time.Time `form:"from_date" time_format:"2006-01-02" json:"from_date" binding:"omitempty"`
	ToDate   *time.Time `form:"to_date" time_format:"2006-01-02" json:"to_date" binding:"omitempty"`
}

type DashboardRevenueQuery struct {
	RoomTypeID int64 `form:"room_type_id" binding:"omitempty" json:"room_type_id,string"`
}
