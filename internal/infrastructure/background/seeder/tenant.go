package seeder

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Seeder) tenantSeeder() error {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	var err error

	domain := os.Getenv("TA_DOMAIN")
	if domain == "" {
		return fmt.Errorf("domain is required")
	}

	schemaName := os.Getenv("TA_SCHEMA")
	if schemaName == "" {
		return fmt.Errorf("schema name is required")
	}

	name := os.Getenv("TA_NAME")
	if name == "" {
		return fmt.Errorf("name is required")
	}

	var logo *string
	if v, ok := os.LookupEnv("TA_LOGO"); ok && strings.TrimSpace(v) != "" {
		logo = &v
	}

	phone := os.Getenv("TA_PHONE")
	if phone == "" {
		return fmt.Errorf("phone is required")
	}

	city := os.Getenv("TA_CITY")
	if city == "" {
		return fmt.Errorf("city is required")
	}

	if schemaName, err = s.normalizeAndValidateSchema(schemaName); err != nil {
		return err
	}

	id, err := s.idPro.NextID()
	if err != nil {
		s.logPro.Error("Generate tenant id failed", port.String("schema", schemaName), port.String("domain", domain), port.Error(err))
		return err
	}

	tenant := &entity.TenantEntity{
		BaseEntity: entity.BaseEntity{
			ID: id,
		},
		Name:       name,
		Logo:       logo,
		City:       city,
		Phone:      phone,
		Domain:     domain,
		SchemaName: schemaName,
		IsActive:   true,
	}

	return s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tenant).Error; err != nil {
			s.logPro.Error("Create tenant error", port.String("schema", schemaName), port.String("domain", domain), port.Int64("tenant_id", tenant.ID), port.String("raw", err.Error()))
			if errors.Is(err, exceptions.ErrDomainAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrSchemaNameAlreadyExists) {
				return err
			}
			s.logPro.Error("Create tenant failed", port.String("schema", schemaName), port.String("domain", domain), port.Int64("tenant_id", tenant.ID), port.Error(err))
			return err
		}

		if err := tx.Exec(fmt.Sprintf(`CREATE SCHEMA "%s"`, schemaName)).Error; err != nil {
			s.logPro.Error("Create schema failed", port.String("schema", schemaName), port.String("domain", domain), port.Int64("tenant_id", tenant.ID), port.Error(err))
			return err
		}

		if err := tx.Exec(fmt.Sprintf("SET search_path TO %s", schemaName)).Error; err != nil {
			s.logPro.Error("Set search_path failed", port.String("schema", schemaName), port.Int64("tenant_id", tenant.ID), port.Error(err))
			return err
		}

		if err := tx.AutoMigrate(autoMigrateModels...); err != nil {
			s.logPro.Error("Auto migrate failed", port.String("schema", schemaName), port.Int64("tenant_id", tenant.ID), port.Error(err))
			return err
		}

		pricingCalendarRes, err := ensureCurrentAndNextYearPricingCalendar(ctx, tx, time.Now())
		if err != nil {
			s.logPro.Error("Create pricing calendar failed", port.String("schema", schemaName), port.Error(err))
			return err
		}
		if pricingCalendarRes.insertedDays > 0 {
			s.logPro.Info(
				"Created pricing calendar",
				port.String("schema", schemaName),
				port.Int("inserted_days", pricingCalendarRes.insertedDays),
				port.Any("inserted_years", pricingCalendarRes.insertedYears),
			)
		}

		if err := s.seedPermissions(ctx, tx); err != nil {
			s.logPro.Error("Create default permissions failed", port.String("schema", schemaName), port.Error(err))
			return err
		}

		if err := s.seedSources(ctx, tx); err != nil {
			s.logPro.Error("Create default sources failed", port.String("schema", schemaName), port.Error(err))
			return err
		}

		if err := s.seedDepartments(ctx, tx); err != nil {
			s.logPro.Error("Create default departments failed", port.String("schema", schemaName), port.Error(err))
			return err
		}

		s.logPro.Info("Create tenant completed successfully",
			port.String("schema", schemaName),
			port.String("domain", domain),
			port.Int64("tenant_id", tenant.ID),
			port.Int("models_count", len(autoMigrateModels)),
		)

		return nil
	})
}

func (s *Seeder) normalizeAndValidateSchema(schema string) (string, error) {
	schemaRegex := regexp.MustCompile(`^[a-z][a-z0-9]*$`)

	schema = strings.ToLower(schema)

	if !schemaRegex.MatchString(schema) || len(schema) > 63 {
		return "", fmt.Errorf("invalid schema name: %q", schema)
	}

	return schema, nil
}

const pricingCalendarInsertBatchSize = 500

type pricingCalendarEnsureResult struct {
	insertedYears []int
	insertedDays  int
}

func ensureCurrentAndNextYearPricingCalendar(ctx context.Context, tx *gorm.DB, now time.Time) (*pricingCalendarEnsureResult, error) {
	return ensurePricingCalendarYears(ctx, tx, []int{now.Year(), now.Year() + 1})
}

func ensurePricingCalendarYears(ctx context.Context, tx *gorm.DB, years []int) (*pricingCalendarEnsureResult, error) {
	normalizedYears := normalizeYears(years)
	if len(normalizedYears) == 0 {
		return &pricingCalendarEnsureResult{}, nil
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
		return &pricingCalendarEnsureResult{}, nil
	}

	rows := make([]*entity.PricingCalendarEntity, 0, totalDaysForYears(missingYears))
	for _, year := range missingYears {
		rows = append(rows, buildPricingCalendarRowsForYear(year)...)
	}

	if err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		CreateInBatches(rows, pricingCalendarInsertBatchSize).Error; err != nil {
		return nil, err
	}

	return &pricingCalendarEnsureResult{
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
