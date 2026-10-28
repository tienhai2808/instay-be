package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type PatternList []string

func (p *PatternList) UnmarshalJSON(data []byte) error {
	raw := strings.TrimSpace(string(data))
	if raw == "" || raw == "null" {
		*p = nil
		return nil
	}

	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		single = strings.TrimSpace(single)
		if single == "" {
			*p = nil
			return nil
		}
		*p = PatternList{single}
		return nil
	}

	var multiple []string
	if err := json.Unmarshal(data, &multiple); err == nil {
		res := make([]string, 0, len(multiple))
		for _, item := range multiple {
			item = strings.TrimSpace(item)
			if item == "" {
				continue
			}
			res = append(res, item)
		}
		*p = res
		return nil
	}

	return fmt.Errorf("pattern list must be string or []string")
}

type TenantSettings struct {
	IMAP struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"imap"`

	EmailBooking struct {
		Subject            PatternList `json:"subject"`
		BookingNumber      PatternList `json:"booking_number"`
		Guest              PatternList `json:"guest"`
		CheckIn            PatternList `json:"check_in"`
		CheckOut           PatternList `json:"check_out"`
		BookedOn           PatternList `json:"booked_on"`
		RoomsBooked        PatternList `json:"rooms_booked"`
		Source             PatternList `json:"source"`
		TotalNetPrice      PatternList `json:"total_net_price"`
		TotalSellPrice     PatternList `json:"total_sell_price"`
		GuestNumber        PatternList `json:"guest_number"`
		RoomNumber         PatternList `json:"room_number"`
		PromotionName      PatternList `json:"promotion_name"`
		MealPlan           PatternList `json:"meal_plan"`
		BookingConditions  PatternList `json:"booking_conditions"`
		BookingPreferences PatternList `json:"booking_preferences"`
	} `json:"email_booking"`

	EmailCancel struct {
		Subject PatternList `json:"subject"`
		Bodies  PatternList `json:"bodies"`
	} `json:"email_cancel"`
}

type Tenant struct {
	Base

	Name       string
	Logo       *string
	Phone      string
	Domain     string
	SchemaName string
	Settings   TenantSettings
	IsActive   bool
}
