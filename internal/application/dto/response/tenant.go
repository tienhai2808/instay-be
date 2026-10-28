package response

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

type BasicTenantResponse struct {
	Logo  *string `json:"logo"`
	Name  string  `json:"name"`
	Phone string  `json:"phone"`
	City  string  `json:"city"`
}

type TenantSettingsResponse struct {
	IMAP struct {
		Host        string `json:"host"`
		Port        int    `json:"port"`
		User        string `json:"user"`
		HasPassword bool   `json:"has_password"`
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
