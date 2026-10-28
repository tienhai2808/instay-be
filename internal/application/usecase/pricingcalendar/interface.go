package pricingcalendar

import (
	"context"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
)

type PricingCalendarUseCase interface {
	List(ctx context.Context, schema string) ([]*response.PricingCalendarResponse, error)

	UpdateBulk(ctx context.Context, schema string, req request.UpdatePricingCalendarsRequest) error
}
