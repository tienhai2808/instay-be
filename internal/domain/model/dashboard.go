package model

import "time"

type DashboardGuestNationalPoint struct {
	National   string
	GuestCount int64
}

type DashboardGuestNational struct {
	Chart []*DashboardGuestNationalPoint
}

type DashboardSupportTicketTime struct {
	UnderOrEqual5Minutes     int64
	UnderOrEqual10Minutes    int64
	UnderOrEqual15Minutes    int64
	UnderOrEqual30Minutes    int64
	Over30Minutes            int64
	AverageResolutionSeconds float64
}

type DashboardServiceTypeCountPoint struct {
	ServiceTypeID   int64
	ServiceTypeName string
	OrderCount      int64
	OrderPercent    float64
}

type DashboardServiceTypeCount struct {
	Chart []*DashboardServiceTypeCountPoint
}

type DashboardReviewSummary struct {
	ReviewCount   int64
	AverageStar   float64
	PositiveCount int64
	NegativeCount int64
	NeutralCount  int64
}

type DashboardReviewCategoryCountPoint struct {
	Category    string
	ReviewCount int64
}

type DashboardReviewCategoryCount struct {
	Chart []*DashboardReviewCategoryCountPoint
}

type DashboardReviewCategorySentimentPoint struct {
	Category      string
	PositiveCount int64
	NegativeCount int64
	NeutralCount  int64
}

type DashboardReviewCategorySentiment struct {
	Chart []*DashboardReviewCategorySentimentPoint
}

type DashboardReviewSentimentCount struct {
	PositiveCount int64
	NegativeCount int64
	NeutralCount  int64
}

type DashboardReviewSentimentTimePoint struct {
	Date          time.Time
	PositiveCount int64
	NegativeCount int64
	NeutralCount  int64
}

type DashboardReviewSentimentTime struct {
	Granularity DashboardRevenueGranularity
	FromDate    time.Time
	ToDate      time.Time
	Chart       []*DashboardReviewSentimentTimePoint
}
