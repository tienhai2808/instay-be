package response

import "time"

type SimpleReviewResponse struct {
	ID        int64     `json:"id,string"`
	Star      uint8     `json:"star"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

type ReviewSummaryResponse struct {
	All              int64   `json:"all"`
	Average          float32 `json:"average"`
	FiveStarPercent  float32 `json:"five_star_percent"`
	FourStarPercent  float32 `json:"four_star_percent"`
	ThreeStarPercent float32 `json:"three_star_percent"`
	TwoStarPercent   float32 `json:"two_star_percent"`
	OneStarPercent   float32 `json:"one_star_percent"`
}

type ReviewResponse struct {
	ID          int64                      `json:"id,string"`
	Star        uint8                      `json:"star"`
	Content     string                     `json:"content"`
	CreatedAt   time.Time                  `json:"created_at"`
	BookingRoom *SimpleBookingRoomResponse `json:"booking_room"`
}

type BasicReviewResponse struct {
	ID      int64  `json:"id,string"`
	Star    uint8  `json:"star"`
	Content string `json:"content"`
}
