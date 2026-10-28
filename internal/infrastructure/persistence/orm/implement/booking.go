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

type bookingRepositoryImpl struct {
	BaseRepository
}

type bookingDashboardRevenueRangeResult struct {
	MinCheckIn *time.Time `gorm:"column:min_check_in"`
	MaxCheckIn *time.Time `gorm:"column:max_check_in"`
}

type bookingDashboardTodayOperationResult struct {
	TodayCheckIn  int64 `gorm:"column:today_check_in"`
	TodayCheckOut int64 `gorm:"column:today_check_out"`
}

type bookingRevenueBucketConfig struct {
	Granularity model.DashboardRevenueGranularity
	BucketExpr  string
	SeriesStart time.Time
	SeriesEnd   time.Time
	Step        string
}

func NewBookingRepository(baseRepo BaseRepository) repository.BookingRepository {
	return &bookingRepositoryImpl{baseRepo}
}

func (r *bookingRepositoryImpl) Create(ctx context.Context, schema string, booking *model.Booking) error {
	bookingEntity := mapper.ToBookingEntity(booking)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingTableName)).
		Create(bookingEntity).Error
}

func (r *bookingRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingTableName)).
		Model(&entity.BookingEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrBookingNotFound
	}

	return nil
}

func (r *bookingRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingTableName)).
		Where("id = ?", id).
		Delete(&entity.BookingEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrBookingNotFound
	}

	return nil
}

func (r *bookingRepositoryImpl) CountStatusSummary(ctx context.Context, schema string, branchIDs ...[]int64) (*model.BookingStatusSummary, error) {
	var summary mapper.BookingStatusSummaryProjection

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingTableName)).
		Select(
			"COUNT(*) as all_count, "+
				"COUNT(*) FILTER (WHERE status = ?) as checked_in_count, "+
				"COUNT(*) FILTER (WHERE status = ?) as confirmed_count, "+
				"COUNT(*) FILTER (WHERE status = ?) as checked_out_count, "+
				"COUNT(*) FILTER (WHERE status = ?) as cancelled_count",
			model.BookingStatusCheckedIn,
			model.BookingStatusConfirmed,
			model.BookingStatusCheckedOut,
			model.BookingStatusCancelled,
		)

	if len(branchIDs) > 0 && len(branchIDs[0]) > 0 {
		db = db.Where("branch_id IN ?", branchIDs[0])
	}

	if err := db.Scan(&summary).Error; err != nil {
		return nil, err
	}

	return mapper.ToBookingStatusSummaryModel(summary), nil
}

func (r *bookingRepositoryImpl) FindDashboardOverview(ctx context.Context, schema string, query model.DashboardOverviewQuery) (*model.DashboardBookingOverview, error) {
	var summary mapper.BookingDashboardOverviewProjection

	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	endDateExclusive := query.ToDate.AddDate(0, 0, 1)

	totalRooms, err := r.countDashboardRooms(ctx, schema, query.BranchIDs)
	if err != nil {
		return nil, err
	}

	branchClause, branchArgs := dashboardBookingBranchClause(query.BranchIDs)
	args := []any{
		model.BookingStatusConfirmed,
		model.BookingStatusCheckedIn,
		model.BookingStatusCheckedOut,
		model.BookingStatusConfirmed,
		model.BookingStatusCheckedIn,
		model.BookingStatusCheckedOut,
		endDateExclusive,
		query.FromDate,
		model.BookingStatusCancelled,
		model.BookingStatusConfirmed,
		model.BookingStatusCheckedIn,
		model.BookingStatusCheckedOut,
		totalRooms,
		endDateExclusive,
		query.FromDate,
	}
	args = append(args, branchArgs...)

	if err := r.getDB(ctx).Raw(
		fmt.Sprintf(
			`SELECT
				COALESCE(SUM(CASE
					WHEN b.status IN (?, ?, ?) THEN b.total_sell_price
					ELSE 0
				END), 0) AS total_revenue,
				COALESCE(SUM(CASE
					WHEN b.status IN (?, ?, ?) THEN
						b.room_number::numeric * GREATEST(
							0,
							LEAST(b.check_out::date, ?::date) - GREATEST(b.check_in::date, ?::date)
						)
					ELSE 0
				END), 0) AS sold_room_nights,
				COUNT(*) AS total_booking_count,
				COUNT(*) FILTER (WHERE b.status = ?) AS cancelled_booking_count,
				COALESCE(AVG(CASE
					WHEN b.status IN (?, ?, ?) THEN GREATEST(1, b.check_out::date - b.check_in::date)::numeric
				END), 0) AS average_days_stay,
					?::bigint AS total_rooms
			FROM %s AS b
			WHERE b.check_in < ? AND b.check_out > ?%s`,
			bookingTable,
			branchClause,
		),
		args...,
	).Scan(&summary).Error; err != nil {
		return nil, err
	}

	dayCount := int64(query.ToDate.Sub(query.FromDate)/(24*time.Hour)) + 1
	occupancyRate := 0.0
	if summary.TotalRooms > 0 && dayCount > 0 {
		occupancyRate = (summary.SoldRoomNights / float64(summary.TotalRooms*dayCount)) * 100
	}

	cancellationRate := 0.0
	if summary.TotalBookingCount > 0 {
		cancellationRate = (float64(summary.CancelledBookingCount) / float64(summary.TotalBookingCount)) * 100
	}

	return mapper.ToDashboardBookingOverviewModel(summary, occupancyRate, cancellationRate), nil
}

func (r *bookingRepositoryImpl) FindDashboardRevenue(ctx context.Context, schema string, query model.DashboardRevenueQuery) (*model.DashboardBookingRevenue, error) {
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	whereClause, whereArgs := buildDashboardRevenueWhereClause(query)

	rangeSQL := fmt.Sprintf(
		`SELECT
			MIN(b.check_in::date) AS min_check_in,
			MAX(b.check_in::date) AS max_check_in
		FROM %s AS b
		WHERE %s`,
		bookingTable,
		whereClause,
	)

	var bookingRange bookingDashboardRevenueRangeResult
	if err := r.getDB(ctx).Raw(rangeSQL, whereArgs...).Scan(&bookingRange).Error; err != nil {
		return nil, err
	}

	if bookingRange.MinCheckIn == nil || bookingRange.MaxCheckIn == nil {
		return &model.DashboardBookingRevenue{
			Granularity: model.DashboardRevenueGranularityDay,
			Chart:       make([]*model.DashboardBookingRevenuePoint, 0),
		}, nil
	}

	minCheckIn := utils.NormalizeDateOnly(*bookingRange.MinCheckIn)
	maxCheckIn := utils.NormalizeDateOnly(*bookingRange.MaxCheckIn)
	bucketConfig := buildBookingRevenueBucketConfig(minCheckIn, maxCheckIn)

	aggregateWhereClause := whereClause + " AND b.check_in >= ? AND b.check_in < ?"
	aggregateArgs := make([]any, 0, 5+len(whereArgs))
	aggregateArgs = append(aggregateArgs, bucketConfig.SeriesStart, bucketConfig.SeriesEnd, bucketConfig.Step)
	aggregateArgs = append(aggregateArgs, whereArgs...)
	aggregateArgs = append(aggregateArgs, minCheckIn, maxCheckIn.AddDate(0, 0, 1))

	aggregateSQL := fmt.Sprintf(
		`WITH series AS (
			SELECT generate_series(?::date, ?::date, ?::interval)::date AS bucket_start
		),
		aggregated AS (
			SELECT
				%s AS bucket_start,
				COALESCE(SUM(b.total_sell_price), 0) AS revenue
			FROM %s AS b
			WHERE %s
			GROUP BY 1
		)
		SELECT
			s.bucket_start,
			COALESCE(a.revenue, 0) AS revenue
		FROM series s
		LEFT JOIN aggregated a ON a.bucket_start = s.bucket_start
		ORDER BY s.bucket_start ASC`,
		bucketConfig.BucketExpr,
		bookingTable,
		aggregateWhereClause,
	)

	var pointsResult []*mapper.BookingDashboardRevenuePointProjection
	if err := r.getDB(ctx).Raw(aggregateSQL, aggregateArgs...).Scan(&pointsResult).Error; err != nil {
		return nil, err
	}

	return &model.DashboardBookingRevenue{
		Granularity: bucketConfig.Granularity,
		FromDate:    minCheckIn,
		ToDate:      maxCheckIn,
		Chart:       mapper.ToDashboardBookingRevenuePoints(pointsResult),
	}, nil
}

func (r *bookingRepositoryImpl) FindDashboardRoomTypeRevenue(ctx context.Context, schema string, query model.DashboardOverviewQuery) (*model.DashboardRoomTypeRevenue, error) {
	roomTypeTable := r.schemaTable(schema, entity.RoomTypeTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	branchRoomTypeTable := r.schemaTable(schema, entity.BranchRoomTypeTableName)
	endDateExclusive := query.ToDate.AddDate(0, 0, 1)
	bookingBranchClause, bookingBranchArgs := dashboardBookingBranchClause(query.BranchIDs)
	roomTypeBranchClause := ""
	if len(query.BranchIDs) > 0 {
		roomTypeBranchClause = fmt.Sprintf(
			`WHERE EXISTS (SELECT 1 FROM %s brt WHERE brt.room_type_id = rt.id AND brt.branch_id IN ?)`,
			branchRoomTypeTable,
		)
	}

	sql := fmt.Sprintf(
		`SELECT
			rt.id AS room_type_id,
			rt.name AS room_type_name,
			COALESCE(br.revenue, 0) AS revenue
		FROM %s AS rt
		LEFT JOIN (
			SELECT
				b.room_type_id,
				COALESCE(SUM(b.total_sell_price), 0) AS revenue
			FROM %s AS b
			WHERE b.status IN (?, ?, ?)
				AND b.check_in < ? AND b.check_out > ?%s
			GROUP BY b.room_type_id
		) AS br ON br.room_type_id = rt.id
		%s
		ORDER BY revenue DESC, rt.name ASC`,
		roomTypeTable,
		bookingTable,
		bookingBranchClause,
		roomTypeBranchClause,
	)

	args := []any{
		model.BookingStatusConfirmed,
		model.BookingStatusCheckedIn,
		model.BookingStatusCheckedOut,
		endDateExclusive,
		query.FromDate,
	}
	args = append(args, bookingBranchArgs...)
	if len(query.BranchIDs) > 0 {
		args = append(args, query.BranchIDs)
	}

	var rows []*mapper.BookingDashboardRoomTypeRevenueProjection
	if err := r.getDB(ctx).Raw(
		sql,
		args...,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardRoomTypeRevenueModel(rows), nil
}

func (r *bookingRepositoryImpl) FindDashboardSourceRevenue(ctx context.Context, schema string, query model.DashboardOverviewQuery) (*model.DashboardSourceRevenue, error) {
	sourceTable := r.schemaTable(schema, entity.SourceTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	endDateExclusive := query.ToDate.AddDate(0, 0, 1)
	branchClause, branchArgs := dashboardBookingBranchClause(query.BranchIDs)

	sql := fmt.Sprintf(
		`SELECT
			s.id AS source_id,
			s.name AS source_name,
			COALESCE(br.revenue, 0) AS revenue,
			CASE
				WHEN SUM(COALESCE(br.revenue, 0)) OVER () = 0 THEN 0
				ELSE COALESCE(br.revenue, 0) * 100.0 / SUM(COALESCE(br.revenue, 0)) OVER ()
			END AS revenue_percent
		FROM %s AS s
		LEFT JOIN (
			SELECT
				b.source_id,
				COALESCE(SUM(b.total_sell_price), 0) AS revenue
			FROM %s AS b
			WHERE b.status IN (?, ?, ?)
				AND b.check_in < ? AND b.check_out > ?%s
			GROUP BY b.source_id
		) AS br ON br.source_id = s.id
		ORDER BY revenue DESC, s.name ASC`,
		sourceTable,
		bookingTable,
		branchClause,
	)

	args := []any{
		model.BookingStatusConfirmed,
		model.BookingStatusCheckedIn,
		model.BookingStatusCheckedOut,
		endDateExclusive,
		query.FromDate,
	}
	args = append(args, branchArgs...)

	var rows []*mapper.BookingDashboardSourceRevenueProjection
	if err := r.getDB(ctx).Raw(
		sql,
		args...,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardSourceRevenueModel(rows), nil
}

func (r *bookingRepositoryImpl) CountDashboardTodayCheckInOut(ctx context.Context, schema string, fromDate, toDate time.Time, branchIDs []int64) (int64, int64, error) {
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	branchClause, branchArgs := dashboardBookingBranchClause(branchIDs)

	sql := fmt.Sprintf(
		`SELECT
			(SELECT COUNT(*)
			 FROM %s AS b
			 WHERE b.check_in >= ? AND b.check_in < ?
				AND b.status <> ?%s) AS today_check_in,
			(SELECT COUNT(*)
			 FROM %s AS b
			 WHERE b.check_out >= ? AND b.check_out < ?
				AND b.status <> ?%s) AS today_check_out`,
		bookingTable,
		branchClause,
		bookingTable,
		branchClause,
	)

	args := []any{fromDate, toDate, model.BookingStatusCancelled}
	args = append(args, branchArgs...)
	args = append(args, fromDate, toDate, model.BookingStatusCancelled)
	args = append(args, branchArgs...)

	var result bookingDashboardTodayOperationResult
	if err := r.getDB(ctx).Raw(
		sql,
		args...,
	).Scan(&result).Error; err != nil {
		return 0, 0, err
	}

	return result.TodayCheckIn, result.TodayCheckOut, nil
}

func (r *bookingRepositoryImpl) FindTimeline(ctx context.Context, schema string, query model.BookingTimelineQuery) (*model.BookingTimeline, error) {
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)
	roomTypeTable := r.schemaTable(schema, entity.RoomTypeTableName)
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	sourceTable := r.schemaTable(schema, entity.SourceTableName)

	eventsSubQuery := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Select(
			"br.room_id, b.id AS booking_id, b.check_in, b.check_out, b.guest_full_name AS guest_name, b.status, b.booking_number, s.name AS source_name",
		).
		Joins(fmt.Sprintf("JOIN %s b ON b.id = br.booking_id", bookingTable)).
		Joins(fmt.Sprintf("JOIN %s s ON s.id = b.source_id", sourceTable)).
		Where(
			"(b.status = ? AND br.status = ?) OR "+
				"(b.status = ? AND br.status = ?) OR "+
				"(b.status = ? AND br.status = ?)",
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved,
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased,
			model.BookingStatusCheckedOut, model.BookingRoomStatusReleased,
		).
		Where("b.check_in < ? AND b.check_out > ?", query.ToDate, query.FromDate)

	db := r.getDB(ctx).
		Table(roomTable+" AS r").
		Select(
			"rt.id AS room_type_id, rt.name AS room_type_name, r.id AS room_id, r.name AS room_name, "+
				"e.booking_id, e.check_in, e.check_out, e.guest_name, e.source_name, e.status, e.booking_number",
		).
		Joins(fmt.Sprintf("JOIN %s rt ON rt.id = r.room_type_id", roomTypeTable)).
		Joins(fmt.Sprintf("JOIN %s f ON f.id = r.floor_id", floorTable)).
		Joins("LEFT JOIN (?) AS e ON e.room_id = r.id", eventsSubQuery)

	if query.RoomTypeID != 0 {
		db = db.Where("r.room_type_id = ?", query.RoomTypeID)
	}

	if query.BranchID != 0 {
		db = db.Where("f.branch_id = ?", query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		db = db.Where("f.branch_id IN ?", query.BranchIDs)
	}

	var rows []*mapper.BookingTimelineProjection
	if err := db.Order("rt.name ASC").
		Order("r.name ASC").
		Order("e.check_in ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToBookingTimelineModel(query.FromDate, query.ToDate, rows), nil
}

func buildDashboardRevenueWhereClause(query model.DashboardRevenueQuery) (string, []any) {
	whereClause := "b.status IN (?, ?, ?)"
	args := []any{
		model.BookingStatusConfirmed,
		model.BookingStatusCheckedIn,
		model.BookingStatusCheckedOut,
	}

	if query.RoomTypeID != 0 {
		whereClause += " AND b.room_type_id = ?"
		args = append(args, query.RoomTypeID)
	}

	if query.BranchID != 0 {
		whereClause += " AND b.branch_id = ?"
		args = append(args, query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		whereClause += " AND b.branch_id IN ?"
		args = append(args, query.BranchIDs)
	}

	return whereClause, args
}

func (r *bookingRepositoryImpl) countDashboardRooms(ctx context.Context, schema string, branchIDs []int64) (int64, error) {
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	db := r.getDB(ctx).
		Table(roomTable + " AS r").
		Model(&entity.RoomEntity{})

	if len(branchIDs) > 0 {
		db = db.
			Joins(fmt.Sprintf("JOIN %s f ON f.id = r.floor_id", floorTable)).
			Where("f.branch_id IN ?", branchIDs)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func dashboardBookingBranchClause(branchIDs []int64) (string, []any) {
	if len(branchIDs) == 0 {
		return "", nil
	}

	return " AND b.branch_id IN ?", []any{branchIDs}
}

func buildBookingRevenueBucketConfig(minCheckIn, maxCheckIn time.Time) bookingRevenueBucketConfig {
	spanDays := int(maxCheckIn.Sub(minCheckIn)/(24*time.Hour)) + 1

	granularity := chooseBookingRevenueGranularity(spanDays)
	switch granularity {
	case model.DashboardRevenueGranularityWeek:
		return bookingRevenueBucketConfig{
			Granularity: granularity,
			BucketExpr:  "date_trunc('week', b.check_in)::date",
			SeriesStart: truncateToISOWeek(minCheckIn),
			SeriesEnd:   truncateToISOWeek(maxCheckIn),
			Step:        "1 week",
		}
	case model.DashboardRevenueGranularityMonth:
		return bookingRevenueBucketConfig{
			Granularity: granularity,
			BucketExpr:  "date_trunc('month', b.check_in)::date",
			SeriesStart: time.Date(minCheckIn.Year(), minCheckIn.Month(), 1, 0, 0, 0, 0, time.UTC),
			SeriesEnd:   time.Date(maxCheckIn.Year(), maxCheckIn.Month(), 1, 0, 0, 0, 0, time.UTC),
			Step:        "1 month",
		}
	case model.DashboardRevenueGranularityYear:
		return bookingRevenueBucketConfig{
			Granularity: granularity,
			BucketExpr:  "date_trunc('year', b.check_in)::date",
			SeriesStart: time.Date(minCheckIn.Year(), time.January, 1, 0, 0, 0, 0, time.UTC),
			SeriesEnd:   time.Date(maxCheckIn.Year(), time.January, 1, 0, 0, 0, 0, time.UTC),
			Step:        "1 year",
		}
	default:
		return bookingRevenueBucketConfig{
			Granularity: model.DashboardRevenueGranularityDay,
			BucketExpr:  "b.check_in::date",
			SeriesStart: minCheckIn,
			SeriesEnd:   maxCheckIn,
			Step:        "1 day",
		}
	}
}

func chooseBookingRevenueGranularity(spanDays int) model.DashboardRevenueGranularity {
	switch {
	case spanDays <= 120:
		return model.DashboardRevenueGranularityDay
	case spanDays <= 720:
		return model.DashboardRevenueGranularityWeek
	case spanDays <= 3650:
		return model.DashboardRevenueGranularityMonth
	default:
		return model.DashboardRevenueGranularityYear
	}
}

func truncateToISOWeek(date time.Time) time.Time {
	normalized := utils.NormalizeDateOnly(date)
	weekday := int(normalized.Weekday())
	offset := (weekday + 6) % 7
	return normalized.AddDate(0, 0, -offset)
}

func (r *bookingRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.BookingQuery) ([]*model.Booking, int64, error) {
	var bookingEntities []*entity.BookingEntity
	var total int64

	orderSvcTable := r.schemaTable(schema, entity.BookingTableName)

	db := r.getDB(ctx).
		Table(orderSvcTable).
		Model(&entity.BookingEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Booking{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("Source", r.preloadWithSchema(schema, entity.SourceTableName))
	db = db.Preload("Branch", r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"))
	db = db.Preload("RoomType", r.preloadWithSchema(schema, entity.RoomTypeTableName))
	db = db.Preload("BookingRooms", r.preloadWithSchema(schema, entity.BookingRoomTableName))
	db = db.Preload("BookingRooms.Room", r.preloadWithSchema(schema, entity.RoomTableName))

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("id", "booking_number", "status", "guest_full_name", "guest_email", "guest_phone", "check_in", "check_out", "branch_id", "source_id", "room_type_id", "total_sell_price", "room_number", "predict").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&bookingEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToBookings(bookingEntities), total, nil
}

func (r *bookingRepositoryImpl) FindByIDWithBookingRooms(ctx context.Context, schema string, id int64) (*model.Booking, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Branch",
			Scope:    r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"),
		},
		Preload{
			Relation: "BookingRooms",
			Scope:    r.preloadWithSchema(schema, entity.BookingRoomTableName),
		},
		Preload{
			Relation: "BookingRooms.Room",
			Scope:    r.preloadWithSchema(schema, entity.RoomTableName, "id", "floor_id", "room_type_id"),
		},
		Preload{
			Relation: "BookingRooms.Room.Floor",
			Scope:    r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id"),
		},
	)
}

func (r *bookingRepositoryImpl) FindByBookingNumber(ctx context.Context, schema, bookingNumber string) (*model.Booking, error) {
	var bookingEntity entity.BookingEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingTableName)).
		Where("booking_number = ?", bookingNumber).
		First(&bookingEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToBooking(&bookingEntity), nil
}

func (r *bookingRepositoryImpl) FindByIDWithSourceAndBookingRooms(ctx context.Context, schema string, id int64) (*model.Booking, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Source",
			Scope:    r.preloadWithSchema(schema, entity.SourceTableName, "id", "name"),
		},
		Preload{
			Relation: "Branch",
			Scope:    r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"),
		},
		Preload{
			Relation: "BookingRooms",
			Scope:    r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "booking_id"),
		},
	)
}

func (r *bookingRepositoryImpl) FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.Booking, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Source",
			Scope:    r.preloadWithSchema(schema, entity.SourceTableName, "id", "name"),
		},
		Preload{
			Relation: "Branch",
			Scope:    r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"),
		},
		Preload{
			Relation: "RoomType",
			Scope:    r.preloadWithSchema(schema, entity.RoomTypeTableName, "id", "name"),
		},
		Preload{
			Relation: "BookingRooms",
			Scope:    r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "code", "room_id", "booking_id", "status"),
		},
		Preload{
			Relation: "BookingRooms.Room",
			Scope:    r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"),
		},
		Preload{
			Relation: "CreatedBy",
			Scope:    r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"),
		},
		Preload{
			Relation: "UpdatedBy",
			Scope:    r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"),
		},
	)
}

func (r *bookingRepositoryImpl) findByIDBase(tx *gorm.DB, schema string, id int64, preloads ...Preload) (*model.Booking, error) {
	var bookingEntity entity.BookingEntity

	tx = tx.Table(r.schemaTable(schema, entity.BookingTableName))

	for _, preload := range preloads {
		if preload.Scope != nil {
			tx = tx.Preload(preload.Relation, preload.Scope)
		} else {
			tx = tx.Preload(preload.Relation)
		}
	}

	if err := tx.Where("id = ?", id).
		First(&bookingEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToBooking(&bookingEntity), nil
}

func (r *bookingRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.BookingQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where(
			"booking_number ILIKE ? OR guest_full_name ILIKE ?",
			term, term,
		)
	}

	if query.SourceID != 0 {
		db = db.Where("source_id = ?", query.SourceID)
	}

	if query.BranchID != 0 {
		db = db.Where("branch_id = ?", query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		db = db.Where("branch_id IN ?", query.BranchIDs)
	}

	if query.UnassignedBranch != nil {
		if *query.UnassignedBranch {
			db = db.Where("branch_id IS NULL")
		} else {
			db = db.Where("branch_id IS NOT NULL")
		}
	}

	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("check_in >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("check_in < ?", query.ToDate.AddDate(0, 0, 1))
	}

	return db
}

func (r *bookingRepositoryImpl) applySorting(db *gorm.DB, query model.BookingQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"check_in":   "check_in",
		"check_out":  "check_out",
		"status":     "status",
		"created_at": "created_at",
	}

	sortField := "check_in"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}
