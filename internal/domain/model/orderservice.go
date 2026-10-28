package model

import "time"

type OrderServiceStatus string

const (
	OrderServiceStatusConfirmed OrderServiceStatus = "confirmed"
	OrderServiceStatusPending   OrderServiceStatus = "pending"
	OrderServiceStatusCancelled OrderServiceStatus = "cancelled"
)

type OrderService struct {
	Base
	Audit

	ServiceID     int64
	BookingRoomID int64
	Quantity      uint32
	UnitPrice     float64
	TotalPrice    float64
	Status        OrderServiceStatus
	GuestNote     *string
	StaffNote     *string
	CancelReason  *string
	EstimatedAt   time.Time

	Service     *Service
	BookingRoom *BookingRoom
}

type OrderServiceQuery struct {
	Page      uint32
	Limit     uint32
	Sort      string
	Order     string
	Status    OrderServiceStatus
	BranchIDs []int64
	FromDate  time.Time
	ToDate    time.Time
}

func NewOrderService(
	id, serviceID, bookingRoomID int64,
	quantity uint32,
	unitPrice float64,
	totalPrice float64,
	status OrderServiceStatus,
	guestNote, staffNote, cancelReason *string,
	estimatedAt time.Time,
	createdByID, updatedByID *int64,
) *OrderService {
	return &OrderService{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		ServiceID:     serviceID,
		BookingRoomID: bookingRoomID,
		Quantity:      quantity,
		UnitPrice:     unitPrice,
		TotalPrice:    totalPrice,
		Status:        status,
		GuestNote:     guestNote,
		StaffNote:     staffNote,
		CancelReason:  cancelReason,
		EstimatedAt:   estimatedAt,
	}
}
