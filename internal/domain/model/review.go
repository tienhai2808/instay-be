package model

import "time"

type ReviewAspect struct {
	Category  string  `json:"category"`
	Sentiment string  `json:"sentiment"`
	Score     float32 `json:"score"`
}

type Review struct {
	Base

	BookingRoomID int64
	Star          uint8
	Content       string
	Aspects       []ReviewAspect

	BookingRoom *BookingRoom
}

type ReviewQuery struct {
	Page      uint32
	Limit     uint32
	Sort      string
	Order     string
	FromDate  *time.Time
	ToDate    *time.Time
	Star      uint8
	BranchIDs []int64
}

type ReviewStarSummary struct {
	All       int64
	FiveStar  int64
	FourStar  int64
	ThreeStar int64
	TwoStar   int64
	OneStar   int64
}

func NewReview(id, bookingID int64, star uint8, content string) *Review {
	return &Review{
		Base: Base{
			ID: id,
		},
		BookingRoomID: bookingID,
		Star:          star,
		Content:       content,
	}
}
