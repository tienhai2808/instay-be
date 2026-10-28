package implement

import (
	"context"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"gorm.io/gorm/clause"
)

const pricingCalendarUpsertBatchSize = 500

type pricingCalendarRepositoryImpl struct {
	BaseRepository
}

func NewPricingCalendarRepository(baseRepo BaseRepository) repository.PricingCalendarRepository {
	return &pricingCalendarRepositoryImpl{baseRepo}
}

func (r *pricingCalendarRepositoryImpl) FindAllWithPricingTag(ctx context.Context, schema string) ([]*model.PricingCalendar, error) {
	var pricingCalendarEntities []*entity.PricingCalendarEntity

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingCalendarTableName)).
		Model(&entity.PricingCalendarEntity{})

	db = db.Preload("PricingTag", r.preloadWithSchema(schema, entity.PricingTagTableName, "id", "name"))

	if err := db.
		Order("date ASC").
		Find(&pricingCalendarEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToPricingCalendars(pricingCalendarEntities), nil
}

func (r *pricingCalendarRepositoryImpl) FindAllByDateRange(ctx context.Context, schema string, fromDate, toDate time.Time) ([]*model.PricingCalendar, error) {
	var pricingCalendarEntities []*entity.PricingCalendarEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingCalendarTableName)).
		Where("date >= ? AND date < ?", utils.NormalizeDateOnly(fromDate), utils.NormalizeDateOnly(toDate)).
		Order("date ASC").
		Find(&pricingCalendarEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToPricingCalendars(pricingCalendarEntities), nil
}

func (r *pricingCalendarRepositoryImpl) EnsureDateRange(ctx context.Context, schema string, fromDate, toDate time.Time) error {
	startDate := utils.NormalizeDateOnly(fromDate)
	endDate := utils.NormalizeDateOnly(toDate)

	if !startDate.Before(endDate) {
		return nil
	}

	pricingCalendars := make([]*model.PricingCalendar, 0, int(endDate.Sub(startDate).Hours()/24))
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		pricingCalendars = append(pricingCalendars, &model.PricingCalendar{Date: d})
	}
	pricingCalendarEntities := mapper.ToPricingCalendarEntities(pricingCalendars)

	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingCalendarTableName)).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(pricingCalendarEntities, pricingCalendarUpsertBatchSize).Error
}

func (r *pricingCalendarRepositoryImpl) FindAllByPricingTagID(ctx context.Context, schema string, pricingTagID int64) ([]*model.PricingCalendar, error) {
	var pricingCalendarEntities []*entity.PricingCalendarEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingCalendarTableName)).
		Where("pricing_tag_id = ?", pricingTagID).
		Order("date ASC").
		Find(&pricingCalendarEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToPricingCalendars(pricingCalendarEntities), nil
}

func (r *pricingCalendarRepositoryImpl) UpdateAllByPricingTagIDs(ctx context.Context, schema string, pricingTagIDs []int64, updateData any) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingCalendarTableName)).
		Where("pricing_tag_id IN ?", pricingTagIDs).
		Updates(updateData).Error
}

func (r *pricingCalendarRepositoryImpl) UpsertBulk(ctx context.Context, schema string, pricingCalendars []*model.PricingCalendar) error {
	entities := mapper.ToPricingCalendarEntities(pricingCalendars)

	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.PricingCalendarTableName)).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "date"}},
			DoUpdates: clause.AssignmentColumns([]string{"pricing_tag_id"}),
		}).
		CreateInBatches(entities, pricingCalendarUpsertBatchSize).Error
}
