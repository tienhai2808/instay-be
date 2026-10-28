package model

import "time"

type Guest struct {
	Base
	Audit

	FullName       string
	DOB            time.Time
	National       string
	PassportNumber string
	CheckIn        time.Time
	CheckOut       time.Time
	BookingRoomID  int64

	BookingRoom *BookingRoom
}

type GuestQuery struct {
	Page      uint32
	Limit     uint32
	Sort      string
	Order     string
	Search    string
	FromDate  time.Time
	ToDate    time.Time
	BranchIDs []int64
}

func NewGuest(
	id, bookingRoomID int64,
	fullName, national, passportNumber string,
	dob, checkIn, checkOut time.Time,
	createdByID, updatedByID *int64,
) *Guest {
	return &Guest{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		FullName:       fullName,
		National:       national,
		PassportNumber: passportNumber,
		DOB:            dob,
		CheckIn:        checkIn,
		CheckOut:       checkOut,
		BookingRoomID:  bookingRoomID,
	}
}
