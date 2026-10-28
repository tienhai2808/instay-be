package response

import "time"

type GuestResponse struct {
	ID             int64                     `json:"id,string"`
	FullName       string                    `json:"full_name"`
	DOB            time.Time                 `json:"dob"`
	National       string                    `json:"national"`
	PassportNumber string                    `json:"passport_number"`
	CheckIn        time.Time                 `json:"check_in"`
	CheckOut       time.Time                 `json:"check_out"`
	BookingRoom    *BasicBookingRoomResponse `json:"booking_room"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
	CreatedBy      *BasicUserResponse        `json:"created_by"`
	UpdatedBy      *BasicUserResponse        `json:"updated_by"`
}

type SimpleGuestResponse struct {
	ID             int64     `json:"id,string"`
	FullName       string    `json:"full_name"`
	DOB            time.Time `json:"dob"`
	National       string    `json:"national"`
	PassportNumber string    `json:"passport_number"`
}
