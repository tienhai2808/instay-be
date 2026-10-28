package response

import "time"

type DashboardHotelOperationResponse struct {
	TodayCheckIn      int64 `json:"today_check_in"`
	TodayCheckOut     int64 `json:"today_check_out"`
	PendingRequest    int64 `json:"pending_request"`
	RequestInProgress int64 `json:"request_in_progress"`
	UnreadMessage     int64 `json:"unread_message"`
}

type DashboardBookingOverviewResponse struct {
	TotalRevenue                  float64 `json:"total_revenue"`
	TotalRevenueChangePercent     float32 `json:"total_revenue_change_percent"`
	OccupancyRate                 float32 `json:"occupancy_rate"`
	OccupancyRateChangePercent    float32 `json:"occupancy_rate_change_percent"`
	CancellationRate              float32 `json:"cancellation_rate"`
	CancellationRateChangePercent float32 `json:"cancellation_rate_change_percent"`
	AverageDaysStay               float32 `json:"average_days_stay"`
	AverageDaysStayChangePercent  float32 `json:"average_days_stay_change_percent"`
}

type DashboardBookingRevenuePointResponse struct {
	Date    time.Time `json:"date"`
	Revenue float64   `json:"revenue"`
}

type DashboardBookingRevenueResponse struct {
	Granularity string                                  `json:"granularity"`
	FromDate    time.Time                               `json:"from_date"`
	ToDate      time.Time                               `json:"to_date"`
	Chart       []*DashboardBookingRevenuePointResponse `json:"chart"`
}

type DashboardRoomTypeRevenuePointResponse struct {
	RoomTypeID   int64   `json:"room_type_id,string"`
	RoomTypeName string  `json:"room_type_name"`
	Revenue      float64 `json:"revenue"`
}

type DashboardRoomTypeRevenueResponse struct {
	Chart []*DashboardRoomTypeRevenuePointResponse `json:"chart"`
}

type DashboardSourceRevenuePointResponse struct {
	SourceID       int64   `json:"source_id,string"`
	SourceName     string  `json:"source_name"`
	Revenue        float64 `json:"revenue"`
	RevenuePercent float32 `json:"revenue_percent"`
}

type DashboardSourceRevenueResponse struct {
	Chart []*DashboardSourceRevenuePointResponse `json:"chart"`
}

type DashboardGuestNationalPointResponse struct {
	National   string `json:"national"`
	GuestCount int64  `json:"guest_count"`
}

type DashboardGuestNationalResponse struct {
	Chart []*DashboardGuestNationalPointResponse `json:"chart"`
}

type DashboardSupportTicketTimePointResponse struct {
	Tag           string  `json:"tag"`
	TicketCount   int64   `json:"ticket_count"`
	TicketPercent float32 `json:"ticket_percent"`
}

type DashboardSupportTicketTimeResponse struct {
	AverageResolutionMinutes float64                                    `json:"average_resolution_minutes"`
	Chart                    []*DashboardSupportTicketTimePointResponse `json:"chart"`
}

type DashboardServiceTypeCountPointResponse struct {
	ServiceTypeID   int64   `json:"service_type_id,string"`
	ServiceTypeName string  `json:"service_type_name"`
	OrderCount      int64   `json:"order_count"`
	OrderPercent    float32 `json:"order_percent"`
}

type DashboardServiceTypeCountResponse struct {
	Chart []*DashboardServiceTypeCountPointResponse `json:"chart"`
}

type DashboardReviewSummaryResponse struct {
	ReviewCount  int64   `json:"review_count"`
	AverageStar  float32 `json:"average_star"`
	PositiveRate float32 `json:"positive_rate"`
	NegativeRate float32 `json:"negative_rate"`
}

type DashboardReviewCategoryCountPointResponse struct {
	Category    string `json:"category"`
	ReviewCount int64  `json:"review_count"`
}

type DashboardReviewCategoryCountResponse struct {
	Chart []*DashboardReviewCategoryCountPointResponse `json:"chart"`
}

type DashboardReviewCategorySentimentBreakdownResponse struct {
	Count   int64   `json:"count"`
	Percent float32 `json:"percent"`
}

type DashboardReviewCategorySentimentPointResponse struct {
	Category string                                             `json:"category"`
	Positive *DashboardReviewCategorySentimentBreakdownResponse `json:"positive"`
	Negative *DashboardReviewCategorySentimentBreakdownResponse `json:"negative"`
	Neutral  *DashboardReviewCategorySentimentBreakdownResponse `json:"neutral"`
}

type DashboardReviewCategorySentimentResponse struct {
	Chart []*DashboardReviewCategorySentimentPointResponse `json:"chart"`
}

type DashboardReviewSentimentCountBreakdownResponse struct {
	Count   int64   `json:"count"`
	Percent float32 `json:"percent"`
}

type DashboardReviewSentimentCountResponse struct {
	Positive *DashboardReviewSentimentCountBreakdownResponse `json:"positive"`
	Negative *DashboardReviewSentimentCountBreakdownResponse `json:"negative"`
	Neutral  *DashboardReviewSentimentCountBreakdownResponse `json:"neutral"`
}

type DashboardReviewSentimentTimePointResponse struct {
	Date     time.Time                                       `json:"date"`
	Positive *DashboardReviewSentimentCountBreakdownResponse `json:"positive"`
	Negative *DashboardReviewSentimentCountBreakdownResponse `json:"negative"`
	Neutral  *DashboardReviewSentimentCountBreakdownResponse `json:"neutral"`
}

type DashboardReviewSentimentTimeResponse struct {
	Granularity string                                       `json:"granularity"`
	FromDate    time.Time                                    `json:"from_date"`
	ToDate      time.Time                                    `json:"to_date"`
	Chart       []*DashboardReviewSentimentTimePointResponse `json:"chart"`
}
