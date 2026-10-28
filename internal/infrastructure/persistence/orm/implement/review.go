package implement

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
	"gorm.io/gorm"
)

type reviewRepositoryImpl struct {
	BaseRepository
}

type reviewDashboardSentimentTimeBucketConfig struct {
	Granularity model.DashboardRevenueGranularity
	BucketExpr  string
	SeriesStart time.Time
	SeriesEnd   time.Time
	Step        string
}

func NewReviewRepository(baseRepo BaseRepository) repository.ReviewRepository {
	return &reviewRepositoryImpl{baseRepo}
}

func (r *reviewRepositoryImpl) Create(ctx context.Context, schema string, review *model.Review) error {
	reviewEntity := mapper.ToReviewEntity(review)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName)).
		Create(reviewEntity).Error
}

func (r *reviewRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName)).
		Model(&entity.ReviewEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepositoryImpl) FindAllPaginated(ctx context.Context, schema string, query model.ReviewQuery) ([]*model.Review, int64, error) {
	var reviewEntities []*entity.ReviewEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName) + " AS r").
		Model(&entity.ReviewEntity{})

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Review{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("r.id", "r.content", "r.star", "r.created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&reviewEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToReviews(reviewEntities), total, nil
}

func (r *reviewRepositoryImpl) FindAllWithBookingRoomDetailsPaginated(ctx context.Context, schema string, query model.ReviewQuery) ([]*model.Review, int64, error) {
	var reviewEntities []*entity.ReviewEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName) + " AS r").
		Model(&entity.ReviewEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Review{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = r.applySorting(db, query)
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "booking_id", "room_id"))
	db = db.Preload("BookingRoom.Booking", r.preloadWithSchema(schema, entity.BookingTableName, "id", "booking_number", "guest_full_name", "guest_email", "guest_phone", "check_in", "check_out"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("r.id", "r.booking_room_id", "r.content", "r.star", "r.aspects", "r.created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&reviewEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToReviews(reviewEntities), total, nil
}

func (r *reviewRepositoryImpl) CountStarSummary(ctx context.Context, schema string, branchIDs []int64) (*model.ReviewStarSummary, error) {
	var summary mapper.ReviewStarSummaryProjection

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName)+" AS r").
		Select(
			"COUNT(*) AS all_count, "+
				"COUNT(*) FILTER (WHERE r.star = ?) AS five_star_count, "+
				"COUNT(*) FILTER (WHERE r.star = ?) AS four_star_count, "+
				"COUNT(*) FILTER (WHERE r.star = ?) AS three_star_count, "+
				"COUNT(*) FILTER (WHERE r.star = ?) AS two_star_count, "+
				"COUNT(*) FILTER (WHERE r.star = ?) AS one_star_count",
			5, 4, 3, 2, 1,
		)

	db = r.applyReviewBranchFilter(db, schema, branchIDs)

	if err := db.Scan(&summary).Error; err != nil {
		return nil, err
	}

	return mapper.ToReviewStarSummaryModel(summary), nil
}

func (r *reviewRepositoryImpl) FindDashboardReviewSummary(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewSummary, error) {
	baseReviewsCTE, branchArgs := r.reviewDashboardBaseReviewsCTE(schema, branchIDs)

	sql := fmt.Sprintf(
		`WITH %s,
		sentiments AS (
			SELECT LOWER(COALESCE(aspect.value->>'sentiment', '')) AS sentiment
			FROM base_reviews AS r
			JOIN LATERAL jsonb_array_elements(COALESCE(r.aspects, '[]'::jsonb)) AS aspect(value) ON TRUE
		)
		SELECT
			(SELECT COUNT(*) FROM base_reviews) AS review_count,
			(SELECT COALESCE(AVG(star), 0) FROM base_reviews) AS average_star,
			COALESCE(COUNT(*) FILTER (WHERE sentiment = 'positive'), 0) AS positive_count,
			COALESCE(COUNT(*) FILTER (WHERE sentiment = 'negative'), 0) AS negative_count,
			COALESCE(COUNT(*) FILTER (WHERE sentiment = 'neutral'), 0) AS neutral_count
		FROM sentiments`,
		baseReviewsCTE,
	)

	var summary mapper.ReviewDashboardSummaryProjection
	if err := r.getDB(ctx).Raw(sql, branchArgs...).Scan(&summary).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardReviewSummaryModel(summary), nil
}

func (r *reviewRepositoryImpl) FindDashboardReviewCategoryCount(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewCategoryCount, error) {
	baseReviewsCTE, branchArgs := r.reviewDashboardBaseReviewsCTE(schema, branchIDs)
	categoryExpr := buildReviewDashboardCategoryExpr("category_tag")

	sql := fmt.Sprintf(
		`WITH %s,
		aspect_rows AS (
			SELECT UPPER(COALESCE(aspect.value->>'category', '')) AS category_tag
			FROM base_reviews AS r
			JOIN LATERAL jsonb_array_elements(COALESCE(r.aspects, '[]'::jsonb)) AS aspect(value) ON TRUE
		),
		mapped AS (
			SELECT %s AS category
			FROM aspect_rows
		)
		SELECT
			category,
			COUNT(*) AS review_count
		FROM mapped
		WHERE category IS NOT NULL
		GROUP BY category`,
		baseReviewsCTE,
		categoryExpr,
	)

	var rows []*mapper.ReviewDashboardCategoryCountProjection
	if err := r.getDB(ctx).Raw(sql, branchArgs...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardReviewCategoryCountModel(rows), nil
}

func (r *reviewRepositoryImpl) FindDashboardReviewCategorySentiment(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewCategorySentiment, error) {
	baseReviewsCTE, branchArgs := r.reviewDashboardBaseReviewsCTE(schema, branchIDs)
	categoryExpr := buildReviewDashboardCategoryExpr("category_tag")

	sql := fmt.Sprintf(
		`WITH %s,
		aspect_rows AS (
			SELECT
				UPPER(COALESCE(aspect.value->>'category', '')) AS category_tag,
				LOWER(COALESCE(aspect.value->>'sentiment', '')) AS sentiment
			FROM base_reviews AS r
			JOIN LATERAL jsonb_array_elements(COALESCE(r.aspects, '[]'::jsonb)) AS aspect(value) ON TRUE
		),
		mapped AS (
			SELECT
				%s AS category,
				sentiment
			FROM aspect_rows
			WHERE sentiment IN ('positive', 'negative', 'neutral')
		)
		SELECT
			category,
			COUNT(*) FILTER (WHERE sentiment = 'positive') AS positive_count,
			COUNT(*) FILTER (WHERE sentiment = 'negative') AS negative_count,
			COUNT(*) FILTER (WHERE sentiment = 'neutral') AS neutral_count
		FROM mapped
		WHERE category IS NOT NULL
		GROUP BY category`,
		baseReviewsCTE,
		categoryExpr,
	)

	var rows []*mapper.ReviewDashboardCategorySentimentProjection
	if err := r.getDB(ctx).Raw(sql, branchArgs...).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardReviewCategorySentimentModel(rows), nil
}

func (r *reviewRepositoryImpl) FindDashboardReviewSentimentCount(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewSentimentCount, error) {
	baseReviewsCTE, branchArgs := r.reviewDashboardBaseReviewsCTE(schema, branchIDs)

	sql := fmt.Sprintf(
		`WITH %s,
		sentiments AS (
			SELECT LOWER(COALESCE(aspect.value->>'sentiment', '')) AS sentiment
			FROM base_reviews AS r
			JOIN LATERAL jsonb_array_elements(COALESCE(r.aspects, '[]'::jsonb)) AS aspect(value) ON TRUE
		)
		SELECT
			COALESCE(COUNT(*) FILTER (WHERE sentiment = 'positive'), 0) AS positive_count,
			COALESCE(COUNT(*) FILTER (WHERE sentiment = 'negative'), 0) AS negative_count,
			COALESCE(COUNT(*) FILTER (WHERE sentiment = 'neutral'), 0) AS neutral_count
		FROM sentiments`,
		baseReviewsCTE,
	)

	var data mapper.ReviewDashboardSentimentCountProjection
	if err := r.getDB(ctx).Raw(sql, branchArgs...).Scan(&data).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardReviewSentimentCountModel(data), nil
}

func (r *reviewRepositoryImpl) FindDashboardReviewSentimentTime(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardReviewSentimentTime, error) {
	baseReviewsCTE, branchArgs := r.reviewDashboardBaseReviewsCTE(schema, branchIDs)

	rangeSQL := fmt.Sprintf(
		`WITH %s
		SELECT
			MIN(created_date) AS min_date,
			MAX(created_date) AS max_date
		FROM (
			SELECT r.created_at::date AS created_date
			FROM base_reviews AS r
			JOIN LATERAL jsonb_array_elements(COALESCE(r.aspects, '[]'::jsonb)) AS aspect(value) ON TRUE
			WHERE LOWER(COALESCE(aspect.value->>'sentiment', '')) IN ('positive', 'negative', 'neutral')
		) AS sentiment_rows`,
		baseReviewsCTE,
	)

	var rangeResult mapper.ReviewDashboardSentimentTimeRangeProjection
	if err := r.getDB(ctx).Raw(rangeSQL, branchArgs...).Scan(&rangeResult).Error; err != nil {
		return nil, err
	}

	if rangeResult.MinDate == nil || rangeResult.MaxDate == nil {
		return &model.DashboardReviewSentimentTime{
			Granularity: model.DashboardRevenueGranularityDay,
			Chart:       make([]*model.DashboardReviewSentimentTimePoint, 0),
		}, nil
	}

	minDate := utils.NormalizeDateOnly(*rangeResult.MinDate)
	maxDate := utils.NormalizeDateOnly(*rangeResult.MaxDate)
	bucketConfig := buildReviewDashboardSentimentTimeBucketConfig(minDate, maxDate)

	aggregateSQL := fmt.Sprintf(
		`WITH %s,
		series AS (
			SELECT generate_series(?::date, ?::date, ?::interval)::date AS bucket_start
		),
		aggregated AS (
			SELECT
				%s AS bucket_start,
				COUNT(*) FILTER (WHERE LOWER(COALESCE(aspect.value->>'sentiment', '')) = 'positive') AS positive_count,
				COUNT(*) FILTER (WHERE LOWER(COALESCE(aspect.value->>'sentiment', '')) = 'negative') AS negative_count,
				COUNT(*) FILTER (WHERE LOWER(COALESCE(aspect.value->>'sentiment', '')) = 'neutral') AS neutral_count
			FROM base_reviews AS r
			JOIN LATERAL jsonb_array_elements(COALESCE(r.aspects, '[]'::jsonb)) AS aspect(value) ON TRUE
			WHERE r.created_at >= ? AND r.created_at < ?
				AND LOWER(COALESCE(aspect.value->>'sentiment', '')) IN ('positive', 'negative', 'neutral')
			GROUP BY 1
		)
		SELECT
			s.bucket_start,
			COALESCE(a.positive_count, 0) AS positive_count,
			COALESCE(a.negative_count, 0) AS negative_count,
			COALESCE(a.neutral_count, 0) AS neutral_count
		FROM series AS s
		LEFT JOIN aggregated AS a ON a.bucket_start = s.bucket_start
		ORDER BY s.bucket_start ASC`,
		baseReviewsCTE,
		bucketConfig.BucketExpr,
	)

	args := make([]any, 0, len(branchArgs)+5)
	args = append(args, branchArgs...)
	args = append(args,
		bucketConfig.SeriesStart,
		bucketConfig.SeriesEnd,
		bucketConfig.Step,
		minDate,
		maxDate.AddDate(0, 0, 1),
	)

	var pointsResult []*mapper.ReviewDashboardSentimentTimePointProjection
	if err := r.getDB(ctx).Raw(
		aggregateSQL,
		args...,
	).Scan(&pointsResult).Error; err != nil {
		return nil, err
	}

	return &model.DashboardReviewSentimentTime{
		Granularity: bucketConfig.Granularity,
		FromDate:    minDate,
		ToDate:      maxDate,
		Chart:       mapper.ToDashboardReviewSentimentTimePoints(pointsResult),
	}, nil
}

func (r *reviewRepositoryImpl) FindByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) (*model.Review, error) {
	var reviewEntity entity.ReviewEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName)).
		Where("booking_room_id = ?", bookingRoomID).
		First(&reviewEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToReview(&reviewEntity), nil
}

func (r *reviewRepositoryImpl) UpdateByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ReviewTableName)).
		Model(&entity.ReviewEntity{}).
		Where("booking_room_id = ?", bookingRoomID).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrReviewNotFound
	}

	return nil
}

func (r *reviewRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.ReviewQuery) *gorm.DB {
	if query.FromDate != nil {
		db = db.Where("r.created_at >= ?", *query.FromDate)
	}

	if query.ToDate != nil {
		db = db.Where("r.created_at < ?", query.ToDate.AddDate(0, 0, 1))
	}

	if query.Star != 0 {
		db = db.Where("r.star = ?", query.Star)
	}

	db = r.applyReviewBranchFilter(db, schema, query.BranchIDs)

	return db
}

func (r *reviewRepositoryImpl) applySorting(db *gorm.DB, query model.ReviewQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at": "r.created_at",
		"star":       "r.star",
	}

	sortField := "r.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}

func (r *reviewRepositoryImpl) applyReviewBranchFilter(db *gorm.DB, schema string, branchIDs []int64) *gorm.DB {
	if len(branchIDs) == 0 {
		return db
	}

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	return db.
		Joins(fmt.Sprintf(`INNER JOIN %s br_scope ON br_scope.id = r.booking_room_id`, bookingRoomTable)).
		Joins(fmt.Sprintf(`INNER JOIN %s rm_scope ON rm_scope.id = br_scope.room_id`, roomTable)).
		Joins(fmt.Sprintf(`INNER JOIN %s f_scope ON f_scope.id = rm_scope.floor_id`, floorTable)).
		Where("f_scope.branch_id IN ?", branchIDs)
}

func (r *reviewRepositoryImpl) reviewDashboardBaseReviewsCTE(schema string, branchIDs []int64) (string, []any) {
	reviewTable := r.schemaTable(schema, entity.ReviewTableName)
	if len(branchIDs) == 0 {
		return fmt.Sprintf("base_reviews AS (SELECT r.* FROM %s AS r)", reviewTable), nil
	}

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	return fmt.Sprintf(
		`base_reviews AS (
			SELECT r.*
			FROM %s AS r
			JOIN %s AS br_dash ON br_dash.id = r.booking_room_id
			JOIN %s AS r_dash ON r_dash.id = br_dash.room_id
			JOIN %s AS f_dash ON f_dash.id = r_dash.floor_id
			WHERE f_dash.branch_id IN ?
		)`,
		reviewTable,
		bookingRoomTable,
		roomTable,
		floorTable,
	), []any{branchIDs}
}

func buildReviewDashboardCategoryExpr(categoryField string) string {
	return fmt.Sprintf(
		`CASE
			WHEN UPPER(split_part(%[1]s, '#', 2)) = 'PRICES' THEN 'price'
			WHEN UPPER(%[1]s) = 'SERVICE#GENERAL' THEN 'employee'
			WHEN UPPER(split_part(%[1]s, '#', 1)) = 'ROOMS' THEN 'room'
			WHEN UPPER(split_part(%[1]s, '#', 1)) = 'ROOM_AMENITIES' THEN 'amenities'
			WHEN UPPER(split_part(%[1]s, '#', 1)) = 'FACILITIES' THEN 'utility'
			WHEN UPPER(split_part(%[1]s, '#', 1)) IN ('HOTEL', 'FOOD&DRINKS') OR UPPER(%[1]s) = 'LOCATION#GENERAL' THEN 'service'
			ELSE NULL
		END`,
		categoryField,
	)
}

func buildReviewDashboardSentimentTimeBucketConfig(minDate, maxDate time.Time) reviewDashboardSentimentTimeBucketConfig {
	spanDays := int(maxDate.Sub(minDate)/(24*time.Hour)) + 1
	granularity := chooseBookingRevenueGranularity(spanDays)

	switch granularity {
	case model.DashboardRevenueGranularityWeek:
		return reviewDashboardSentimentTimeBucketConfig{
			Granularity: granularity,
			BucketExpr:  "date_trunc('week', r.created_at)::date",
			SeriesStart: truncateToISOWeek(minDate),
			SeriesEnd:   truncateToISOWeek(maxDate),
			Step:        "1 week",
		}
	case model.DashboardRevenueGranularityMonth:
		return reviewDashboardSentimentTimeBucketConfig{
			Granularity: granularity,
			BucketExpr:  "date_trunc('month', r.created_at)::date",
			SeriesStart: time.Date(minDate.Year(), minDate.Month(), 1, 0, 0, 0, 0, time.UTC),
			SeriesEnd:   time.Date(maxDate.Year(), maxDate.Month(), 1, 0, 0, 0, 0, time.UTC),
			Step:        "1 month",
		}
	case model.DashboardRevenueGranularityYear:
		return reviewDashboardSentimentTimeBucketConfig{
			Granularity: granularity,
			BucketExpr:  "date_trunc('year', r.created_at)::date",
			SeriesStart: time.Date(minDate.Year(), time.January, 1, 0, 0, 0, 0, time.UTC),
			SeriesEnd:   time.Date(maxDate.Year(), time.January, 1, 0, 0, 0, 0, time.UTC),
			Step:        "1 year",
		}
	default:
		return reviewDashboardSentimentTimeBucketConfig{
			Granularity: model.DashboardRevenueGranularityDay,
			BucketExpr:  "r.created_at::date",
			SeriesStart: minDate,
			SeriesEnd:   maxDate,
			Step:        "1 day",
		}
	}
}
