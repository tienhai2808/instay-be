package implement

import (
	"context"
	"errors"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type bookingRoomRepositoryImpl struct {
	BaseRepository
}

func NewBookingRoomRepository(baseRepo BaseRepository) repository.BookingRoomRepository {
	return &bookingRoomRepositoryImpl{baseRepo}
}

func (r *bookingRoomRepositoryImpl) Create(ctx context.Context, schema string, bookingRoom *model.BookingRoom) error {
	bookingRoomEntity := mapper.ToBookingRoomEntity(bookingRoom)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingRoomTableName)).
		Create(&bookingRoomEntity).Error
}

func (r *bookingRoomRepositoryImpl) CreateBatch(ctx context.Context, schema string, bookingRooms []*model.BookingRoom) error {
	bookingRoomEntities := mapper.ToBookingRoomEntities(bookingRooms)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingRoomTableName)).
		Create(&bookingRoomEntities).Error
}

func (r *bookingRoomRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingRoomTableName)).
		Model(&entity.BookingRoomEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrBookingRoomNotFound
	}

	return nil
}

func (r *bookingRoomRepositoryImpl) Updates(ctx context.Context, schema string, ids []int64, updateData map[string]any) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingRoomTableName)).
		Model(&entity.BookingRoomEntity{}).
		Where("id IN ?", ids).
		Updates(updateData).Error
}

func (r *bookingRoomRepositoryImpl) DeleteAllByBookingID(ctx context.Context, schema string, bookingID int64) error {
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BookingRoomTableName)).
		Where("booking_id = ?", bookingID).
		Delete(&entity.BookingRoomEntity{}).Error
}

func (r *bookingRoomRepositoryImpl) FindByIDWithRoom(ctx context.Context, schema string, id int64) (*model.BookingRoom, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Room",
			Scope:    r.preloadWithSchema(schema, entity.RoomTableName),
		},
		Preload{
			Relation: "Room.Floor",
			Scope:    r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id"),
		},
	)
}

func (r *bookingRoomRepositoryImpl) FindAllActiveIDs(ctx context.Context, schema string) ([]int64, error) {
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)

	ids := make([]int64, 0)
	if err := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Joins("JOIN "+bookingTable+" AS b ON b.id = br.booking_id").
		Where("br.status = ?", model.BookingRoomStatusReleased).
		Where("b.status = ?", model.BookingStatusCheckedIn).
		Where("b.check_out > NOW()").
		Pluck("br.id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *bookingRoomRepositoryImpl) FindAllActiveIDsByBranchIDs(ctx context.Context, schema string, branchIDs []int64) ([]int64, error) {
	if len(branchIDs) == 0 {
		return make([]int64, 0), nil
	}

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	bookingTable := r.schemaTable(schema, entity.BookingTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	ids := make([]int64, 0)
	if err := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Joins("JOIN "+bookingTable+" AS b ON b.id = br.booking_id").
		Joins("JOIN "+roomTable+" AS r ON r.id = br.room_id").
		Joins("JOIN "+floorTable+" AS f ON f.id = r.floor_id").
		Where("br.status = ?", model.BookingRoomStatusReleased).
		Where("b.status = ?", model.BookingStatusCheckedIn).
		Where("b.check_out > NOW()").
		Where("f.branch_id IN ?", branchIDs).
		Pluck("br.id", &ids).Error; err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *bookingRoomRepositoryImpl) FindByIDWithBookingAndRoom(ctx context.Context, schema string, id int64) (*model.BookingRoom, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Room",
			Scope:    r.preloadWithSchema(schema, entity.RoomTableName),
		},
		Preload{
			Relation: "Room.Floor",
			Scope:    r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id"),
		},
		Preload{
			Relation: "Booking",
			Scope:    r.preloadWithSchema(schema, entity.BookingTableName, "id", "booking_number", "guest_full_name", "guest_phone", "guest_email", "check_in", "check_out", "status"),
		},
	)
}

func (r *bookingRoomRepositoryImpl) FindByIDWithBooking(ctx context.Context, schema string, id int64) (*model.BookingRoom, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Booking",
			Scope:    r.preloadWithSchema(schema, entity.BookingTableName, "id", "status", "check_in", "check_out"),
		},
	)
}

func (r *bookingRoomRepositoryImpl) ExistsActiveByRoomIDExceptID(ctx context.Context, schema string, id, roomID int64) (bool, error) {
	var count int64

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)

	if err := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Where("br.room_id = ?", roomID).
		Where("br.id <> ?", id).
		Where("br.status = ?", model.BookingRoomStatusReserved).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *bookingRoomRepositoryImpl) ExistsActiveByRoomID(ctx context.Context, schema string, roomID int64) (bool, error) {
	var count int64

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)

	if err := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Where("br.room_id = ?", roomID).
		Where("br.status = ?", model.BookingRoomStatusReserved).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *bookingRoomRepositoryImpl) FindActiveRoomIDsByRoomIDs(ctx context.Context, schema string, roomIDs []int64) ([]int64, error) {
	if len(roomIDs) == 0 {
		return make([]int64, 0), nil
	}

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)

	activeRoomIDs := make([]int64, 0, len(roomIDs))
	if err := r.getDB(ctx).
		Table(bookingRoomTable+" AS br").
		Where("br.room_id IN ?", roomIDs).
		Where("br.status = ?", model.BookingRoomStatusReserved).
		Distinct("br.room_id").
		Pluck("br.room_id", &activeRoomIDs).Error; err != nil {
		return nil, err
	}

	return activeRoomIDs, nil
}

func (r *bookingRoomRepositoryImpl) FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.BookingRoom, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Room",
			Scope:    r.preloadWithSchema(schema, entity.RoomTableName),
		},
		Preload{
			Relation: "Room.Floor",
			Scope:    r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id"),
		},
		Preload{
			Relation: "OrderServices",
			Scope:    r.preloadWithSchema(schema, entity.OrderServiceTableName),
		},
		Preload{
			Relation: "OrderServices.Service",
			Scope:    r.preloadWithSchema(schema, entity.ServiceTableName),
		},
		Preload{
			Relation: "SupportTickets",
			Scope:    r.preloadWithSchema(schema, entity.SupportTicketTableName),
		},
		Preload{
			Relation: "SupportTickets.SupportCategory",
			Scope:    r.preloadWithSchema(schema, entity.SupportCategoryTableName),
		},
		Preload{
			Relation: "Guests",
			Scope:    r.preloadWithSchema(schema, entity.GuestTableName),
		},
		Preload{
			Relation: "Review",
			Scope:    r.preloadWithSchema(schema, entity.ReviewTableName),
		},
	)
}

func (r *bookingRoomRepositoryImpl) findByIDBase(tx *gorm.DB, schema string, id int64, preloads ...Preload) (*model.BookingRoom, error) {
	var bookingRoomEntity entity.BookingRoomEntity

	tx = tx.Table(r.schemaTable(schema, entity.BookingRoomTableName))

	for _, preload := range preloads {
		if preload.Scope != nil {
			tx = tx.Preload(preload.Relation, preload.Scope)
		} else {
			tx = tx.Preload(preload.Relation)
		}
	}

	if err := tx.Where("id = ?", id).
		First(&bookingRoomEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToBookingRoom(&bookingRoomEntity), nil
}
