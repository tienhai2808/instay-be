package job

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type rotatePricingCalendarJob struct {
	logPro     port.LogProvider
	db         *gorm.DB
	schemaRepo repository.SchemeRepository
}

type rotatePricingCalendarEnsureResult struct {
	insertedYears []int
	insertedDays  int
}

const rotatePricingCalendarInsertBatchSize = 500

func NewRotatePricingCalendarJob(
	logPro port.LogProvider,
	db *gorm.DB,
	schemaRepo repository.SchemeRepository,
) Job {
	return &rotatePricingCalendarJob{
		logPro,
		db,
		schemaRepo,
	}
}

func (j *rotatePricingCalendarJob) Name() string {
	return "rotate_pricing_calendar"
}

func (j *rotatePricingCalendarJob) Run() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	startTime := time.Now()
	currentYear := startTime.Year()
	oldYear := currentYear - 1
	nextYear := currentYear + 1

	schemas, err := j.schemaRepo.FindAll(ctx)
	if err != nil {
		j.logPro.Error("Find all schema failed", port.Error(err))
		return
	}

	totalDeleted := int64(0)
	totalInserted := int64(0)
	successCount := 0
	failCount := 0

	for _, schema := range schemas {
		var deletedDays int64
		var insertedDays int
		var insertedYears []int

		err := j.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if err := setSchemaSearchPath(ctx, tx, schema); err != nil {
				return err
			}

			deletedDaysRes, err := deletePricingCalendarYear(ctx, tx, oldYear)
			if err != nil {
				return err
			}
			deletedDays = deletedDaysRes

			ensureRes, err := ensurePricingCalendarYears(ctx, tx, []int{nextYear})
			if err != nil {
				return err
			}
			insertedDays = ensureRes.insertedDays
			insertedYears = ensureRes.insertedYears

			return nil
		})
		if err != nil {
			j.logPro.Error(
				"rotate pricing calendar failed",
				port.String("schema", schema),
				port.Error(err),
			)
			failCount++
			continue
		}

		successCount++
		totalDeleted += deletedDays
		totalInserted += int64(insertedDays)

		if deletedDays > 0 || insertedDays > 0 {
			j.logPro.Info(
				"Rotated pricing calendar in schema",
				port.String("schema", schema),
				port.Int64("deleted_days", deletedDays),
				port.Int("inserted_days", insertedDays),
				port.Any("inserted_years", insertedYears),
			)
		}
	}

	j.logPro.Info(
		"Rotate pricing calendar completed",
		port.Int("total_schemas", len(schemas)),
		port.Int("success_count", successCount),
		port.Int("fail_count", failCount),
		port.Int("deleted_year", oldYear),
		port.Int("inserted_year", nextYear),
		port.Int64("total_deleted_days", totalDeleted),
		port.Int64("total_inserted_days", totalInserted),
		port.Duration("duration", time.Since(startTime)),
	)
}

func setSchemaSearchPath(ctx context.Context, tx *gorm.DB, schema string) error {
	return tx.WithContext(ctx).Exec(fmt.Sprintf(`SET search_path TO "%s"`, schema)).Error
}

func deletePricingCalendarYear(ctx context.Context, tx *gorm.DB, year int) (int64, error) {
	result := tx.WithContext(ctx).
		Where("date >= ? AND date < ?", startOfYear(year), startOfYear(year+1)).
		Delete(&entity.PricingCalendarEntity{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func ensurePricingCalendarYears(ctx context.Context, tx *gorm.DB, years []int) (*rotatePricingCalendarEnsureResult, error) {
	normalizedYears := normalizeYears(years)
	if len(normalizedYears) == 0 {
		return &rotatePricingCalendarEnsureResult{}, nil
	}

	missingYears := make([]int, 0, len(normalizedYears))
	for _, year := range normalizedYears {
		exists, err := pricingCalendarYearExists(ctx, tx, year)
		if err != nil {
			return nil, err
		}
		if !exists {
			missingYears = append(missingYears, year)
		}
	}

	if len(missingYears) == 0 {
		return &rotatePricingCalendarEnsureResult{}, nil
	}

	rows := make([]*entity.PricingCalendarEntity, 0, totalDaysForYears(missingYears))
	for _, year := range missingYears {
		rows = append(rows, buildPricingCalendarRowsForYear(year)...)
	}

	if err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(rows, rotatePricingCalendarInsertBatchSize).Error; err != nil {
		return nil, err
	}

	return &rotatePricingCalendarEnsureResult{
		insertedYears: missingYears,
		insertedDays:  len(rows),
	}, nil
}

func pricingCalendarYearExists(ctx context.Context, tx *gorm.DB, year int) (bool, error) {
	var existing entity.PricingCalendarEntity
	err := tx.WithContext(ctx).
		Select("date").
		Where("date >= ? AND date < ?", startOfYear(year), startOfYear(year+1)).
		Take(&existing).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return false, err
}

func buildPricingCalendarRowsForYear(year int) []*entity.PricingCalendarEntity {
	start := startOfYear(year)
	end := startOfYear(year + 1)

	rows := make([]*entity.PricingCalendarEntity, 0, daysInYear(year))
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		rows = append(rows, &entity.PricingCalendarEntity{Date: d})
	}

	return rows
}

func normalizeYears(years []int) []int {
	if len(years) == 0 {
		return nil
	}

	yearSet := make(map[int]struct{}, len(years))
	for _, year := range years {
		if year <= 0 {
			continue
		}
		yearSet[year] = struct{}{}
	}

	normalized := make([]int, 0, len(yearSet))
	for year := range yearSet {
		normalized = append(normalized, year)
	}
	sort.Ints(normalized)

	return normalized
}

func totalDaysForYears(years []int) int {
	total := 0
	for _, year := range years {
		total += daysInYear(year)
	}
	return total
}

func daysInYear(year int) int {
	return int(startOfYear(year+1).Sub(startOfYear(year)).Hours() / 24)
}

func startOfYear(year int) time.Time {
	return time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
}
