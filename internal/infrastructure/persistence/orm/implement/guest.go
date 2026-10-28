package implement

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type guestRepositoryImpl struct {
	BaseRepository
}

func NewGuestRepository(baseRepo BaseRepository) repository.GuestRepository {
	return &guestRepositoryImpl{baseRepo}
}

func (r *guestRepositoryImpl) Create(ctx context.Context, schema string, guest *model.Guest) error {
	guestEntity := mapper.ToGuestEntity(guest)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.GuestTableName)).
		Create(guestEntity).Error
}

func (r *guestRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.GuestTableName)).
		Model(&entity.GuestEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrGuestNotFound
	}

	return nil
}

func (r *guestRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.GuestTableName)).
		Where("id = ?", id).
		Delete(&entity.GuestEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrGuestNotFound
	}

	return nil
}

func (r *guestRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.GuestQuery) ([]*model.Guest, int64, error) {
	var guestEntities []*entity.GuestEntity
	var total int64

	guestTable := r.schemaTable(schema, entity.GuestTableName)

	db := r.getDB(ctx).
		Table(guestTable).
		Model(&entity.GuestEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Guest{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "room_id"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&guestEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToGuests(guestEntities), total, nil
}

func (r *guestRepositoryImpl) FindByIDWithBookingRoomDetails(ctx context.Context, schema string, id int64) (*model.Guest, error) {
	var guestEntity entity.GuestEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.GuestTableName)).
		Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "room_id")).
		Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "floor_id")).
		Preload("BookingRoom.Room.Floor", r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id")).
		Where("id = ?", id).
		First(&guestEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToGuest(&guestEntity), nil
}

func (r *guestRepositoryImpl) FindAllWithDetails(ctx context.Context, schema string, query model.GuestQuery) ([]*model.Guest, error) {
	var guestEntities []*entity.GuestEntity

	guestTable := r.schemaTable(schema, entity.GuestTableName)

	db := r.getDB(ctx).
		Table(guestTable).
		Model(&entity.GuestEntity{})

	db = r.applyFilters(db, schema, query)
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "room_id"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))

	db = r.applySorting(db, query)

	if err := db.Find(&guestEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToGuests(guestEntities), nil
}

func (r *guestRepositoryImpl) FindAllByBookingRoomID(ctx context.Context, schema string, bookingRoomID int64) ([]*model.Guest, error) {
	var guestEntities []*entity.GuestEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.GuestTableName)).
		Where("booking_room_id = ?", bookingRoomID).
		Find(&guestEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToGuests(guestEntities), nil
}

func (r *guestRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.GuestQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where(
			"full_name ILIKE ? OR passport_number ILIKE ?",
			term, term,
		)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("check_in >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("check_in < ?", query.ToDate.AddDate(0, 0, 1))
	}

	if len(query.BranchIDs) > 0 {
		bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
		roomTable := r.schemaTable(schema, entity.RoomTableName)
		floorTable := r.schemaTable(schema, entity.FloorTableName)

		db = db.
			Joins(fmt.Sprintf(`INNER JOIN %s br_scope ON br_scope.id = booking_room_id`, bookingRoomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s r_scope ON r_scope.id = br_scope.room_id`, roomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s f_scope ON f_scope.id = r_scope.floor_id`, floorTable)).
			Where("f_scope.branch_id IN ?", query.BranchIDs)
	}

	return db
}

func (r *guestRepositoryImpl) applySorting(db *gorm.DB, query model.GuestQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"check_in":  "check_in",
		"check_out": "check_out",
		"national":  "national",
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

func (r *guestRepositoryImpl) FindDashboardGuestNational(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardGuestNational, error) {
	guestTable := r.schemaTable(schema, entity.GuestTableName)

	var rows []*mapper.GuestDashboardNationalProjection
	db := r.getDB(ctx).
		Table(guestTable + " AS g").
		Select(
			"NULLIF(TRIM(g.national), '') AS national, COUNT(*) AS guest_count",
		).
		Where("NULLIF(TRIM(g.national), '') IS NOT NULL")

	if len(branchIDs) > 0 {
		bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
		roomTable := r.schemaTable(schema, entity.RoomTableName)
		floorTable := r.schemaTable(schema, entity.FloorTableName)

		db = db.
			Joins(fmt.Sprintf(`INNER JOIN %s br_dash ON br_dash.id = g.booking_room_id`, bookingRoomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s r_dash ON r_dash.id = br_dash.room_id`, roomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s f_dash ON f_dash.id = r_dash.floor_id`, floorTable)).
			Where("f_dash.branch_id IN ?", branchIDs)
	}

	if err := db.
		Group("NULLIF(TRIM(g.national), '')").
		Order("guest_count DESC, national ASC").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardGuestNationalModel(rows), nil
}
