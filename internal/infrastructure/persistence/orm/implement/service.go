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

type serviceRepositoryImpl struct {
	BaseRepository
}

func NewServiceRepository(baseRepo BaseRepository) repository.ServiceRepository {
	return &serviceRepositoryImpl{baseRepo}
}

func (r *serviceRepositoryImpl) Create(ctx context.Context, schema string, service *model.Service) error {
	svcEntity := mapper.ToServiceEntity(service)
	branchIDs := serviceBranchIDs(service)

	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table(r.schemaTable(schema, entity.ServiceTableName)).
			Omit("Branches").
			Create(svcEntity).Error; err != nil {
			return err
		}

		return r.createBranchServices(tx, schema, service.ID, branchIDs)
	})
}

func (r *serviceRepositoryImpl) FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.Service, error) {
	service, err := r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "ServiceType",
			Scope:    r.preloadWithSchema(schema, entity.ServiceTypeTableName, "id", "name", "icon"),
		},
		Preload{
			Relation: "ServiceImages",
			Scope:    r.preloadWithSchema(schema, entity.ServiceImageTableName),
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
	if err != nil || service == nil {
		return service, err
	}
	if err := r.attachBranches(ctx, schema, []*model.Service{service}); err != nil {
		return nil, err
	}

	return service, nil
}

func (r *serviceRepositoryImpl) FindBySlugWithServiceTypeAndServiceImages(ctx context.Context, schema string, slug string) (*model.Service, error) {
	var svcEntity entity.ServiceEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTableName)).
		Preload("ServiceType", r.preloadWithSchema(schema, entity.ServiceTypeTableName, "id", "name", "icon")).
		Preload("ServiceImages", r.preloadWithSchema(schema, entity.ServiceImageTableName)).
		Where("slug = ?", slug).
		First(&svcEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	service := mapper.ToService(&svcEntity)
	if err := r.attachBranches(ctx, schema, []*model.Service{service}); err != nil {
		return nil, err
	}

	return service, nil
}

func (r *serviceRepositoryImpl) FindByIDWithServiceImages(ctx context.Context, schema string, id int64) (*model.Service, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "ServiceImages",
			Scope:    r.preloadWithSchema(schema, entity.ServiceImageTableName),
		},
	)
}

func (r *serviceRepositoryImpl) FindByIDWithServiceType(ctx context.Context, schema string, id int64) (*model.Service, error) {
	return r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "ServiceType",
			Scope:    r.preloadWithSchema(schema, entity.ServiceTypeTableName),
		},
	)
}

func (r *serviceRepositoryImpl) ExistsByIDAndBranchID(ctx context.Context, schema string, id, branchID int64) (bool, error) {
	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchServiceTableName)).
		Where("service_id = ? AND branch_id = ?", id, branchID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *serviceRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs []int64) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.
			Table(r.schemaTable(schema, entity.ServiceTableName)).
			Model(&entity.ServiceEntity{}).
			Where("id = ?", id).
			Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return exceptions.ErrServiceNotFound
		}

		return r.replaceBranchServices(tx, schema, id, branchIDs)
	})
}

func (r *serviceRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTableName)).
		Where("id = ?", id).
		Delete(&entity.ServiceEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrServiceNotFound
	}

	return nil
}

func (r *serviceRepositoryImpl) DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error) {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTableName)).
		Where("id IN ?", ids).
		Delete(&entity.ServiceEntity{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *serviceRepositoryImpl) FindByIDsWithServiceImages(ctx context.Context, schema string, ids []int64) ([]*model.Service, error) {
	var svcEntities []*entity.ServiceEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.ServiceTableName)).
		Preload("ServiceImages", r.preloadWithSchema(schema, entity.ServiceImageTableName)).
		Where("id IN ?", ids).
		Find(&svcEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToServices(svcEntities), nil
}

func (r *serviceRepositoryImpl) FindAllWithServiceTypePaginated(ctx context.Context, schema string, query model.ServiceQuery) ([]*model.Service, int64, error) {
	var serviceEntities []*entity.ServiceEntity
	var total int64

	serviceTable := r.schemaTable(schema, entity.ServiceTableName)

	db := r.getDB(ctx).
		Table(serviceTable + " AS svc").
		Model(&entity.ServiceEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Service{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("ServiceType", r.preloadWithSchema(schema, entity.ServiceTypeTableName, "id", "name", "icon"))

	db = r.applySorting(db, schema, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("svc.id", "svc.name", "svc.price", "svc.is_active", "svc.service_type_id", "svc.created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&serviceEntities).Error; err != nil {
		return nil, 0, err
	}

	services := mapper.ToServices(serviceEntities)
	if err := r.attachBranches(ctx, schema, services); err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *serviceRepositoryImpl) FindAllWithServiceTypeAndThumbnailPaginated(ctx context.Context, schema string, query model.ServiceQuery) ([]*model.Service, int64, error) {
	var serviceEntities []*entity.ServiceEntity
	var total int64

	serviceTable := r.schemaTable(schema, entity.ServiceTableName)

	db := r.getDB(ctx).
		Table(serviceTable + " AS svc").
		Model(&entity.ServiceEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Service{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("ServiceImages", func(db *gorm.DB) *gorm.DB {
		return db.
			Table(r.schemaTable(schema, entity.ServiceImageTableName)).
			Where("is_thumbnail = true")
	})

	db = r.applySorting(db, schema, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("svc.id", "svc.name", "svc.slug", "svc.price", "svc.discount_type", "svc.discount_value", "svc.discount_start_at", "svc.discount_end_at", "svc.created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&serviceEntities).Error; err != nil {
		return nil, 0, err
	}

	services := mapper.ToServices(serviceEntities)
	if err := r.attachBranches(ctx, schema, services); err != nil {
		return nil, 0, err
	}

	return services, total, nil
}

func (r *serviceRepositoryImpl) FindAllActiveByServiceTypeSlugWithThumbnail(ctx context.Context, schema string, serviceTypeSlug string, branchID int64) ([]*model.Service, error) {
	var serviceEntities []*entity.ServiceEntity

	serviceTable := r.schemaTable(schema, entity.ServiceTableName)
	serviceTypeTable := r.schemaTable(schema, entity.ServiceTypeTableName)
	serviceImageTable := r.schemaTable(schema, entity.ServiceImageTableName)

	db := r.getDB(ctx).
		Table(serviceTable+" AS svc").
		Joins("INNER JOIN "+serviceTypeTable+" st ON st.id = svc.service_type_id").
		Where("st.slug = ?", serviceTypeSlug).
		Where("svc.is_active = true")

	if branchID != 0 {
		branchServiceTable := r.schemaTable(schema, entity.BranchServiceTableName)
		db = db.Joins(fmt.Sprintf("JOIN %s bs ON bs.service_id = svc.id", branchServiceTable)).
			Where("bs.branch_id = ?", branchID)
	}

	if err := db.
		Preload("ServiceImages", func(db *gorm.DB) *gorm.DB {
			return db.
				Table(serviceImageTable).
				Where("is_thumbnail = true")
		}).
		Find(&serviceEntities).Error; err != nil {
		return nil, err
	}

	services := mapper.ToServices(serviceEntities)
	if err := r.attachBranches(ctx, schema, services); err != nil {
		return nil, err
	}

	return services, nil
}

func (r *serviceRepositoryImpl) findByIDBase(tx *gorm.DB, schema string, id int64, preloads ...Preload) (*model.Service, error) {
	var svcEntity entity.ServiceEntity

	tx = tx.Table(r.schemaTable(schema, entity.ServiceTableName))

	for _, preload := range preloads {
		if preload.Scope != nil {
			tx = tx.Preload(preload.Relation, preload.Scope)
		} else {
			tx = tx.Preload(preload.Relation)
		}
	}

	if err := tx.Where("id = ?", id).
		First(&svcEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToService(&svcEntity), nil
}

func (r *serviceRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.ServiceQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("svc.name ILIKE ?", term)
	}

	if query.IsActive != nil {
		db = db.Where("svc.is_active = ?", *query.IsActive)
	}

	if query.HasDiscount != nil {
		discountCondition := "svc.discount_type IS NOT NULL AND svc.discount_value IS NOT NULL AND svc.discount_value > 0 AND (svc.discount_start_at IS NULL OR svc.discount_start_at <= CURRENT_DATE) AND (svc.discount_end_at IS NULL OR svc.discount_end_at >= CURRENT_DATE)"
		if *query.HasDiscount {
			db = db.Where(discountCondition)
		} else {
			db = db.Where("NOT (" + discountCondition + ")")
		}
	}

	if query.ServiceTypeID != 0 {
		db = db.Where("svc.service_type_id = ?", query.ServiceTypeID)
	}

	if query.BranchID != 0 {
		branchServiceTable := r.schemaTable(schema, entity.BranchServiceTableName)
		db = db.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s bs WHERE bs.service_id = svc.id AND bs.branch_id = ?)", branchServiceTable), query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		branchServiceTable := r.schemaTable(schema, entity.BranchServiceTableName)
		db = db.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s bs WHERE bs.service_id = svc.id AND bs.branch_id IN ?)", branchServiceTable), query.BranchIDs)
	}

	return db
}

func (r *serviceRepositoryImpl) createBranchServices(tx *gorm.DB, schema string, serviceID int64, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return nil
	}

	relations := make([]*entity.BranchServiceEntity, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		relations = append(relations, &entity.BranchServiceEntity{
			BranchID:  branchID,
			ServiceID: serviceID,
		})
	}

	return tx.
		Table(r.schemaTable(schema, entity.BranchServiceTableName)).
		Create(&relations).Error
}

func (r *serviceRepositoryImpl) replaceBranchServices(tx *gorm.DB, schema string, serviceID int64, branchIDs []int64) error {
	if err := tx.
		Table(r.schemaTable(schema, entity.BranchServiceTableName)).
		Where("service_id = ?", serviceID).
		Delete(&entity.BranchServiceEntity{}).Error; err != nil {
		return err
	}

	return r.createBranchServices(tx, schema, serviceID, branchIDs)
}

func (r *serviceRepositoryImpl) attachBranches(ctx context.Context, schema string, services []*model.Service) error {
	if len(services) == 0 {
		return nil
	}

	serviceIDs := make([]int64, 0, len(services))
	servicesByID := make(map[int64]*model.Service, len(services))
	for _, service := range services {
		if service == nil {
			continue
		}
		serviceIDs = append(serviceIDs, service.ID)
		servicesByID[service.ID] = service
		service.Branches = make([]*model.Branch, 0)
	}

	type branchRow struct {
		ServiceID int64  `gorm:"column:service_id"`
		BranchID  int64  `gorm:"column:branch_id"`
		Name      string `gorm:"column:name"`
	}

	var rows []*branchRow
	joinTable := r.schemaTable(schema, entity.BranchServiceTableName)
	branchTable := r.schemaTable(schema, entity.BranchTableName)
	if err := r.getDB(ctx).
		Table(joinTable+" AS bs").
		Select("bs.service_id, b.id AS branch_id, b.name").
		Joins(fmt.Sprintf("JOIN %s b ON b.id = bs.branch_id", branchTable)).
		Where("bs.service_id IN ?", serviceIDs).
		Order("b.name ASC").
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		service, ok := servicesByID[row.ServiceID]
		if !ok {
			continue
		}
		service.Branches = append(service.Branches, &model.Branch{
			Base: model.Base{ID: row.BranchID},
			Name: row.Name,
		})
	}

	return nil
}

func serviceBranchIDs(service *model.Service) []int64 {
	if service == nil || len(service.Branches) == 0 {
		return nil
	}

	branchIDs := make([]int64, 0, len(service.Branches))
	for _, branch := range service.Branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs
}

func (r *serviceRepositoryImpl) applySorting(db *gorm.DB, schema string, query model.ServiceQuery) *gorm.DB {
	orderServiceTable := r.schemaTable(schema, entity.OrderServiceTableName)

	orderedCountExpr := "(SELECT COUNT(*) FROM " + orderServiceTable + " AS os WHERE os.service_id = svc.id)"
	orderedQuantityExpr := "(SELECT COALESCE(SUM(os.quantity), 0) FROM " + orderServiceTable + " AS os WHERE os.service_id = svc.id)"

	allowedSorts := map[string]string{
		"created_at":       "svc.created_at",
		"name":             "svc.name",
		"price":            "svc.price",
		"ordered_count":    orderedCountExpr,
		"ordered_quantity": orderedQuantityExpr,
	}

	sortField := "svc.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}
