package implement

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type chatRepositoryImpl struct {
	BaseRepository
}

func NewChatRepository(baseRepo BaseRepository) repository.ChatRepository {
	return &chatRepositoryImpl{baseRepo}
}

func (r *chatRepositoryImpl) Create(ctx context.Context, schema string, chat *model.Chat) error {
	chatEntity := mapper.ToChatEntity(chat)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ChatTableName)).
		Create(chatEntity).Error
}

func (r *chatRepositoryImpl) FindAllByBookingRoomIDWithDepartmentPaginated(ctx context.Context, schema string, bookingRoomID int64, query model.ChatQuery) ([]*model.Chat, int64, error) {
	var chatEntities []*entity.ChatEntity
	var total int64

	chatTable := r.schemaTable(schema, entity.ChatTableName)
	db := r.getDB(ctx).
		Table(chatTable+" AS c").
		Model(&entity.ChatEntity{}).
		Where("c.booking_room_id = ?", bookingRoomID)

	db = r.applyGuestFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Chat{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = db.Preload("Department", r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("c.*").
		Order("COALESCE(c.last_message_at, c.created_at) DESC").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&chatEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToChats(chatEntities), total, nil
}

func (r *chatRepositoryImpl) FindAllWithBookingRoomDetailsPaginated(ctx context.Context, schema string, departmentID *int64, query model.ChatQuery) ([]*model.Chat, int64, error) {
	var chatEntities []*entity.ChatEntity
	var total int64

	chatTable := r.schemaTable(schema, entity.ChatTableName)
	db := r.getDB(ctx).
		Table(chatTable + " AS c").
		Model(&entity.ChatEntity{})

	db = r.applyAdminFilters(db, schema, departmentID, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Chat{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "booking_id", "room_id"))
	db = db.Preload("BookingRoom.Booking", r.preloadWithSchema(schema, entity.BookingTableName, "id", "booking_number", "guest_full_name", "guest_email", "guest_phone", "check_in", "check_out"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))
	db = db.Preload("BookingRoom.Room.Floor", r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("c.*").
		Order("COALESCE(c.last_message_at, c.created_at) DESC").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&chatEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToChats(chatEntities), total, nil
}

func (r *chatRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ChatTableName)).
		Model(&entity.ChatEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrChatNotFound
	}

	return nil
}

func (r *chatRepositoryImpl) FindByIDWithDepartment(ctx context.Context, schema string, id int64) (*model.Chat, error) {
	var chatEntity entity.ChatEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ChatTableName)).
		Preload("Department", r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name")).
		Where("id = ?", id).
		First(&chatEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToChat(&chatEntity), nil
}

func (r *chatRepositoryImpl) FindByIDWithBookingRoomDetails(ctx context.Context, schema string, id int64) (*model.Chat, error) {
	var chatEntity entity.ChatEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ChatTableName)).
		Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "booking_id", "room_id")).
		Preload("BookingRoom.Booking", r.preloadWithSchema(schema, entity.BookingTableName, "id", "booking_number", "guest_full_name", "guest_email", "guest_phone", "check_in", "check_out")).
		Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name")).
		Preload("BookingRoom.Room.Floor", r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id")).
		Where("id = ?", id).
		First(&chatEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToChat(&chatEntity), nil
}

func (r *chatRepositoryImpl) FindByBookingRoomIDAndDepartmentID(ctx context.Context, schema string, bookingRoomID, departmentID int64) (*model.Chat, error) {
	var chatEntity entity.ChatEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ChatTableName)).
		Where("booking_room_id = ? AND department_id = ?", bookingRoomID, departmentID).
		First(&chatEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToChat(&chatEntity), nil
}

func (r *chatRepositoryImpl) FindIDsByBranchIDs(ctx context.Context, schema string, branchIDs []int64) ([]int64, error) {
	if len(branchIDs) == 0 {
		return make([]int64, 0), nil
	}

	chatTable := r.schemaTable(schema, entity.ChatTableName)
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	ids := make([]int64, 0)
	if err := r.getDB(ctx).
		Table(chatTable+" AS c").
		Select("c.id").
		Joins(fmt.Sprintf(`JOIN %s br ON br.id = c.booking_room_id`, bookingRoomTable)).
		Joins(fmt.Sprintf(`JOIN %s r ON r.id = br.room_id`, roomTable)).
		Joins(fmt.Sprintf(`JOIN %s f ON f.id = r.floor_id`, floorTable)).
		Where("f.branch_id IN ?", branchIDs).
		Scan(&ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *chatRepositoryImpl) FindAllExternalBySinceAndCursor(ctx context.Context, schema string, since, cursorAt *time.Time, cursorID *int64, limit uint32) ([]*model.Chat, bool, error) {
	if limit == 0 {
		limit = 100
	}
	if limit > 100 {
		limit = 100
	}

	fetchLimit := int(limit) + 1
	sortExpr := "COALESCE(c.last_message_at, c.created_at)"

	var chatEntities []*entity.ChatEntity
	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ChatTableName) + " AS c").
		Model(&entity.ChatEntity{})

	if since != nil {
		db = db.Where(sortExpr+" >= ?", *since)
	}

	if cursorAt != nil && cursorID != nil {
		db = db.Where(
			"("+sortExpr+" < ? OR ("+sortExpr+" = ? AND c.id < ?))",
			*cursorAt,
			*cursorAt,
			*cursorID,
		)
	}

	db = db.Session(&gorm.Session{})
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "booking_id", "room_id"))
	db = db.Preload("BookingRoom.Booking", r.preloadWithSchema(schema, entity.BookingTableName, "id", "booking_number"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))

	if err := db.
		Select("c.*").
		Order(sortExpr + " DESC").
		Order("c.id DESC").
		Limit(fetchLimit).
		Find(&chatEntities).Error; err != nil {
		return nil, false, err
	}

	hasNext := len(chatEntities) > int(limit)
	if hasNext {
		chatEntities = chatEntities[:int(limit)]
	}

	return mapper.ToChats(chatEntities), hasNext, nil
}

func (r *chatRepositoryImpl) applyGuestFilters(db *gorm.DB, schema string, query model.ChatQuery) *gorm.DB {
	if query.Search != "" {
		deptTable := r.schemaTable(schema, entity.DepartmentTableName)
		term := "%" + query.Search + "%"
		db = db.
			Joins(fmt.Sprintf("INNER JOIN %s d ON d.id = c.department_id", deptTable)).
			Where("d.name ILIKE ?", term)
	}

	return db
}

func (r *chatRepositoryImpl) applyAdminFilters(db *gorm.DB, schema string, departmentID *int64, query model.ChatQuery) *gorm.DB {
	if departmentID != nil {
		db = db.
			Where("c.department_id = ?", *departmentID).
			Where("c.expired_at > CURRENT_TIMESTAMP")
	}

	if query.Search != "" {
		bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
		bookingTable := r.schemaTable(schema, entity.BookingTableName)
		roomTable := r.schemaTable(schema, entity.RoomTableName)
		term := "%" + query.Search + "%"

		db = db.
			Joins(fmt.Sprintf("INNER JOIN %s br ON br.id = c.booking_room_id", bookingRoomTable)).
			Joins(fmt.Sprintf("INNER JOIN %s b ON b.id = br.booking_id", bookingTable)).
			Joins(fmt.Sprintf("INNER JOIN %s r ON r.id = br.room_id", roomTable)).
			Where("(b.booking_number ILIKE ? OR r.name ILIKE ?)", term, term)
	}

	if len(query.BranchIDs) > 0 {
		bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
		roomTable := r.schemaTable(schema, entity.RoomTableName)
		floorTable := r.schemaTable(schema, entity.FloorTableName)

		db = db.
			Joins(fmt.Sprintf(`INNER JOIN %s br_scope ON br_scope.id = c.booking_room_id`, bookingRoomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s r_scope ON r_scope.id = br_scope.room_id`, roomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s f_scope ON f_scope.id = r_scope.floor_id`, floorTable)).
			Where("f_scope.branch_id IN ?", query.BranchIDs)
	}

	return db
}
