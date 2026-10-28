package model

import "time"

type (
	SupportTicketPriority string
	SupportTicketStatus   string
)

const (
	SupportTicketPriorityLow                 SupportTicketPriority = "low"
	SupportTicketPriorityMedium              SupportTicketPriority = "medium"
	SupportTicketPriorityHigh                SupportTicketPriority = "high"
	SupportTicketPriorityUrgent              SupportTicketPriority = "urgent"
	SupportTicketStatusPending               SupportTicketStatus   = "pending"
	SupportTicketStatusAccepted              SupportTicketStatus   = "accepted"
	SupportTicketStatusCancelled             SupportTicketStatus   = "cancelled"
	SupportTicketStatusDone                  SupportTicketStatus   = "done"
	SupportTicketPriorityLowDisplayNameVI                          = "Thấp"
	SupportTicketPriorityLowDisplayNameEN                          = "Low"
	SupportTicketPriorityMediumDisplayNameVI                       = "Trung bình"
	SupportTicketPriorityMediumDisplayNameEN                       = "Medium"
	SupportTicketPriorityHighDisplayNameVI                         = "Cao"
	SupportTicketPriorityHighDisplayNameEN                         = "High"
	SupportTicketPriorityUrgentDisplayNameVI                       = "Gấp"
	SupportTicketPriorityUrgentDisplayNameEN                       = "Urgent"
)

type SupportTicket struct {
	Base

	BookingRoomID     int64
	Content           string
	Status            SupportTicketStatus
	Priority          SupportTicketPriority
	ResolutionSeconds *int64
	SupportCategoryID int64
	UpdatedByID       *int64

	SupportCategory *SupportCategory
	BookingRoom     *BookingRoom
	UpdatedBy       *User
}

type SupportTicketQuery struct {
	Page      uint32
	Limit     uint32
	Sort      string
	Order     string
	Status    SupportTicketStatus
	Priority  SupportTicketPriority
	FromDate  time.Time
	ToDate    time.Time
	BranchIDs []int64
}

func NewSupportTicket(
	id int64,
	bookingRoomID int64,
	content string,
	status SupportTicketStatus,
	priority SupportTicketPriority,
	supportCategoryID int64,
) *SupportTicket {
	return &SupportTicket{
		Base: Base{
			ID: id,
		},
		BookingRoomID:     bookingRoomID,
		Content:           content,
		Status:            status,
		Priority:          priority,
		SupportCategoryID: supportCategoryID,
	}
}
