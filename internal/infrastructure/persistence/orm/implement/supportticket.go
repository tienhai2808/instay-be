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

type supportTicketRepositoryImpl struct {
	BaseRepository
}

type supportTicketDashboardOperationResult struct {
	PendingRequest    int64 `gorm:"column:pending_request"`
	RequestInProgress int64 `gorm:"column:request_in_progress"`
}

func NewSupportTicketRepository(baseRepo BaseRepository) repository.SupportTicketRepository {
	return &supportTicketRepositoryImpl{baseRepo}
}

func (r *supportTicketRepositoryImpl) Create(ctx context.Context, schema string, supportTicket *model.SupportTicket) error {
	supTickEntity := mapper.ToSupportTicketEntity(supportTicket)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)).
		Create(supTickEntity).Error
}

func (r *supportTicketRepositoryImpl) FindAllByBookingRoomIDWithSupportCategory(ctx context.Context, schema string, bookingRoomID int64) ([]*model.SupportTicket, error) {
	var supTicketEntities []*entity.SupportTicketEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)).
		Preload("SupportCategory", r.preloadWithSchema(schema, entity.SupportCategoryTableName, "id", "name")).
		Where("booking_room_id = ?", bookingRoomID).
		Find(&supTicketEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToSupportTickets(supTicketEntities), nil
}

func (r *supportTicketRepositoryImpl) FindByIDWithSupportCategory(ctx context.Context, schema string, id int64) (*model.SupportTicket, error) {
	var supportTicketEntity entity.SupportTicketEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)).
		Preload("SupportCategory", r.preloadWithSchema(schema, entity.SupportCategoryTableName, "id", "name", "department_id")).
		Where("id = ?", id).
		First(&supportTicketEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToSupportTicket(&supportTicketEntity), nil
}

func (r *supportTicketRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)).
		Model(&entity.SupportTicketEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrSupportTicketNotFound
	}

	return nil
}

func (r *supportTicketRepositoryImpl) FindAllWithSupportCategoryPaginated(ctx context.Context, schema string, departmentID *int64, query model.SupportTicketQuery) ([]*model.SupportTicket, int64, error) {
	var supportTicketEntities []*entity.SupportTicketEntity
	var total int64

	supTickTable := r.schemaTable(schema, entity.SupportTicketTableName)

	db := r.getDB(ctx).
		Table(supTickTable + " AS st").
		Model(&entity.SupportTicketEntity{})

	db = r.applyFilters(db, schema, departmentID, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.SupportTicket{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("SupportCategory", r.preloadWithSchema(schema, entity.SupportCategoryTableName, "id", "name"))
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "code", "room_id"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("st.*").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&supportTicketEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToSupportTickets(supportTicketEntities), total, nil
}

func (r *supportTicketRepositoryImpl) applyFilters(db *gorm.DB, schema string, departmentID *int64, query model.SupportTicketQuery) *gorm.DB {
	if departmentID != nil {
		supCatTable := r.schemaTable(schema, entity.SupportCategoryTableName)

		db = db.
			Joins(fmt.Sprintf(`INNER JOIN %s sc ON sc.id = st.support_category_id`, supCatTable)).
			Where("sc.department_id = ?", *departmentID)
	}

	if query.Status != "" {
		db = db.Where("st.status = ?", query.Status)
	}

	if query.Priority != "" {
		db = db.Where("st.priority = ?", query.Priority)
	}

	if len(query.BranchIDs) > 0 {
		bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
		roomTable := r.schemaTable(schema, entity.RoomTableName)
		floorTable := r.schemaTable(schema, entity.FloorTableName)

		db = db.
			Joins(fmt.Sprintf(`INNER JOIN %s br_scope ON br_scope.id = st.booking_room_id`, bookingRoomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s r_scope ON r_scope.id = br_scope.room_id`, roomTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s f_scope ON f_scope.id = r_scope.floor_id`, floorTable)).
			Where("f_scope.branch_id IN ?", query.BranchIDs)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("st.created_at >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("st.created_at < ?", query.ToDate.AddDate(0, 0, 1))
	}

	return db
}

func (r *supportTicketRepositoryImpl) applySorting(db *gorm.DB, query model.SupportTicketQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at": "st.created_at",
		"priority":   "st.priority",
	}

	sortField := "st.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}

func (r *supportTicketRepositoryImpl) CountDashboardPendingAndInProgress(ctx context.Context, schema string, branchIDs []int64) (int64, int64, error) {
	var result supportTicketDashboardOperationResult

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)+" AS st").
		Select(
			"COUNT(*) FILTER (WHERE st.status = ?) AS pending_request, "+
				"COUNT(*) FILTER (WHERE st.status = ?) AS request_in_progress",
			model.SupportTicketStatusPending,
			model.SupportTicketStatusAccepted,
		)

	db = r.applySupportTicketBranchFilter(db, schema, branchIDs)

	if err := db.Scan(&result).Error; err != nil {
		return 0, 0, err
	}

	return result.PendingRequest, result.RequestInProgress, nil
}

func (r *supportTicketRepositoryImpl) FindDashboardSupportTicketTime(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardSupportTicketTime, error) {
	var result mapper.SupportTicketDashboardTimeProjection

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.SupportTicketTableName)+" AS st").
		Select(
			"COUNT(*) FILTER (WHERE st.resolution_seconds IS NOT NULL AND st.resolution_seconds <= ?) AS under_or_equal_5_minutes, "+
				"COUNT(*) FILTER (WHERE st.resolution_seconds > ? AND st.resolution_seconds <= ?) AS under_or_equal_10_minutes, "+
				"COUNT(*) FILTER (WHERE st.resolution_seconds > ? AND st.resolution_seconds <= ?) AS under_or_equal_15_minutes, "+
				"COUNT(*) FILTER (WHERE st.resolution_seconds > ? AND st.resolution_seconds <= ?) AS under_or_equal_30_minutes, "+
				"COUNT(*) FILTER (WHERE st.resolution_seconds > ?) AS over_30_minutes, "+
				"COALESCE(AVG(st.resolution_seconds) FILTER (WHERE st.resolution_seconds IS NOT NULL AND st.resolution_seconds > 0), 0) AS average_resolution_seconds",
			5*60,
			5*60, 10*60,
			10*60, 15*60,
			15*60, 30*60,
			30*60,
		)

	db = r.applySupportTicketBranchFilter(db, schema, branchIDs)

	if err := db.Scan(&result).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardSupportTicketTimeModel(result), nil
}

func (r *supportTicketRepositoryImpl) applySupportTicketBranchFilter(db *gorm.DB, schema string, branchIDs []int64) *gorm.DB {
	if len(branchIDs) == 0 {
		return db
	}

	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)

	return db.
		Joins(fmt.Sprintf(`INNER JOIN %s br_dash ON br_dash.id = st.booking_room_id`, bookingRoomTable)).
		Joins(fmt.Sprintf(`INNER JOIN %s r_dash ON r_dash.id = br_dash.room_id`, roomTable)).
		Joins(fmt.Sprintf(`INNER JOIN %s f_dash ON f_dash.id = r_dash.floor_id`, floorTable)).
		Where("f_dash.branch_id IN ?", branchIDs)
}
