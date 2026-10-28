package service

import "github.com/InstaySystem/is_v2-be/internal/domain/model"

type PricingService struct{}

func (PricingService) CalculateFinalPrice(basePrice float64, pricingType model.RoomTypePricePricingType, pricingValue float64) float64 {
	switch pricingType {
	case model.RoomTypePricePricingTypeAmount:
		return basePrice + pricingValue
	case model.RoomTypePricePricingTypePercent:
		return basePrice + (basePrice * pricingValue / 100)
	default:
		return basePrice
	}
}
