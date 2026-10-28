package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type RoomResponse struct {
	ID        int64                  `json:"id,string"`
	Name      string                 `json:"name"`
	RoomType  *BasicRoomTypeResponse `json:"room_type"`
	Floor     *BasicFloorResponse    `json:"floor"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	CreatedBy *BasicUserResponse     `json:"created_by"`
	UpdatedBy *BasicUserResponse     `json:"updated_by"`
}

type RoomTimelineStatusSummaryResponse struct {
	All             int64 `json:"all"`
	Available       int64 `json:"available"`
	Occupied        int64 `json:"occupied"`
	UpcomingCheckIn int64 `json:"upcoming_checkin"`
	OverdueCheckout int64 `json:"overdue_checkout"`
	OverdueCheckIn  int64 `json:"overdue_checkin"`
}

type RoomTypeResponse struct {
	ID             int64                    `json:"id,string"`
	Name           string                   `json:"name"`
	Price          float64                  `json:"price"`
	SurfaceArea    float32                  `json:"surface_area"`
	MaxPeople      uint8                    `json:"max_people"`
	CreatedAt      time.Time                `json:"created_at"`
	UpdatedAt      time.Time                `json:"updated_at"`
	CreatedBy      *BasicUserResponse       `json:"created_by"`
	UpdatedBy      *BasicUserResponse       `json:"updated_by"`
	RoomsCount     int64                    `json:"rooms_count"`
	RoomTypePrices []*RoomTypePriceResponse `json:"room_type_prices"`
	Branches       []*BasicBranchResponse   `json:"branches"`
}

type RoomTypePriceCalculationResponse struct {
	TotalDays    int64   `json:"total_days"`
	TotalPrice   float64 `json:"total_price"`
	AveragePrice float64 `json:"average_price"`
}

type BasicRoomTypeResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type FloorResponse struct {
	ID        int64                `json:"id,string"`
	Name      string               `json:"name"`
	Branch    *BasicBranchResponse `json:"branch"`
	CreatedAt time.Time            `json:"created_at"`
	UpdatedAt time.Time            `json:"updated_at"`
	CreatedBy *BasicUserResponse   `json:"created_by"`
	UpdatedBy *BasicUserResponse   `json:"updated_by"`
}

type BasicRoomResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type RoomTypePriceResponse struct {
	ID           int64                    `json:"id,string"`
	RoomTypeID   int64                    `json:"room_type_id,string"`
	PricingTagID int64                    `json:"pricing_tag_id,string"`
	PricingType  string                   `json:"pricing_type"`
	PricingValue float64                  `json:"pricing_value"`
	PricingTag   *BasicPricingTagResponse `json:"pricing_tag"`
}

type PricingTagResponse struct {
	ID        int64              `json:"id,string"`
	Name      string             `json:"name"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedBy *BasicUserResponse `json:"created_by"`
	UpdatedBy *BasicUserResponse `json:"updated_by"`
}

type PricingCalendarResponse struct {
	Date         time.Time                `json:"date"`
	PricingTagID *int64                   `json:"pricing_tag_id,string"`
	PricingTag   *BasicPricingTagResponse `json:"pricing_tag"`
}

type RoomTimelineResponse struct {
	ServerNow time.Time                          `json:"server_now"`
	Status    *RoomTimelineStatusSummaryResponse `json:"status"`
	Timeline  []*RoomTimelineRoomTypeResponse    `json:"timeline"`
}

type RoomTimelineRoomTypeResponse struct {
	RoomTypeID   int64                       `json:"room_type_id,string"`
	RoomTypeName string                      `json:"room_type_name"`
	RoomsCount   int64                       `json:"rooms_count"`
	Rooms        []*RoomTimelineRoomResponse `json:"rooms"`
}

type BasicPricingTagResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type BasicFloorResponse struct {
	ID     int64                `json:"id,string"`
	Name   string               `json:"name"`
	Branch *BasicBranchResponse `json:"branch"`
}

type RoomTimelineRoomResponse struct {
	RoomID   int64                        `json:"room_id,string"`
	RoomName string                       `json:"room_name"`
	Price    float64                      `json:"price"`
	Booking  *RoomTimelineBookingResponse `json:"booking"`
}

type RoomTimelineBookingResponse struct {
	BookingID     int64               `json:"booking_id,string"`
	BookingNumber string              `json:"booking_number"`
	Status        model.BookingStatus `json:"status"`
	SourceName    string              `json:"source_name"`
	GuestName     string              `json:"guest_name"`
	CheckIn       time.Time           `json:"check_in"`
	CheckOut      time.Time           `json:"check_out"`
}

type RoomStatusSummaryResponse struct {
	All       int64 `json:"all"`
	Available int64 `json:"available"`
	Occupied  int64 `json:"occupied"`
	Booked    int64 `json:"booked"`
}

type SimpleRoomResponse struct {
	ID          int64  `json:"id,string"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}
