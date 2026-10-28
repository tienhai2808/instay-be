package model

import "time"

type CustomerGender string

const (
	CustomerGenderMale   CustomerGender = "male"
	CustomerGenderFemale CustomerGender = "female"
)

type Customer struct {
	ID            int64
	CreatedAt     time.Time
	BirthYear     uint32
	Gender        CustomerGender
	Phone         string
	BookingRoomID int64

	BookingRoom *BookingRoom
}

type CustomerQuery struct {
	Page      uint32
	Limit     uint32
	Sort      string
	Order     string
	Search    string
	BranchIDs []int64
}

func NewCustomer(
	id, bookingRoomID int64,
	birthYear uint32,
	gender CustomerGender,
	phone string,
) *Customer {
	return &Customer{
		ID:            id,
		BirthYear:     birthYear,
		Gender:        gender,
		Phone:         phone,
		BookingRoomID: bookingRoomID,
	}
}
