package model

import "time"

type (
	NotificationReceiverType string
	NotificationContentType  string
)

const (
	NotificationReceiverTypeGuest        NotificationReceiverType = "guest"
	NotificationReceiverTypeStaff        NotificationReceiverType = "staff"
	NotificationContentTypeOrderService  NotificationContentType  = "order_service"
	NotificationContentTypeSupportTicket NotificationContentType  = "support_ticket"
	NotificationContentTypeBooking       NotificationContentType  = "booking"
	NotificationContentTypeEvent         NotificationContentType  = "event"
)

type NotificationQuery struct {
	Page      uint32
	Limit     uint32
	BranchIDs []int64
}

type Notification struct {
	ID            int64
	DepartmentID  *int64
	BookingRoomID *int64
	ReceiverType  NotificationReceiverType
	ContentType   NotificationContentType
	ContentKey    string
	ContentVI     string
	ContentEN     string
	IsRead        bool
	CreatedAt     time.Time

	Department  *Department
	BookingRoom *BookingRoom
	Users       []*User
}

func NewNotification(
	id int64,
	departmentID, bookingRoomID *int64,
	receiverType NotificationReceiverType,
	contentType NotificationContentType,
	contentKey string,
	contentVI, contentEN string,
	isRead bool,
) *Notification {
	return &Notification{
		ID:            id,
		DepartmentID:  departmentID,
		BookingRoomID: bookingRoomID,
		ReceiverType:  receiverType,
		ContentType:   contentType,
		ContentKey:    contentKey,
		ContentVI:     contentVI,
		ContentEN:     contentEN,
		IsRead:        isRead,
	}
}
