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

type dailyRoomTypeRateRepositoryImpl struct {
	BaseRepository
}

const dailyRoomTypeRateUpsertBatchSize = 500

func NewDailyRoomTypeRateRepository(baseRepo BaseRepository) repository.DailyRoomTypeRateRepository {
	return &dailyRoomTypeRateRepositoryImpl{baseRepo}
}

func (r *dailyRoomTypeRateRepositoryImpl) UpsertBulk(ctx context.Context, schema string, dailyRoomTypeRates []*model.DailyRoomTypeRate) error {
	entities := mapper.ToDailyRoomTypeRateEntities(dailyRoomTypeRates)

	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DailyRoomTypeRateTableName)).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "room_type_id"},
				{Name: "pricing_calendar_date"},
			},
			DoUpdates: clause.AssignmentColumns([]string{"final_price"}),
		}).
		CreateInBatches(entities, dailyRoomTypeRateUpsertBatchSize).Error
}

func (r *dailyRoomTypeRateRepositoryImpl) SumFinalPriceByRoomTypeIDAndDateRange(ctx context.Context, schema string, roomTypeID int64, fromDate, toDate time.Time) (float64, int64, error) {
	from := utils.NormalizeDateOnly(fromDate)
	to := utils.NormalizeDateOnly(toDate)

	if !from.Before(to) {
		return 0, 0, nil
	}

	var result struct {
		SumFinalPrice float64 `gorm:"column:sum_final_price"`
		TotalDays     int64   `gorm:"column:total_days"`
	}

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DailyRoomTypeRateTableName)).
		Select("COALESCE(SUM(final_price), 0) AS sum_final_price, COUNT(*) AS total_days").
		Where("room_type_id = ? AND pricing_calendar_date >= ? AND pricing_calendar_date < ?", roomTypeID, from, to).
		Scan(&result).Error; err != nil {
		return 0, 0, err
	}

	return result.SumFinalPrice, result.TotalDays, nil
}
