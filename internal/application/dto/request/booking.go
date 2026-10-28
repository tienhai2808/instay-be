package request

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type VerifyBookingRoomRequest struct {
	Code string `json:"code" binding:"required"`
}

type CreateBookingRoomRequest struct {
	BookingID int64 `json:"booking_id,string" binding:"required"`
	RoomID    int64 `json:"room_id,string" binding:"required"`
}

type UpdateBookingRoomRequest struct {
	RoomID int64 `json:"room_id,string" binding:"required"`
}

type SourcePaginationQuery struct {
	Page   uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit  uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search string `form:"search" json:"search"`
}

type CreateSourceRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type UpdateSourceRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

type CreateBookingRequest struct {
	BookingNumber      *string   `json:"booking_number" binding:"omitempty,min=1"`
	BranchID           *int64    `json:"branch_id,string" binding:"omitempty,gt=0"`
	GuestFullName      string    `json:"guest_full_name" binding:"required,min=1"`
	GuestEmail         string    `json:"guest_email" binding:"required,email"`
	GuestPhone         string    `json:"guest_phone" binding:"required,min=10"`
	CheckIn            time.Time `json:"check_in" binding:"required"`
	CheckOut           time.Time `json:"check_out" binding:"required"`
	RoomTypeID         int64     `json:"room_type_id,string" binding:"required"`
	RoomNumber         uint32    `json:"room_number" binding:"required,gt=0"`
	GuestNumber        string    `json:"guest_number" binding:"required"`
	TotalNetPrice      *float64  `json:"total_net_price" binding:"omitempty,gt=0"`
	TotalSellPrice     *float64  `json:"total_sell_price" binding:"omitempty,gt=0"`
	PromotionName      *string   `json:"promotion_name" binding:"omitempty,min=2"`
	MealPlan           *string   `json:"meal_plan" binding:"omitempty,min=2"`
	BookingPreferences *string   `json:"booking_preferences" binding:"omitempty,min=2"`
	BookingConditions  *string   `json:"booking_conditions" binding:"omitempty,min=2"`
	Note               *string   `json:"note" binding:"omitempty,min=2"`
	SourceID           int64     `json:"source_id,string" binding:"required"`
}

type UpdateBookingRequest struct {
	BranchID           *int64    `json:"branch_id,string" binding:"omitempty,gt=0"`
	GuestFullName      string    `json:"guest_full_name" binding:"required,min=1"`
	GuestEmail         string    `json:"guest_email" binding:"required,email"`
	GuestPhone         string    `json:"guest_phone" binding:"required,min=10"`
	CheckIn            time.Time `json:"check_in" binding:"required"`
	CheckOut           time.Time `json:"check_out" binding:"required"`
	RoomTypeID         int64     `json:"room_type_id,string" binding:"required"`
	RoomNumber         uint32    `json:"room_number" binding:"required,gt=0"`
	GuestNumber        string    `json:"guest_number" binding:"required"`
	PromotionName      *string   `json:"promotion_name" binding:"omitempty,min=2"`
	MealPlan           *string   `json:"meal_plan" binding:"omitempty,min=2"`
	BookingPreferences *string   `json:"booking_preferences" binding:"omitempty,min=2"`
	BookingConditions  *string   `json:"booking_conditions" binding:"omitempty,min=2"`
	Note               *string   `json:"note" binding:"omitempty,min=2"`
}

type UpdateBookingStatusRequest struct {
	Status model.BookingStatus `json:"status" binding:"oneof=cancelled checked-in checked-out"`
}

type BookingPaginationQuery struct {
	Page             uint32              `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit            uint32              `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort             string              `form:"sort" json:"sort"`
	Order            string              `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	FromDate         time.Time           `form:"from_date" time_format:"2006-01-02"`
	ToDate           time.Time           `form:"to_date" time_format:"2006-01-02"`
	Status           model.BookingStatus `form:"status" binding:"omitempty,oneof=confirmed checked-in checked-out cancelled"`
	Search           string              `form:"search" json:"search"`
	SourceID         int64               `form:"source_id" binding:"omitempty" json:"source_id,string"`
	BranchID         int64               `form:"branch_id" binding:"omitempty,gt=0" json:"branch_id,string"`
	UnassignedBranch *bool               `form:"unassigned_branch" binding:"omitempty" json:"unassigned_branch"`
}

type BookingTimelineQuery struct {
	FromDate   time.Time `form:"from_date" time_format:"2006-01-02" binding:"required"`
	ToDate     time.Time `form:"to_date" time_format:"2006-01-02"`
	RoomTypeID int64     `form:"room_type_id" binding:"omitempty" json:"room_type_id,string"`
	BranchID   int64     `form:"branch_id" binding:"omitempty,gt=0" json:"branch_id,string"`
}
