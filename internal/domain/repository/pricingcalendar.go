package repository

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

type PricingCalendarRepository interface {
	FindAllWithPricingTag(ctx context.Context, schema string) ([]*model.PricingCalendar, error)

	FindAllByDateRange(ctx context.Context, schema string, fromDate, toDate time.Time) ([]*model.PricingCalendar, error)

	EnsureDateRange(ctx context.Context, schema string, fromDate, toDate time.Time) error

	FindAllByPricingTagID(ctx context.Context, schema string, pricingTagID int64) ([]*model.PricingCalendar, error)

	UpdateAllByPricingTagIDs(ctx context.Context, schema string, pricingTagIDs []int64, updateData any) error

	UpsertBulk(ctx context.Context, schema string, pricingCalendars []*model.PricingCalendar) error
}
