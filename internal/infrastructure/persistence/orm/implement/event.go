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

type eventRepositoryImpl struct {
	BaseRepository
}

func NewEventRepository(baseRepo BaseRepository) repository.EventRepository {
	return &eventRepositoryImpl{baseRepo}
}

func (r *eventRepositoryImpl) Create(ctx context.Context, schema string, event *model.Event) error {
	eventEntity := mapper.ToEventEntity(event)
	branchIDs := eventBranchIDs(event)

	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table(r.schemaTable(schema, entity.EventTableName)).
			Omit("Branches").
			Create(eventEntity).Error; err != nil {
			return err
		}

		return r.createBranchEvents(tx, schema, event.ID, branchIDs)
	})
}

func (r *eventRepositoryImpl) FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.Event, error) {
	var eventEntity entity.EventEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.EventTableName)).
		Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name")).
		Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name")).
		Where("id = ?", id).
		First(&eventEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	event := mapper.ToEvent(&eventEntity)
	if err := r.attachBranches(ctx, schema, []*model.Event{event}); err != nil {
		return nil, err
	}

	return event, nil
}

func (r *eventRepositoryImpl) FindBySlug(ctx context.Context, schema string, slug string) (*model.Event, error) {
	var eventEntity entity.EventEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.EventTableName)).
		Where("slug = ?", slug).
		First(&eventEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	event := mapper.ToEvent(&eventEntity)
	if err := r.attachBranches(ctx, schema, []*model.Event{event}); err != nil {
		return nil, err
	}

	return event, nil
}

func (r *eventRepositoryImpl) FindByIDs(ctx context.Context, schema string, ids []int64) ([]*model.Event, error) {
	var eventEntities []*entity.EventEntity

	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.EventTableName)).
		Where("id IN ?", ids).
		Find(&eventEntities).Error; err != nil {
		return nil, err
	}

	events := mapper.ToEvents(eventEntities)
	if err := r.attachBranches(ctx, schema, events); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *eventRepositoryImpl) FindAllPaginated(ctx context.Context, schema string, query model.EventQuery) ([]*model.Event, int64, error) {
	var eventEntities []*entity.EventEntity
	var total int64

	eventTable := r.schemaTable(schema, entity.EventTableName)

	db := r.getDB(ctx).
		Table(eventTable + " AS e").
		Model(&entity.EventEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Event{}, 0, nil
	}

	db = db.Session(&gorm.Session{})
	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("e.id", "e.slug", "e.title", "e.excerpt", "e.thumbnail", "e.is_active", "e.published_at", "e.created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&eventEntities).Error; err != nil {
		return nil, 0, err
	}

	events := mapper.ToEvents(eventEntities)
	if err := r.attachBranches(ctx, schema, events); err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs []int64) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.
			Table(r.schemaTable(schema, entity.EventTableName)).
			Model(&entity.EventEntity{}).
			Where("id = ?", id).
			Updates(updateData)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return exceptions.ErrEventNotFound
		}

		return r.replaceBranchEvents(tx, schema, id, branchIDs)
	})
}

func (r *eventRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.EventTableName)).
		Where("id = ?", id).
		Delete(&entity.EventEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrEventNotFound
	}

	return nil
}

func (r *eventRepositoryImpl) DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error) {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.EventTableName)).
		Where("id IN ?", ids).
		Delete(&entity.EventEntity{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *eventRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.EventQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("e.title ILIKE ?", term)
	}

	if query.IsActive != nil {
		db = db.Where("e.is_active = ?", *query.IsActive)
	}

	if query.BranchID != 0 {
		db = r.applyBranchFilter(db, schema, []int64{query.BranchID}, query.IncludeGlobal)
	} else if len(query.BranchIDs) > 0 {
		db = r.applyBranchFilter(db, schema, query.BranchIDs, query.IncludeGlobal)
	}

	return db
}

func (r *eventRepositoryImpl) applyBranchFilter(db *gorm.DB, schema string, branchIDs []int64, includeGlobal bool) *gorm.DB {
	branchEventTable := r.schemaTable(schema, entity.BranchEventTableName)
	globalCondition := fmt.Sprintf("NOT EXISTS (SELECT 1 FROM %s be_all WHERE be_all.event_id = e.id)", branchEventTable)
	branchCondition := fmt.Sprintf("EXISTS (SELECT 1 FROM %s be WHERE be.event_id = e.id AND be.branch_id IN ?)", branchEventTable)

	if includeGlobal {
		return db.Where(globalCondition+" OR "+branchCondition, branchIDs)
	}

	return db.Where(branchCondition, branchIDs)
}

func (r *eventRepositoryImpl) createBranchEvents(tx *gorm.DB, schema string, eventID int64, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return nil
	}

	relations := make([]*entity.BranchEventEntity, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		relations = append(relations, &entity.BranchEventEntity{
			BranchID: branchID,
			EventID:  eventID,
		})
	}

	return tx.
		Table(r.schemaTable(schema, entity.BranchEventTableName)).
		Create(&relations).Error
}

func (r *eventRepositoryImpl) replaceBranchEvents(tx *gorm.DB, schema string, eventID int64, branchIDs []int64) error {
	if err := tx.
		Table(r.schemaTable(schema, entity.BranchEventTableName)).
		Where("event_id = ?", eventID).
		Delete(&entity.BranchEventEntity{}).Error; err != nil {
		return err
	}

	return r.createBranchEvents(tx, schema, eventID, branchIDs)
}

func (r *eventRepositoryImpl) attachBranches(ctx context.Context, schema string, events []*model.Event) error {
	if len(events) == 0 {
		return nil
	}

	eventIDs := make([]int64, 0, len(events))
	eventsByID := make(map[int64]*model.Event, len(events))
	for _, event := range events {
		if event == nil {
			continue
		}
		eventIDs = append(eventIDs, event.ID)
		eventsByID[event.ID] = event
		event.Branches = make([]*model.Branch, 0)
	}

	type branchRow struct {
		EventID  int64  `gorm:"column:event_id"`
		BranchID int64  `gorm:"column:branch_id"`
		Name     string `gorm:"column:name"`
	}

	var rows []*branchRow
	joinTable := r.schemaTable(schema, entity.BranchEventTableName)
	branchTable := r.schemaTable(schema, entity.BranchTableName)
	if err := r.getDB(ctx).
		Table(joinTable+" AS be").
		Select("be.event_id, b.id AS branch_id, b.name").
		Joins(fmt.Sprintf("JOIN %s b ON b.id = be.branch_id", branchTable)).
		Where("be.event_id IN ?", eventIDs).
		Order("b.name ASC").
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		event, ok := eventsByID[row.EventID]
		if !ok {
			continue
		}
		event.Branches = append(event.Branches, &model.Branch{
			Base: model.Base{ID: row.BranchID},
			Name: row.Name,
		})
	}

	return nil
}

func eventBranchIDs(event *model.Event) []int64 {
	if event == nil || len(event.Branches) == 0 {
		return nil
	}

	branchIDs := make([]int64, 0, len(event.Branches))
	for _, branch := range event.Branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs
}

func (r *eventRepositoryImpl) applySorting(db *gorm.DB, query model.EventQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at":   "e.created_at",
		"title":        "e.title",
		"published_at": "e.published_at",
		"is_active":    "e.is_active",
	}

	sortField := "e.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}
