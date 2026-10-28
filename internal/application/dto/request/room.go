package request

import "time"

type RoomPaginationQuery struct {
	Page       uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit      uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Sort       string `form:"sort" json:"sort"`
	Order      string `form:"order" binding:"omitempty,oneof=asc desc" json:"order"`
	BranchID   int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
	FloorID    int64  `form:"floor_id" binding:"omitempty" json:"floor_id,string"`
	RoomTypeID int64  `form:"room_type_id" binding:"omitempty" json:"room_type_id,string"`
	Search     string `form:"search" json:"search"`
}

type RoomTimelineQuery struct {
	Search     string `form:"search" json:"search"`
	BranchID   int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
	RoomTypeID int64  `form:"room_type_id" binding:"omitempty" json:"room_type_id,string"`
	FloorID    int64  `form:"floor_id" binding:"omitempty" json:"floor_id,string"`
}

type CreateRoomRequest struct {
	Name       string `json:"name" binding:"required,min=2"`
	FloorID    int64  `json:"floor_id,string" binding:"required"`
	RoomTypeID int64  `json:"room_type_id,string" binding:"required"`
}

type PricingTagPaginationQuery struct {
	Page   uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit  uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search string `form:"search" json:"search"`
}

type RoomTypePaginationQuery struct {
	Page     uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search   string `form:"search" json:"search"`
	BranchID int64  `form:"branch_id" binding:"omitempty,gt=0" json:"branch_id,string"`
}

type CreateRoomTypeRequest struct {
	Name        string     `json:"name" binding:"required,min=2"`
	Price       float64    `json:"price" binding:"required,gt=0"`
	SurfaceArea float32    `json:"surface_area" binding:"required,gt=0"`
	MaxPeople   uint8      `json:"max_people" binding:"required,gt=0"`
	BranchIDs   Int64Slice `json:"branch_ids" binding:"required,min=1,dive,gt=0"`
}

type UpdateRoomTypeRequest struct {
	Name        string     `json:"name" binding:"required,min=2"`
	Price       float64    `json:"price" binding:"required,gt=0"`
	SurfaceArea float32    `json:"surface_area" binding:"required,gt=0"`
	MaxPeople   uint8      `json:"max_people" binding:"required,gt=0"`
	BranchIDs   Int64Slice `json:"branch_ids" binding:"required,min=1,dive,gt=0"`
}

type CalculateRoomTypeRateRequest struct {
	RoomTypeID int64     `json:"room_type_id,string" binding:"required,gt=0"`
	Quantity   uint8     `json:"quantity" binding:"required,gte=1"`
	FromDate   time.Time `json:"from_date" binding:"required"`
	ToDate     time.Time `json:"to_date" binding:"required"`
}

type CreatePricingTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdatePricingTagRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

type UpdatePricingCalendarRequest struct {
	PricingTagID int64       `json:"pricing_tag_id,string" binding:"required,gt=0"`
	Dates        []time.Time `json:"dates" binding:"required,min=1,dive,required"`
}

type UpdatePricingCalendarsRequest struct {
	PricingCalendars []UpdatePricingCalendarRequest `json:"pricing_calendars" binding:"required,min=1,dive"`
}

type UpdateRoomRequest struct {
	Name       string `json:"name" binding:"required,min=2"`
	FloorID    int64  `json:"floor_id,string" binding:"required"`
	RoomTypeID int64  `json:"room_type_id,string" binding:"required"`
}

type FloorPaginationQuery struct {
	Page     uint32 `form:"page" binding:"omitempty,min=1" json:"page"`
	Limit    uint32 `form:"limit" binding:"omitempty,min=1,max=100" json:"limit"`
	Search   string `form:"search" json:"search"`
	BranchID int64  `form:"branch_id" binding:"omitempty" json:"branch_id,string"`
}

type CreateFloorRequest struct {
	Name     string `json:"name" binding:"required,min=1"`
	BranchID int64  `json:"branch_id,string" binding:"required,gt=0"`
}

type UpdateFloorRequest struct {
	Name     string `json:"name" binding:"required,min=1"`
	BranchID int64  `json:"branch_id,string" binding:"required,gt=0"`
}

type CreateRoomTypePriceRequest struct {
	RoomTypeID   int64   `json:"room_type_id,string" binding:"required,gt=0"`
	PricingTagID int64   `json:"pricing_tag_id,string" binding:"required,gt=0"`
	PricingType  string  `json:"pricing_type" binding:"required,oneof=percent amount"`
	PricingValue float64 `json:"pricing_value" binding:"required"`
}

type UpdateRoomTypePriceRequest struct {
	PricingType  string  `json:"pricing_type" binding:"required,oneof=percent amount"`
	PricingValue float64 `json:"pricing_value" binding:"required"`
}
