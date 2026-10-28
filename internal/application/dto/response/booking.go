package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type SimpleBookingRoomResponse struct {
	ID      int64                 `json:"id,string"`
	Booking *BasicBookingResponse `json:"booking"`
	Room    *BasicRoomResponse    `json:"room"`
}

type BookingRoomResponse struct {
	ID             int64                          `json:"id,string"`
	OrderService   []*BasicOrderServiceResponse   `json:"order_services"`
	SupportTickets []*SimpleSupportTicketResponse `json:"support_tickets"`
	Guests         []*SimpleGuestResponse         `json:"guests"`
	Review         *SimpleReviewResponse          `json:"review"`
}

type BasicBookingRoomResponse struct {
	ID   int64              `json:"id,string"`
	Code string             `json:"code"`
	Room *BasicRoomResponse `json:"room"`
}

type BasicSourceResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type BasicBookingResponse struct {
	ID            int64     `json:"id,string"`
	BookingNumber string    `json:"booking_number"`
	GuestFullName string    `json:"guest_full_name"`
	GuestEmail    *string   `json:"guest_email"`
	GuestPhone    *string   `json:"guest_phone"`
	CheckIn       time.Time `json:"check_in"`
	CheckOut      time.Time `json:"check_out"`
}

type SourceResponse struct {
	ID        int64              `json:"id,string"`
	Name      string             `json:"name"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	CreatedBy *BasicUserResponse `json:"created_by"`
	UpdatedBy *BasicUserResponse `json:"updated_by"`
}

type BookingStatusSummaryResponse struct {
	All        int64 `json:"all"`
	CheckedIn  int64 `json:"checked-in"`
	Confirmed  int64 `json:"confirmed"`
	CheckedOut int64 `json:"checked-out"`
	Cancelled  int64 `json:"cancelled"`
}

type BookingTimelineResponse struct {
	StartDate time.Time                          `json:"start_date"`
	EndDate   time.Time                          `json:"end_date"`
	Status    *RoomStatusSummaryResponse         `json:"status"`
	Chart     []*BookingTimelineRoomTypeResponse `json:"chart"`
}

type BookingTimelineRoomTypeResponse struct {
	RoomTypeID   int64                          `json:"room_type_id,string"`
	RoomTypeName string                         `json:"room_type_name"`
	RoomsCount   int64                          `json:"rooms_count"`
	Rooms        []*BookingTimelineRoomResponse `json:"rooms"`
}

type BookingTimelineRoomResponse struct {
	RoomID   int64                           `json:"room_id,string"`
	RoomName string                          `json:"room_name"`
	Events   []*BookingTimelineEventResponse `json:"events"`
}

type BookingTimelineEventResponse struct {
	BookingID     int64               `json:"booking_id,string"`
	CheckIn       time.Time           `json:"check_in"`
	CheckOut      time.Time           `json:"check_out"`
	Title         string              `json:"title"`
	Status        model.BookingStatus `json:"status"`
	BookingNumber string              `json:"booking_number"`
}

type BookingPredictResponse struct {
	CancellationProbability float32 `json:"cancellation_probability"`
	WillCancel              bool    `json:"will_cancel"`
	RiskLevel               string  `json:"risk_level"`
}

type BookingResponse struct {
	ID                 int64                       `json:"id,string"`
	BookingNumber      string                      `json:"booking_number"`
	GuestFullName      string                      `json:"guest_full_name"`
	GuestEmail         *string                     `json:"guest_email"`
	GuestPhone         *string                     `json:"guest_phone"`
	CheckIn            time.Time                   `json:"check_in"`
	CheckOut           time.Time                   `json:"check_out"`
	RoomNumber         uint32                      `json:"room_number"`
	GuestNumber        string                      `json:"guest_number"`
	BookedOn           time.Time                   `json:"booked_on"`
	TotalNetPrice      float64                     `json:"total_net_price"`
	TotalSellPrice     float64                     `json:"total_sell_price"`
	PromotionName      *string                     `json:"promotion_name"`
	MealPlan           *string                     `json:"meal_plan"`
	BookingPreferences *string                     `json:"booking_preferences"`
	BookingConditions  *string                     `json:"booking_conditions"`
	Note               *string                     `json:"note"`
	Status             model.BookingStatus         `json:"status"`
	Branch             *BasicBranchResponse        `json:"branch"`
	Source             *BasicSourceResponse        `json:"source"`
	RoomType           *BasicRoomTypeResponse      `json:"room_type"`
	CreatedAt          time.Time                   `json:"created_at"`
	UpdatedAt          time.Time                   `json:"updated_at"`
	CreatedBy          *BasicUserResponse          `json:"created_by"`
	UpdatedBy          *BasicUserResponse          `json:"updated_by"`
	BookingRooms       []*BasicBookingRoomResponse `json:"booking_rooms"`
	Predict            *BookingPredictResponse     `json:"predict"`
}

type SimpleBookingResponse struct {
	ID             int64                       `json:"id,string"`
	BookingNumber  string                      `json:"booking_number"`
	GuestFullName  string                      `json:"guest_full_name"`
	GuestEmail     *string                     `json:"guest_email"`
	GuestPhone     *string                     `json:"guest_phone"`
	CheckIn        time.Time                   `json:"check_in"`
	CheckOut       time.Time                   `json:"check_out"`
	RoomNumber     uint32                      `json:"room_number"`
	TotalSellPrice float64                     `json:"total_sell_price"`
	RoomType       *BasicRoomTypeResponse      `json:"room_type"`
	Branch         *BasicBranchResponse        `json:"branch"`
	Source         *BasicSourceResponse        `json:"source"`
	BookingRooms   []*BasicBookingRoomResponse `json:"booking_rooms"`
	Status         model.BookingStatus         `json:"status"`
	Predict        *BookingPredictResponse     `json:"predict"`
}
