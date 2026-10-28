package response

import (
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type OrderServiceResponse struct {
	ID           int64                     `json:"id,string"`
	BookingRoom  *BasicBookingRoomResponse `json:"booking_room"`
	Service      *ServiceResponse          `json:"service"`
	Quantity     uint32                    `json:"quantity"`
	UnitPrice    float64                   `json:"unit_price"`
	TotalPrice   float64                   `json:"total_price"`
	Status       model.OrderServiceStatus  `json:"status"`
	GuestNote    *string                   `json:"guest_note"`
	StaffNote    *string                   `json:"staff_note"`
	CancelReason *string                   `json:"cancel_reason"`
	EstimatedAt  time.Time                 `json:"estimated_at"`
	CreatedAt    time.Time                 `json:"created_at"`
	UpdatedAt    time.Time                 `json:"updated_at"`
	CreatedBy    *BasicUserResponse        `json:"created_by"`
	UpdatedBy    *BasicUserResponse        `json:"updated_by"`
}

type BasicOrderServiceResponse struct {
	ID          int64                     `json:"id,string"`
	Service     *BasicServiceResponse     `json:"service"`
	BookingRoom *BasicBookingRoomResponse `json:"booking_room"`
	Quantity    uint32                    `json:"quantity"`
	UnitPrice   float64                   `json:"unit_price"`
	TotalPrice  float64                   `json:"total_price"`
	Status      model.OrderServiceStatus  `json:"status"`
	EstimatedAt time.Time                 `json:"estimated_at"`
	GuestNote   *string                   `json:"guest_note"`
}

type BasicServiceTypeResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
	Icon string `json:"icon"`
	Slug string `json:"slug"`
}

type BasicServiceImageResponse struct {
	Key string `json:"key"`
}

type BasicServiceResponse struct {
	ID   int64  `json:"id,string"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type SimpleOrderServiceResponse struct {
	ID           int64                    `json:"id,string"`
	Service      *ServiceResponse         `json:"service"`
	Quantity     uint32                   `json:"quantity"`
	UnitPrice    float64                  `json:"unit_price"`
	TotalPrice   float64                  `json:"total_price"`
	Status       model.OrderServiceStatus `json:"status"`
	GuestNote    *string                  `json:"guest_note"`
	StaffNote    *string                  `json:"staff_note"`
	CancelReason *string                  `json:"cancel_reason"`
	EstimatedAt  time.Time                `json:"estimated_at"`
	CreatedAt    time.Time                `json:"created_at"`
}

type ServiceImageResponse struct {
	Key         string `json:"key"`
	IsThumbnail bool   `json:"is_thumbnail"`
	SortOrder   uint32 `json:"sort_order"`
}

type ServiceTypeResponse struct {
	ID           int64                    `json:"id,string"`
	Name         string                   `json:"name"`
	Icon         string                   `json:"icon"`
	CreatedAt    time.Time                `json:"created_at"`
	UpdatedAt    time.Time                `json:"updated_at"`
	CreatedBy    *BasicUserResponse       `json:"created_by"`
	UpdatedBy    *BasicUserResponse       `json:"updated_by"`
	Department   *BasicDepartmentResponse `json:"department"`
	ServiceCount int64                    `json:"service_count"`
}

type SimpleServiceTypeResponse struct {
	ID       int64              `json:"id,string"`
	Name     string             `json:"name"`
	Icon     string             `json:"icon"`
	Services []*ServiceResponse `json:"services"`
}

type ServiceResponse struct {
	ID              int64                      `json:"id,string"`
	Name            string                     `json:"name"`
	Slug            string                     `json:"slug"`
	Price           float64                    `json:"price"`
	DiscountType    *model.ServiceDiscountType `json:"discount_type"`
	DiscountValue   *float64                   `json:"discount_value"`
	DiscountStartAt *time.Time                 `json:"discount_start_at"`
	DiscountEndAt   *time.Time                 `json:"discount_end_at"`
	Thumbnail       *BasicServiceImageResponse `json:"thumbnail"`
	Branches        []*BasicBranchResponse     `json:"branches"`
}

type ExtendedServiceResponse struct {
	ID              int64                      `json:"id,string"`
	Name            string                     `json:"name"`
	Price           float64                    `json:"price"`
	DiscountType    *model.ServiceDiscountType `json:"discount_type"`
	DiscountValue   *float64                   `json:"discount_value"`
	DiscountStartAt *time.Time                 `json:"discount_start_at"`
	DiscountEndAt   *time.Time                 `json:"discount_end_at"`
	Description     string                     `json:"description"`
	ServiceType     *BasicServiceTypeResponse  `json:"service_type"`
	Branches        []*BasicBranchResponse     `json:"branches"`
	Images          []*ServiceImageResponse    `json:"images"`
}

type ServiceDetailsResponse struct {
	ID              int64                      `json:"id,string"`
	Name            string                     `json:"name"`
	Price           float64                    `json:"price"`
	DiscountType    *model.ServiceDiscountType `json:"discount_type"`
	DiscountValue   *float64                   `json:"discount_value"`
	DiscountStartAt *time.Time                 `json:"discount_start_at"`
	DiscountEndAt   *time.Time                 `json:"discount_end_at"`
	IsActive        bool                       `json:"is_active"`
	Description     string                     `json:"description"`
	ServiceType     *BasicServiceTypeResponse  `json:"service_type"`
	Branches        []*BasicBranchResponse     `json:"branches"`
	Images          []*ServiceImageResponse    `json:"images"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	CreatedBy       *BasicUserResponse         `json:"created_by"`
	UpdatedBy       *BasicUserResponse         `json:"updated_by"`
}

type SimpleServiceResponse struct {
	ID          int64                     `json:"id,string"`
	Name        string                    `json:"name"`
	ServiceType *BasicServiceTypeResponse `json:"service_type"`
	Branches    []*BasicBranchResponse    `json:"branches"`
	Price       float64                   `json:"price"`
	IsActive    bool                      `json:"is_active"`
	CreatedAt   time.Time                 `json:"created_at"`
}
