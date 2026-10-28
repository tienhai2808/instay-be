package model

import "time"

type ServiceDiscountType string

const (
	ServiceDiscountTypePercent ServiceDiscountType = "percent"
	ServiceDiscountTypeAmount  ServiceDiscountType = "amount"
)

type Service struct {
	Base
	Audit

	Name            string
	Slug            string
	Price           float64
	IsActive        bool
	Description     string
	DiscountType    *ServiceDiscountType
	DiscountValue   *float64
	DiscountStartAt *time.Time
	DiscountEndAt   *time.Time
	ServiceTypeID   int64

	ServiceType   *ServiceType
	ServiceImages []*ServiceImage
	Branches      []*Branch
}

type ServiceQuery struct {
	Page          uint32
	Limit         uint32
	Sort          string
	Order         string
	ServiceTypeID int64
	BranchID      int64
	BranchIDs     []int64
	IsActive      *bool
	HasDiscount   *bool
	Search        string
}

func NewService(
	id, svcTypeID int64,
	name, slug, description string,
	price float64,
	isActive bool,
	createdByID, updatedByID *int64,
) *Service {
	return &Service{
		Base: Base{
			ID: id,
		},
		Audit: Audit{
			CreatedByID: createdByID,
			UpdatedByID: updatedByID,
		},
		Name:          name,
		Slug:          slug,
		Description:   description,
		Price:         price,
		IsActive:      isActive,
		ServiceTypeID: svcTypeID,
	}
}

func (s *Service) IsDiscountActiveAt(at time.Time) bool {
	if s == nil || s.DiscountType == nil || s.DiscountValue == nil {
		return false
	}

	if s.DiscountStartAt != nil && at.Before(*s.DiscountStartAt) {
		return false
	}
	if s.DiscountEndAt != nil && at.After(*s.DiscountEndAt) {
		return false
	}

	return true
}

func (s *Service) UnitPriceAt(at time.Time) float64 {
	if s == nil {
		return 0
	}

	base := s.Price
	if !s.IsDiscountActiveAt(at) {
		return base
	}

	switch *s.DiscountType {
	case ServiceDiscountTypePercent:
		discounted := base * (1 - (*s.DiscountValue / 100))
		if discounted < 0 {
			return 0
		}
		return discounted
	case ServiceDiscountTypeAmount:
		discounted := base - *s.DiscountValue
		if discounted < 0 {
			return 0
		}
		return discounted
	default:
		return base
	}
}
