package request

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

type UpdateTenantRequest struct {
	Name  string  `json:"name" binding:"required,min=2,max=150"`
	Logo  *string `json:"logo" binding:"omitempty,max=150"`
	Phone string  `json:"phone" binding:"required,min=1,max=20"`
	City  string  `json:"city" binding:"required,min=1,max=100"`
}

type UpdateTenantSettingsRequest struct {
	IMAP struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"imap"`

	EmailBooking struct {
		Subject            model.PatternList `json:"subject"`
		BookingNumber      model.PatternList `json:"booking_number"`
		Guest              model.PatternList `json:"guest"`
		CheckIn            model.PatternList `json:"check_in"`
		CheckOut           model.PatternList `json:"check_out"`
		BookedOn           model.PatternList `json:"booked_on"`
		RoomsBooked        model.PatternList `json:"rooms_booked"`
		Source             model.PatternList `json:"source"`
		TotalNetPrice      model.PatternList `json:"total_net_price"`
		TotalSellPrice     model.PatternList `json:"total_sell_price"`
		GuestNumber        model.PatternList `json:"guest_number"`
		RoomNumber         model.PatternList `json:"room_number"`
		PromotionName      model.PatternList `json:"promotion_name"`
		MealPlan           model.PatternList `json:"meal_plan"`
		BookingConditions  model.PatternList `json:"booking_conditions"`
		BookingPreferences model.PatternList `json:"booking_preferences"`
	} `json:"email_booking"`

	EmailCancel struct {
		Subject model.PatternList `json:"subject"`
		Bodies  model.PatternList `json:"bodies"`
	} `json:"email_cancel"`
}
