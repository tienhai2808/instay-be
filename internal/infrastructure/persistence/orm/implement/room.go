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
	"gorm.io/gorm"
)

type roomRepositoryImpl struct {
	BaseRepository
}

func NewRoomRepository(baseRepo BaseRepository) repository.RoomRepository {
	return &roomRepositoryImpl{baseRepo}
}

func (r *roomRepositoryImpl) Create(ctx context.Context, schema string, room *model.Room) error {
	roomEntity := mapper.ToRoomEntity(room)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Create(roomEntity).Error
}

func (r *roomRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.Room, error) {
	var roomEntity entity.RoomEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Where("id = ?", id).
		First(&roomEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToRoom(&roomEntity), nil
}

func (r *roomRepositoryImpl) FindByIDWithFloor(ctx context.Context, schema string, id int64) (*model.Room, error) {
	var roomEntity entity.RoomEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Preload("Floor", r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id")).
		Where("id = ?", id).
		First(&roomEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToRoom(&roomEntity), nil
}

func (r *roomRepositoryImpl) FindAllByRoomTypeID(ctx context.Context, schema string, roomTypeID int64) ([]*model.Room, error) {
	var roomEntities []*entity.RoomEntity

	roomTable := r.schemaTable(schema, entity.RoomTableName)
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)

	selectExpr := fmt.Sprintf(
		`r.id, r.name, r.room_type_id, NOT EXISTS (
			SELECT 1
			FROM %s br
			JOIN %s b ON b.id = br.booking_id
			WHERE br.room_id = r.id
				AND (
					(b.status = ? AND br.status = ?)
					OR (b.status = ? AND br.status = ?)
				)
		) AS is_available`,
		bookingRoomTable, bookingTable,
	)

	if err := r.getDB(ctx).
		Table(roomTable+" AS r").
		Select(
			selectExpr,
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved,
		).
		Where("r.room_type_id = ?", roomTypeID).
		Order("r.name ASC").
		Scan(&roomEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToRooms(roomEntities), nil
}

func (r *roomRepositoryImpl) ExistsByIDs(ctx context.Context, schema string, ids []int64) ([]int64, error) {
	if len(ids) == 0 {
		return make([]int64, 0), nil
	}

	existingIDs := make([]int64, 0, len(ids))
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Where("id IN ?", ids).
		Pluck("id", &existingIDs).Error; err != nil {
		return nil, err
	}

	return existingIDs, nil
}

func (r *roomRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.RoomQuery) ([]*model.Room, int64, error) {
	var roomEntities []*entity.RoomEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Model(&entity.RoomEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Room{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("RoomType", r.preloadWithSchema(schema, entity.RoomTypeTableName, "id", "name"))
	db = db.Preload("Floor", r.preloadWithSchema(schema, entity.FloorTableName, "id", "name", "branch_id"))
	db = db.Preload("Floor.Branch", r.preloadWithSchema(schema, entity.BranchTableName, "id", "name"))
	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&roomEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToRooms(roomEntities), total, nil
}

func (r *roomRepositoryImpl) CountStatusSummary(ctx context.Context, schema string, query model.RoomStatusQuery) (*model.RoomStatusSummary, error) {
	var summary mapper.RoomStatusSummaryProjection

	roomTable := r.schemaTable(schema, entity.RoomTableName)
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	selectExpr := `COUNT(*) AS all_count,
		COUNT(*) FILTER (WHERE room_status = 'available') AS available_count,
		COUNT(*) FILTER (WHERE room_status = 'occupied') AS occupied_count,
		COUNT(*) FILTER (WHERE room_status = 'booked') AS booked_count`

	stateExpr := fmt.Sprintf(
		`CASE
			WHEN EXISTS (
				SELECT 1
				FROM %s br
				JOIN %s b ON b.id = br.booking_id
				WHERE br.room_id = r.id
					AND b.status = ?
					AND br.status = ?
					AND b.check_out >= ?
			) THEN 'occupied'
			WHEN EXISTS (
				SELECT 1
				FROM %s br
				JOIN %s b ON b.id = br.booking_id
				WHERE br.room_id = r.id
					AND b.status = ?
					AND br.status = ?
					AND b.check_in <= ?
					AND b.check_out >= ?
			) THEN 'booked'
			ELSE 'available'
		END`,
		bookingRoomTable, bookingTable,
		bookingRoomTable, bookingTable,
	)

	args := []any{
		model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, query.SnapshotAt,
		model.BookingStatusConfirmed, model.BookingRoomStatusReserved, query.SnapshotAt, query.SnapshotAt,
	}

	db := r.getDB(ctx).
		Table(roomTable+" AS r").
		Select(stateExpr+" AS room_status", args...)

	if query.RoomTypeID != 0 {
		db = db.Where("r.room_type_id = ?", query.RoomTypeID)
	}

	if query.BranchID != 0 {
		db = db.Joins(fmt.Sprintf("JOIN %s f ON f.id = r.floor_id", floorTable)).
			Where("f.branch_id = ?", query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		db = db.Joins(fmt.Sprintf("JOIN %s f ON f.id = r.floor_id", floorTable)).
			Where("f.branch_id IN ?", query.BranchIDs)
	}

	if err := r.getDB(ctx).
		Table("(?) AS room_states", db).
		Select(selectExpr).
		Scan(&summary).Error; err != nil {
		return nil, err
	}

	return mapper.ToRoomStatusSummaryModel(summary), nil
}

func (r *roomRepositoryImpl) FindTimeline(ctx context.Context, schema string, query model.RoomTimelineQuery) (*model.RoomTimeline, error) {
	snapshotAt := query.SnapshotAt
	if snapshotAt.IsZero() {
		snapshotAt = time.Now()
	}

	upcomingThreshold := query.UpcomingThreshold
	if upcomingThreshold <= 0 {
		upcomingThreshold = time.Hour
	}
	upcomingTo := snapshotAt.Add(upcomingThreshold)

	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)
	roomTypeTable := r.schemaTable(schema, entity.RoomTypeTableName)
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	sourceTable := r.schemaTable(schema, entity.SourceTableName)

	candidateQuery := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Select(
			"br.room_id AS room_id, b.id AS booking_id, b.booking_number AS booking_number, b.status AS booking_status, "+
				"s.name AS source_name, b.guest_full_name AS guest_name, b.check_in AS check_in, b.check_out AS check_out, "+
				"CASE "+
				"WHEN b.status = ? AND br.status = ? AND b.check_out < ? THEN 'overdue_checkout' "+
				"WHEN b.status = ? AND br.status = ? AND b.check_out >= ? THEN 'occupied' "+
				"WHEN b.status = ? AND br.status = ? AND b.check_in <= ? THEN 'overdue_checkin' "+
				"WHEN b.status = ? AND br.status = ? AND b.check_in > ? AND b.check_in <= ? THEN 'upcoming_checkin' "+
				"END AS room_status, "+
				"CASE "+
				"WHEN b.status = ? AND br.status = ? AND b.check_out < ? THEN 1 "+
				"WHEN b.status = ? AND br.status = ? AND b.check_out >= ? THEN 2 "+
				"WHEN b.status = ? AND br.status = ? AND b.check_in <= ? THEN 3 "+
				"WHEN b.status = ? AND br.status = ? AND b.check_in > ? AND b.check_in <= ? THEN 4 "+
				"ELSE 5 "+
				"END AS status_priority, "+
				"CASE WHEN b.status = ? THEN b.check_out ELSE b.check_in END AS status_time",
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, snapshotAt,
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, snapshotAt,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved, snapshotAt,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved, snapshotAt, upcomingTo,
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, snapshotAt,
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, snapshotAt,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved, snapshotAt,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved, snapshotAt, upcomingTo,
			model.BookingStatusCheckedIn,
		).
		Joins(fmt.Sprintf("JOIN %s b ON b.id = br.booking_id", bookingTable)).
		Joins(fmt.Sprintf("JOIN %s s ON s.id = b.source_id", sourceTable)).
		Where(
			"(b.status = ? AND br.status = ? AND b.check_out < ?) "+
				"OR (b.status = ? AND br.status = ? AND b.check_out >= ?) "+
				"OR (b.status = ? AND br.status = ? AND b.check_in <= ?) "+
				"OR (b.status = ? AND br.status = ? AND b.check_in > ? AND b.check_in <= ?)",
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, snapshotAt,
			model.BookingStatusCheckedIn, model.BookingRoomStatusReleased, snapshotAt,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved, snapshotAt,
			model.BookingStatusConfirmed, model.BookingRoomStatusReserved, snapshotAt, upcomingTo,
		)

	rankedCandidateQuery := r.getDB(ctx).
		Table("(?) AS candidates", candidateQuery).
		Select(
			"room_id, booking_id, booking_number, booking_status, source_name, guest_name, check_in, check_out, room_status, " +
				"ROW_NUMBER() OVER (PARTITION BY room_id ORDER BY status_priority ASC, status_time ASC, booking_id ASC) AS row_num",
		)

	db := r.getDB(ctx).
		Table(roomTable+" AS r").
		Select(
			"rt.id AS room_type_id, rt.name AS room_type_name, rt.price AS room_price, "+
				"r.id AS room_id, r.name AS room_name, COALESCE(rc.room_status, 'available') AS room_status, "+
				"rc.booking_id, rc.booking_number, rc.booking_status, rc.source_name, rc.guest_name, rc.check_in, rc.check_out",
		).
		Joins(fmt.Sprintf("JOIN %s rt ON rt.id = r.room_type_id", roomTypeTable)).
		Joins(fmt.Sprintf("JOIN %s f ON f.id = r.floor_id", floorTable)).
		Joins("LEFT JOIN (?) AS rc ON rc.room_id = r.id AND rc.row_num = 1", rankedCandidateQuery)

	if query.RoomTypeID != 0 {
		db = db.Where("r.room_type_id = ?", query.RoomTypeID)
	}

	if query.BranchID != 0 {
		db = db.Where("f.branch_id = ?", query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		db = db.Where("f.branch_id IN ?", query.BranchIDs)
	}

	if query.FloorID != 0 {
		db = db.Where("r.floor_id = ?", query.FloorID)
	}

	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where(
			"r.name ILIKE ? OR COALESCE(rc.guest_name, '') ILIKE ? OR COALESCE(rc.booking_number, '') ILIKE ? OR COALESCE(rc.source_name, '') ILIKE ?",
			term, term, term, term,
		)
	}

	var rows []*mapper.RoomTimelineProjection
	if err := db.Order("rt.name ASC").
		Order("r.name ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToRoomTimelineModel(snapshotAt, rows), nil
}

func (r *roomRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Model(&entity.RoomEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrRoomNotFound
	}

	return nil
}

func (r *roomRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Where("id = ?", id).
		Delete(&entity.RoomEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrRoomNotFound
	}

	return nil
}

func (r *roomRepositoryImpl) DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error) {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Where("id IN ?", ids).
		Delete(&entity.RoomEntity{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *roomRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.RoomQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	if query.BranchID != 0 {
		floorTable := r.schemaTable(schema, entity.FloorTableName)
		db = db.Where(fmt.Sprintf("floor_id IN (SELECT id FROM %s WHERE branch_id = ?)", floorTable), query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		floorTable := r.schemaTable(schema, entity.FloorTableName)
		db = db.Where(fmt.Sprintf("floor_id IN (SELECT id FROM %s WHERE branch_id IN ?)", floorTable), query.BranchIDs)
	}

	if query.FloorID != 0 {
		db = db.Where("floor_id = ?", query.FloorID)
	}

	if query.RoomTypeID != 0 {
		db = db.Where("room_type_id = ?", query.RoomTypeID)
	}

	return db
}

func (r *roomRepositoryImpl) applySorting(db *gorm.DB, query model.RoomQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"name":       "name",
		"created_at": "created_at",
	}

	sortField := "created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}
