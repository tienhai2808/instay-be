package implement

import (
	"context"
	"errors"
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/mapper"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/gorm"
)

type RoomsCountResult struct {
	RoomTypeID int64 `gorm:"column:room_type_id"`
	RoomsCount int64 `gorm:"column:rooms_count"`
}

type roomTypeRepositoryImpl struct {
	BaseRepository
}

func NewRoomTypeRepository(baseRepo BaseRepository) repository.RoomTypeRepository {
	return &roomTypeRepositoryImpl{baseRepo}
}

func (r *roomTypeRepositoryImpl) Create(ctx context.Context, schema string, roomType *model.RoomType) error {
	roomTypeEntity := mapper.ToRoomTypeEntity(roomType)
	branchIDs := roomTypeBranchIDs(roomType)

	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table(r.schemaTable(schema, entity.RoomTypeTableName)).
			Omit("Branches").
			Create(roomTypeEntity).Error; err != nil {
			return err
		}

		return r.createBranchRoomTypes(tx, schema, roomType.ID, branchIDs)
	})
}

func (r *roomTypeRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.RoomType, error) {
	var roomTypeEntity entity.RoomTypeEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypeTableName)).
		Where("id = ?", id).
		First(&roomTypeEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	roomType := mapper.ToRoomType(&roomTypeEntity)
	if err := r.attachBranches(ctx, schema, []*model.RoomType{roomType}); err != nil {
		return nil, err
	}

	return roomType, nil
}

func (r *roomTypeRepositoryImpl) ExistsByIDAndBranchID(ctx context.Context, schema string, id, branchID int64) (bool, error) {
	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.BranchRoomTypeTableName)).
		Where("room_type_id = ? AND branch_id = ?", id, branchID).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *roomTypeRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.RoomTypeQuery) ([]*model.RoomType, int64, error) {
	var roomTypeEntities []*entity.RoomTypeEntity
	var total int64
	roomTypeTable := r.schemaTable(schema, entity.RoomTypeTableName)

	db := r.getDB(ctx).
		Table(roomTypeTable + " AS rt").
		Model(&entity.RoomTypeEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.RoomType{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("RoomTypePrices", r.preloadWithSchema(schema, entity.RoomTypePriceTableName))
	db = db.Preload("RoomTypePrices.PricingTag", r.preloadWithSchema(schema, entity.PricingTagTableName, "id", "name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&roomTypeEntities).Error; err != nil {
		return nil, 0, err
	}

	roomTypes := mapper.ToRoomTypes(roomTypeEntities)
	if err := r.attachBranches(ctx, schema, roomTypes); err != nil {
		return nil, 0, err
	}

	return roomTypes, total, nil
}

func (r *roomTypeRepositoryImpl) SearchByName(ctx context.Context, schema string, keyword string) (*model.RoomType, error) {
	var roomTypeEntity entity.RoomTypeEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypeTableName)).
		Where("name ILIKE ?", "%"+keyword+"%").
		First(&roomTypeEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, exceptions.ErrRoomTypeNotFound
		}
	}

	return mapper.ToRoomType(&roomTypeEntity), nil
}

func (r *roomTypeRepositoryImpl) FindAll(ctx context.Context, schema string, query model.RoomTypeQuery) ([]*model.RoomType, error) {
	var roomTypeEntities []*entity.RoomTypeEntity
	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypeTableName)+" AS rt").
		Select("id", "name").
		Model(&entity.RoomTypeEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Find(&roomTypeEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToRoomTypes(roomTypeEntities), nil
}

func (r *roomTypeRepositoryImpl) CountRoomByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error) {
	var counts []RoomsCountResult
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTableName)).
		Select("room_type_id, COUNT(*) as rooms_count").
		Where("room_type_id IN ?", ids).
		Group("room_type_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.RoomTypeID] = c.RoomsCount
	}

	return countMap, nil
}

func (r *roomTypeRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs []int64) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.
			Table(r.schemaTable(schema, entity.RoomTypeTableName)).
			Model(&entity.RoomTypeEntity{}).
			Where("id = ?", id).
			Updates(updateData)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return exceptions.ErrRoomTypeNotFound
		}

		return r.replaceBranchRoomTypes(tx, schema, id, branchIDs)
	})
}

func (r *roomTypeRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.RoomTypeTableName)).
		Where("id = ?", id).
		Delete(&entity.RoomTypeEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrRoomTypeNotFound
	}

	return nil
}

func (r *roomTypeRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.RoomTypeQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("rt.name ILIKE ?", term)
	}

	if query.BranchID != 0 {
		joinTable := r.schemaTable(schema, entity.BranchRoomTypeTableName)
		db = db.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s brt WHERE brt.room_type_id = rt.id AND brt.branch_id = ?)", joinTable), query.BranchID)
	} else if len(query.BranchIDs) > 0 {
		joinTable := r.schemaTable(schema, entity.BranchRoomTypeTableName)
		db = db.Where(fmt.Sprintf("EXISTS (SELECT 1 FROM %s brt WHERE brt.room_type_id = rt.id AND brt.branch_id IN ?)", joinTable), query.BranchIDs)
	}

	return db
}

func (r *roomTypeRepositoryImpl) createBranchRoomTypes(tx *gorm.DB, schema string, roomTypeID int64, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return nil
	}

	relations := make([]*entity.BranchRoomTypeEntity, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		relations = append(relations, &entity.BranchRoomTypeEntity{
			BranchID:   branchID,
			RoomTypeID: roomTypeID,
		})
	}

	return tx.
		Table(r.schemaTable(schema, entity.BranchRoomTypeTableName)).
		Create(&relations).Error
}

func (r *roomTypeRepositoryImpl) replaceBranchRoomTypes(tx *gorm.DB, schema string, roomTypeID int64, branchIDs []int64) error {
	if err := tx.
		Table(r.schemaTable(schema, entity.BranchRoomTypeTableName)).
		Where("room_type_id = ?", roomTypeID).
		Delete(&entity.BranchRoomTypeEntity{}).Error; err != nil {
		return err
	}

	return r.createBranchRoomTypes(tx, schema, roomTypeID, branchIDs)
}

func (r *roomTypeRepositoryImpl) attachBranches(ctx context.Context, schema string, roomTypes []*model.RoomType) error {
	if len(roomTypes) == 0 {
		return nil
	}

	roomTypeIDs := make([]int64, 0, len(roomTypes))
	roomTypesByID := make(map[int64]*model.RoomType, len(roomTypes))
	for _, roomType := range roomTypes {
		if roomType == nil {
			continue
		}
		roomTypeIDs = append(roomTypeIDs, roomType.ID)
		roomTypesByID[roomType.ID] = roomType
		roomType.Branches = make([]*model.Branch, 0)
	}

	type branchRow struct {
		RoomTypeID int64  `gorm:"column:room_type_id"`
		BranchID   int64  `gorm:"column:branch_id"`
		Name       string `gorm:"column:name"`
	}

	var rows []*branchRow
	joinTable := r.schemaTable(schema, entity.BranchRoomTypeTableName)
	branchTable := r.schemaTable(schema, entity.BranchTableName)
	if err := r.getDB(ctx).
		Table(joinTable+" AS brt").
		Select("brt.room_type_id, b.id AS branch_id, b.name").
		Joins(fmt.Sprintf("JOIN %s b ON b.id = brt.branch_id", branchTable)).
		Where("brt.room_type_id IN ?", roomTypeIDs).
		Order("b.name ASC").
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		roomType, ok := roomTypesByID[row.RoomTypeID]
		if !ok {
			continue
		}
		roomType.Branches = append(roomType.Branches, &model.Branch{
			Base: model.Base{ID: row.BranchID},
			Name: row.Name,
		})
	}

	return nil
}

func roomTypeBranchIDs(roomType *model.RoomType) []int64 {
	if roomType == nil || len(roomType.Branches) == 0 {
		return nil
	}

	branchIDs := make([]int64, 0, len(roomType.Branches))
	for _, branch := range roomType.Branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs
}
