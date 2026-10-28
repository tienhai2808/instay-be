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

type orderServiceRepositoryImpl struct {
	BaseRepository
}

func NewOrderServiceRepository(baseRepo BaseRepository) repository.OrderServiceRepository {
	return &orderServiceRepositoryImpl{baseRepo}
}

func (r *orderServiceRepositoryImpl) Create(ctx context.Context, schema string, orderService *model.OrderService) error {
	orderSvcEntity := mapper.ToOrderServiceEntity(orderService)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.OrderServiceTableName)).
		Create(orderSvcEntity).Error
}

func (r *orderServiceRepositoryImpl) FindAllByBookingRoomIDWithServiceDetails(ctx context.Context, schema string, bookingRoomID int64) ([]*model.OrderService, error) {
	var orderServiceEntities []*entity.OrderServiceEntity
	serviceImageTable := r.schemaTable(schema, entity.ServiceImageTableName)

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.OrderServiceTableName)).
		Preload("Service", r.preloadWithSchema(schema, entity.ServiceTableName)).
		Preload("Service.ServiceImages", func(db *gorm.DB) *gorm.DB {
			return db.
				Table(serviceImageTable).
				Where("is_thumbnail = true")
		}).
		Where("booking_room_id = ?", bookingRoomID).
		Find(&orderServiceEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToOrderServices(orderServiceEntities), nil
}

func (r *orderServiceRepositoryImpl) FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.OrderService, error) {
	var orderServiceEntity entity.OrderServiceEntity
	serviceImageTable := r.schemaTable(schema, entity.ServiceImageTableName)

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.OrderServiceTableName)).
		Preload("Service", r.preloadWithSchema(schema, entity.ServiceTableName)).
		Preload("Service.ServiceImages", func(db *gorm.DB) *gorm.DB {
			return db.
				Table(serviceImageTable).
				Where("is_thumbnail = true")
		}).
		Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName)).
		Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName)).
		Preload("BookingRoom.Room.Floor", r.preloadWithSchema(schema, entity.FloorTableName, "id", "branch_id")).
		Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name")).
		Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name")).
		Where("id = ?", id).
		First(&orderServiceEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToOrderService(&orderServiceEntity), nil
}

func (r *orderServiceRepositoryImpl) FindAllWithServicePaginated(ctx context.Context, schema string, departmentID *int64, query model.OrderServiceQuery) ([]*model.OrderService, int64, error) {
	var orderServiceEntities []*entity.OrderServiceEntity
	var total int64

	orderSvcTable := r.schemaTable(schema, entity.OrderServiceTableName)

	db := r.getDB(ctx).
		Table(orderSvcTable + " AS os").
		Model(&entity.OrderServiceEntity{})

	db = r.applyFilters(db, schema, departmentID, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.OrderService{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("Service", r.preloadWithSchema(schema, entity.ServiceTableName, "id", "name", "slug"))
	db = db.Preload("BookingRoom", r.preloadWithSchema(schema, entity.BookingRoomTableName, "id", "booking_id", "room_id"))
	db = db.Preload("BookingRoom.Room", r.preloadWithSchema(schema, entity.RoomTableName, "id", "name"))

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("os.id", "os.status", "os.guest_note", "os.service_id", "os.booking_room_id", "os.estimated_at", "os.quantity", "os.unit_price", "os.total_price").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&orderServiceEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToOrderServices(orderServiceEntities), total, nil
}

func (r *orderServiceRepositoryImpl) FindByIDWithServiceDetails(ctx context.Context, schema string, id int64) (*model.OrderService, error) {
	var orderServiceEntity entity.OrderServiceEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.OrderServiceTableName)).
		Preload("Service", r.preloadWithSchema(schema, entity.ServiceTableName, "id", "name", "slug", "service_type_id")).
		Preload("Service.ServiceType", r.preloadWithSchema(schema, entity.ServiceTypeTableName, "id", "department_id")).
		Where("id = ?", id).
		First(&orderServiceEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToOrderService(&orderServiceEntity), nil
}

func (r *orderServiceRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.OrderServiceTableName)).
		Model(&entity.OrderServiceEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrOrderServiceNotFound
	}

	return nil
}

func (r *orderServiceRepositoryImpl) applyFilters(db *gorm.DB, schema string, departmentID *int64, query model.OrderServiceQuery) *gorm.DB {
	if departmentID != nil {
		serviceTable := r.schemaTable(schema, entity.ServiceTableName)
		serviceTypeTable := r.schemaTable(schema, entity.ServiceTypeTableName)

		db = db.
			Joins(fmt.Sprintf(`INNER JOIN %s svc ON svc.id = os.service_id`, serviceTable)).
			Joins(fmt.Sprintf(`INNER JOIN %s svc_type ON svc_type.id = svc.service_type_id`, serviceTypeTable)).
			Where("svc_type.department_id = ?", *departmentID)
	}

	if len(query.BranchIDs) > 0 {
		bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
		roomTable := r.schemaTable(schema, entity.RoomTableName)
		floorTable := r.schemaTable(schema, entity.FloorTableName)

		db = db.
			Joins(fmt.Sprintf("JOIN %s br_scope ON br_scope.id = os.booking_room_id", bookingRoomTable)).
			Joins(fmt.Sprintf("JOIN %s r_scope ON r_scope.id = br_scope.room_id", roomTable)).
			Joins(fmt.Sprintf("JOIN %s f_scope ON f_scope.id = r_scope.floor_id", floorTable)).
			Where("f_scope.branch_id IN ?", query.BranchIDs)
	}

	if query.Status != "" {
		db = db.Where("os.status = ?", query.Status)
	}

	if !query.FromDate.IsZero() {
		db = db.Where("os.created_at >= ?", query.FromDate)
	}

	if !query.ToDate.IsZero() {
		db = db.Where("os.created_at < ?", query.ToDate.AddDate(0, 0, 1))
	}

	return db
}

func (r *orderServiceRepositoryImpl) applySorting(db *gorm.DB, query model.OrderServiceQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at":   "os.created_at",
		"estimated_at": "os.estimated_at",
	}

	sortField := "os.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}

func (r *orderServiceRepositoryImpl) FindDashboardServiceTypeCount(ctx context.Context, schema string, branchIDs []int64) (*model.DashboardServiceTypeCount, error) {
	serviceTypeTable := r.schemaTable(schema, entity.ServiceTypeTableName)
	serviceTable := r.schemaTable(schema, entity.ServiceTableName)
	orderServiceTable := r.schemaTable(schema, entity.OrderServiceTableName)
	bookingRoomTable := r.schemaTable(schema, entity.BookingRoomTableName)
	roomTable := r.schemaTable(schema, entity.RoomTableName)
	floorTable := r.schemaTable(schema, entity.FloorTableName)
	branchJoins := ""
	branchClause := ""
	args := []any{
		model.OrderServiceStatusPending,
		model.OrderServiceStatusConfirmed,
	}
	if len(branchIDs) > 0 {
		branchJoins = fmt.Sprintf(
			`JOIN %s AS br_dash ON br_dash.id = os.booking_room_id
			JOIN %s AS r_dash ON r_dash.id = br_dash.room_id
			JOIN %s AS f_dash ON f_dash.id = r_dash.floor_id`,
			bookingRoomTable,
			roomTable,
			floorTable,
		)
		branchClause = " AND f_dash.branch_id IN ?"
		args = append(args, branchIDs)
	}

	sql := fmt.Sprintf(
		`SELECT
			st.id AS service_type_id,
			st.name AS service_type_name,
			COALESCE(osc.order_count, 0) AS order_count,
			CASE
				WHEN SUM(COALESCE(osc.order_count, 0)) OVER () = 0 THEN 0
				ELSE COALESCE(osc.order_count, 0) * 100.0 / SUM(COALESCE(osc.order_count, 0)) OVER ()
			END AS order_percent
		FROM %s AS st
		LEFT JOIN (
			SELECT
				s.service_type_id,
				COUNT(*) AS order_count
			FROM %s AS os
			JOIN %s AS s ON s.id = os.service_id
			%s
			WHERE os.status IN (?, ?)%s
			GROUP BY s.service_type_id
		) AS osc ON osc.service_type_id = st.id
		ORDER BY order_count DESC, st.name ASC`,
		serviceTypeTable,
		orderServiceTable,
		serviceTable,
		branchJoins,
		branchClause,
	)

	var rows []*mapper.OrderServiceDashboardServiceTypeCountProjection
	if err := r.getDB(ctx).Raw(
		sql,
		args...,
	).Scan(&rows).Error; err != nil {
		return nil, err
	}

	return mapper.ToDashboardServiceTypeCountModel(rows), nil
}
