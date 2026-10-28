package implement

import (
	"context"
	"fmt"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type customerRepositoryImpl struct {
	BaseRepository
}

func NewCustomerRepository(baseRepo BaseRepository) repository.CustomerRepository {
	return &customerRepositoryImpl{baseRepo}
}

func (r *customerRepositoryImpl) Create(ctx context.Context, schema string, customer *model.Customer) error {
	customerEntity := mapper.ToCustomerEntity(customer)
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.CustomerTableName)).
		Create(customerEntity).Error
}

func (r *customerRepositoryImpl) FindAllPaginated(ctx context.Context, schema string, query model.CustomerQuery) ([]*model.Customer, int64, error) {
	var customerEntities []*entity.CustomerEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.CustomerTableName) + " AS c").
		Model(&entity.CustomerEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Customer{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&customerEntities).Error; err != nil {
		return nil, 0, err
	}

	return mapper.ToCustomers(customerEntities), total, nil
}

func (r *customerRepositoryImpl) FindAll(ctx context.Context, schema string, query model.CustomerQuery) ([]*model.Customer, error) {
	var customerEntities []*entity.CustomerEntity

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.CustomerTableName) + " AS c").
		Model(&entity.CustomerEntity{})

	db = r.applyFilters(db, schema, query)
	db = r.applySorting(db, query)

	if err := db.Find(&customerEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToCustomers(customerEntities), nil
}

func (r *customerRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.CustomerQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("c.phone ILIKE ?", term)
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

func (r *customerRepositoryImpl) applySorting(db *gorm.DB, query model.CustomerQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at": "c.created_at",
		"birth_year": "c.birth_year",
		"gender":     "c.gender",
		"phone":      "c.phone",
	}

	sortField := "c.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}
