package model

import "time"

type BookingStatus string

const (
	BookingStatusConfirmed  BookingStatus = "confirmed"
	BookingStatusCancelled  BookingStatus = "cancelled"
	BookingStatusCheckedIn  BookingStatus = "checked-in"
	BookingStatusCheckedOut BookingStatus = "checked-out"
)

type BookingPredict struct {
	CancellationProbability float32 `json:"cancellation_probability"`
	WillCancel              bool    `json:"will_cancel"`
	RiskLevel               string  `json:"risk_level"`
}

type Booking struct {
	Base
	Audit

	BookingNumber      string
	GuestFullName      string
	GuestEmail         *string
	GuestPhone         *string
	CheckIn            time.Time
	CheckOut           time.Time
	RoomNumber         uint32
	GuestNumber        string
	BookedOn           time.Time
	TotalNetPrice      float64
	TotalSellPrice     float64
	PromotionName      *string
	MealPlan           *string
	BookingPreferences *string
	BookingConditions  *string
	Note               *string
	Status             BookingStatus
	BranchID           *int64
	SourceID           int64
	RoomTypeID         int64
	Predict            BookingPredict

	Branch       *Branch
	Source       *Source
	RoomType     *RoomType
	BookingRooms []*BookingRoom
}

type BookingQuery struct {
	Page             uint32
	Limit            uint32
	Sort             string
	Order            string
	FromDate         time.Time
	ToDate           time.Time
	Status           BookingStatus
	Search           string
	SourceID         int64
	BranchID         int64
	BranchIDs        []int64
	UnassignedBranch *bool
}

type BookingStatusSummary struct {
	All        int64
	CheckedIn  int64
	Confirmed  int64
	CheckedOut int64
	Cancelled  int64
}

type DashboardOverviewQuery struct {
	FromDate  time.Time
	ToDate    time.Time
	BranchIDs []int64
}

type DashboardBookingOverview struct {
	TotalRevenue     float64
	OccupancyRate    float64
	CancellationRate float64
	AverageDaysStay  float64
}

type DashboardRevenueQuery struct {
	RoomTypeID int64
	BranchID   int64
	BranchIDs  []int64
}

type DashboardRevenueGranularity string

const (
	DashboardRevenueGranularityDay   DashboardRevenueGranularity = "day"
	DashboardRevenueGranularityWeek  DashboardRevenueGranularity = "week"
	DashboardRevenueGranularityMonth DashboardRevenueGranularity = "month"
	DashboardRevenueGranularityYear  DashboardRevenueGranularity = "year"
)

type DashboardBookingRevenuePoint struct {
	Date    time.Time
	Revenue float64
}

type DashboardBookingRevenue struct {
	Granularity DashboardRevenueGranularity
	FromDate    time.Time
	ToDate      time.Time
	Chart       []*DashboardBookingRevenuePoint
}

type DashboardRoomTypeRevenuePoint struct {
	RoomTypeID   int64
	RoomTypeName string
	Revenue      float64
}

type DashboardRoomTypeRevenue struct {
	Chart []*DashboardRoomTypeRevenuePoint
}

type DashboardSourceRevenuePoint struct {
	SourceID       int64
	SourceName     string
	Revenue        float64
	RevenuePercent float64
}

type DashboardSourceRevenue struct {
	Chart []*DashboardSourceRevenuePoint
}

type BookingTimelineQuery struct {
	FromDate   time.Time
	ToDate     time.Time
	RoomTypeID int64
	BranchID   int64
	BranchIDs  []int64
}

type BookingTimeline struct {
	StartDate time.Time
	EndDate   time.Time
	Status    *RoomStatusSummary
	Data      []*BookingTimelineRoomType
}

type BookingTimelineRoomType struct {
	RoomTypeID   int64
	RoomTypeName string
	RoomsCount   int64
	Rooms        []*BookingTimelineRoom
}

type BookingTimelineRoom struct {
	RoomID   int64
	RoomName string
	Events   []*BookingTimelineEvent
}

type BookingTimelineEvent struct {
	BookingID     int64
	CheckIn       time.Time
	CheckOut      time.Time
	Title         string
	Status        BookingStatus
	BookingNumber string
}

func NewBooking(
	id int64,
	bookingNumber, guestFullName, guestNumber string,
	guestEmail, guestPhone *string,
	checkIn, checkOut, bookedOn time.Time,
	roomNumber uint32,
	totalNetPrice, totalSellPrice float64,
	promotionName, mealPlan, bookingPreferences, bookingConditions, note *string,
	status BookingStatus,
	roomTypeID, sourceID int64,
	branchID *int64,
	createdByID, updatedByID *int64,
) *Booking {
	return &Booking{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		BookingNumber:      bookingNumber,
		GuestFullName:      guestFullName,
		GuestEmail:         guestEmail,
		GuestPhone:         guestPhone,
		CheckIn:            checkIn,
		CheckOut:           checkOut,
		RoomNumber:         roomNumber,
		GuestNumber:        guestNumber,
		BookedOn:           bookedOn,
		TotalNetPrice:      totalNetPrice,
		TotalSellPrice:     totalSellPrice,
		PromotionName:      promotionName,
		MealPlan:           mealPlan,
		BookingPreferences: bookingPreferences,
		BookingConditions:  bookingConditions,
		Note:               note,
		Status:             status,
		BranchID:           branchID,
		SourceID:           sourceID,
		RoomTypeID:         roomTypeID,
	}
}
