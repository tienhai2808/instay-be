package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type SupportTicketPriorityResponse struct {
	Name          model.SupportTicketPriority `json:"name"`
	DisplayNameVI string                      `json:"display_name_vi"`
	DisplayNameEN string                      `json:"display_name_en"`
}

type SupportCategoryResponse struct {
	ID                  int64                    `json:"id,string"`
	Name                string                   `json:"name"`
	CreatedAt           time.Time                `json:"created_at"`
	UpdatedAt           time.Time                `json:"updated_at"`
	CreatedBy           *BasicUserResponse       `json:"created_by"`
	UpdatedBy           *BasicUserResponse       `json:"updated_by"`
	Department          *BasicDepartmentResponse `json:"department"`
	SupportTicketsCount int64                    `json:"support_tickets_count"`
}

type BasicSupportCategoryResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
}

type SimpleSupportTicketResponse struct {
	ID              int64                         `json:"id,string"`
	SupportCategory *BasicSupportCategoryResponse `json:"support_category"`
	BookingRoom     *BasicBookingRoomResponse     `json:"booking_room,omitempty"`
	Content         string                        `json:"content"`
	Status          model.SupportTicketStatus     `json:"status"`
	Priority        model.SupportTicketPriority   `json:"priority"`
	CreatedAt       time.Time                     `json:"created_at"`
}

type SupportTicketResponse struct {
	ID                int64                         `json:"id,string"`
	SupportCategory   *BasicSupportCategoryResponse `json:"support_category"`
	BookingRoom       *BasicBookingRoomResponse     `json:"booking_room"`
	Content           string                        `json:"content"`
	Status            model.SupportTicketStatus     `json:"status"`
	Priority          model.SupportTicketPriority   `json:"priority"`
	ResolutionSeconds *int64                        `json:"resolution_seconds"`
	CreatedAt         time.Time                     `json:"created_at"`
	UpdatedAt         time.Time                     `json:"updated_at"`
	UpdatedBy         *BasicUserResponse            `json:"updated_by"`
}
