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

type userRepositoryImpl struct {
	BaseRepository
}

func NewUserRepository(baseRepo BaseRepository) repository.UserRepository {
	return &userRepositoryImpl{baseRepo}
}

func (r *userRepositoryImpl) Create(ctx context.Context, schema string, user *model.User) error {
	userEntity := mapper.ToUserEntity(user)
	branchIDs := userBranchIDs(user)

	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.
			Table(r.schemaTable(schema, entity.UserTableName)).
			Omit("Branches").
			Create(userEntity).Error; err != nil {
			return err
		}

		return r.createUserBranches(tx, schema, user.ID, branchIDs)
	})
}

func (r *userRepositoryImpl) FindByUsernameWithDepartmentDetails(ctx context.Context, schema, username string) (*model.User, error) {
	var userEntity entity.UserEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Preload("Department", r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name")).
		Where("username = ? AND is_active = true", username).
		First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	user := mapper.ToUser(&userEntity)

	if err := r.loadDepartmentPermissions(ctx, schema, user); err != nil {
		return nil, err
	}
	if err := r.attachBranches(ctx, schema, []*model.User{user}); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByIDWithDepartmentDetails(ctx context.Context, schema string, id int64) (*model.User, error) {
	user, err := r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Department",
			Scope:    r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name"),
		},
	)
	if err != nil || user == nil {
		return nil, err
	}

	if err := r.loadDepartmentPermissions(ctx, schema, user); err != nil {
		return nil, err
	}
	if err := r.attachBranches(ctx, schema, []*model.User{user}); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByIDWithDetails(ctx context.Context, schema string, id int64) (*model.User, error) {
	user, err := r.findByIDBase(r.getDB(ctx), schema, id,
		Preload{
			Relation: "Department",
			Scope:    r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name"),
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
	if err != nil || user == nil {
		return user, err
	}
	if err := r.attachBranches(ctx, schema, []*model.User{user}); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.User, error) {
	user, err := r.findByIDBase(r.getDB(ctx), schema, id)
	if err != nil || user == nil {
		return user, err
	}
	if err := r.attachBranches(ctx, schema, []*model.User{user}); err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) FindByEmail(ctx context.Context, schema, email string) (*model.User, error) {
	var userEntity entity.UserEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Where("email = ?", email).
		First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToUser(&userEntity), nil
}

func (r *userRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any, branchIDs ...[]int64) error {
	return r.getDB(ctx).Transaction(func(tx *gorm.DB) error {
		result := tx.
			Table(r.schemaTable(schema, entity.UserTableName)).
			Model(&entity.UserEntity{}).
			Where("id = ?", id).
			Updates(updateData)
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return exceptions.ErrUserNotFound
		}

		if len(branchIDs) == 0 {
			return nil
		}

		return r.replaceUserBranches(tx, schema, id, branchIDs[0])
	})
}

func (r *userRepositoryImpl) ExistsByEmail(ctx context.Context, schema, email string) (bool, error) {
	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Model(&entity.UserEntity{}).
		Where("email = ?", email).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) ExistsActiveAdminExceptID(ctx context.Context, schema string, id int64) (bool, error) {
	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Model(&entity.UserEntity{}).
		Where("role = ? AND is_active = true AND id <> ?", model.RoleAdmin, id).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) ExistsActiveAdmin(ctx context.Context, schema string) (bool, error) {
	var count int64
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Model(&entity.UserEntity{}).
		Where("role = ? AND is_active = true", model.RoleAdmin).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *userRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Where("id = ?", id).
		Delete(&entity.UserEntity{})
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return exceptions.ErrUserNotFound
	}

	return nil
}

func (r *userRepositoryImpl) DeleteByIDs(ctx context.Context, schema string, ids []int64) (int64, error) {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Where("id IN ?", ids).
		Delete(&entity.UserEntity{})
	if result.Error != nil {
		return 0, result.Error
	}

	return result.RowsAffected, nil
}

func (r *userRepositoryImpl) FindAllActiveByDepartmentID(ctx context.Context, schema string, deptID int64) ([]*model.User, error) {
	var userEntities []*entity.UserEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Where("department_id = ? AND is_active = true", deptID).
		Find(&userEntities).Error; err != nil {
		return nil, err
	}

	users := mapper.ToUsers(userEntities)
	if err := r.attachBranches(ctx, schema, users); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *userRepositoryImpl) FindByIDs(ctx context.Context, schema string, ids []int64) ([]*model.User, error) {
	var userEntities []*entity.UserEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Where("id IN ?", ids).
		Find(&userEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToUsers(userEntities), nil
}

func (r *userRepositoryImpl) FindAllWithDepartmentPaginated(ctx context.Context, schema string, query model.UserQuery) ([]*model.User, int64, error) {
	var userEntities []*entity.UserEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName) + " AS u").
		Model(&entity.UserEntity{})

	db = r.applyFilters(db, schema, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.User{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("Department", r.preloadWithSchema(schema, entity.DepartmentTableName, "id", "name"))

	db = r.applySorting(db, query)

	offset := (query.Page - 1) * query.Limit

	if err := db.Select("u.id", "u.role", "u.first_name", "u.last_name", "u.is_active", "u.department_id", "u.created_at").
		Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&userEntities).Error; err != nil {
		return nil, 0, err
	}

	users := mapper.ToUsers(userEntities)
	if err := r.attachBranches(ctx, schema, users); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *userRepositoryImpl) findByIDBase(tx *gorm.DB, schema string, id int64, preloads ...Preload) (*model.User, error) {
	var userEntity entity.UserEntity

	tx = tx.Table(r.schemaTable(schema, entity.UserTableName))

	for _, preload := range preloads {
		if preload.Scope != nil {
			tx = tx.Preload(preload.Relation, preload.Scope)
		} else {
			tx = tx.Preload(preload.Relation)
		}
	}

	if err := tx.Where("id = ?", id).
		First(&userEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToUser(&userEntity), nil
}

func (r *userRepositoryImpl) applyFilters(db *gorm.DB, schema string, query model.UserQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where(
			"u.username ILIKE ? OR u.first_name ILIKE ? OR u.last_name ILIKE ?",
			term, term, term,
		)
	}

	if query.Role != "" {
		db = db.Where("u.role = ?", query.Role)
	}

	if query.IsActive != nil {
		db = db.Where("u.is_active = ?", *query.IsActive)
	}

	if query.DepartmentID != 0 {
		db = db.Where("u.department_id = ?", query.DepartmentID)
	}

	if query.BranchID != 0 {
		userBranchTable := r.schemaTable(schema, entity.UserBranchTableName)
		db = db.Joins(fmt.Sprintf("JOIN %s ub ON ub.user_id = u.id", userBranchTable)).
			Where("ub.branch_id = ?", query.BranchID)
	}

	return db
}

func (r *userRepositoryImpl) applySorting(db *gorm.DB, query model.UserQuery) *gorm.DB {
	allowedSorts := map[string]string{
		"created_at": "u.created_at",
		"first_name": "u.first_name",
		"last_name":  "u.last_name",
	}

	sortField := "u.created_at"
	if field, ok := allowedSorts[query.Sort]; ok {
		sortField = field
	}

	order := "DESC"
	if strings.ToUpper(query.Order) == "ASC" {
		order = "ASC"
	}

	return db.Order(sortField + " " + order)
}

func (r *userRepositoryImpl) loadDepartmentPermissions(ctx context.Context, schema string, user *model.User) error {
	if user == nil || user.Department == nil {
		return nil
	}

	joinTable := r.schemaTable(schema, entity.DepartmentPermissionTableName)
	permTable := r.schemaTable(schema, entity.PermissionTableName)

	var permEntities []*entity.PermissionEntity
	if err := r.getDB(ctx).
		Table(permTable+" AS p").
		Joins(fmt.Sprintf(`JOIN %s dp ON dp.permission_id = p.id`, joinTable)).
		Where("dp.department_id = ?", user.DepartmentID).
		Find(&permEntities).Error; err != nil {
		return err
	}
	user.Department.Permissions = mapper.ToPermissions(permEntities)

	return nil
}

func (r *userRepositoryImpl) createUserBranches(tx *gorm.DB, schema string, userID int64, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return nil
	}

	relations := make([]*entity.UserBranchEntity, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		relations = append(relations, &entity.UserBranchEntity{
			UserID:   userID,
			BranchID: branchID,
		})
	}

	return tx.
		Table(r.schemaTable(schema, entity.UserBranchTableName)).
		Create(&relations).Error
}

func (r *userRepositoryImpl) replaceUserBranches(tx *gorm.DB, schema string, userID int64, branchIDs []int64) error {
	if err := tx.
		Table(r.schemaTable(schema, entity.UserBranchTableName)).
		Where("user_id = ?", userID).
		Delete(&entity.UserBranchEntity{}).Error; err != nil {
		return err
	}

	return r.createUserBranches(tx, schema, userID, branchIDs)
}

func (r *userRepositoryImpl) attachBranches(ctx context.Context, schema string, users []*model.User) error {
	if len(users) == 0 {
		return nil
	}

	userIDs := make([]int64, 0, len(users))
	usersByID := make(map[int64]*model.User, len(users))
	for _, user := range users {
		if user == nil {
			continue
		}
		userIDs = append(userIDs, user.ID)
		usersByID[user.ID] = user
		user.Branches = make([]*model.Branch, 0)
	}
	if len(userIDs) == 0 {
		return nil
	}

	type branchRow struct {
		UserID   int64  `gorm:"column:user_id"`
		BranchID int64  `gorm:"column:branch_id"`
		Name     string `gorm:"column:name"`
	}

	var rows []*branchRow
	joinTable := r.schemaTable(schema, entity.UserBranchTableName)
	branchTable := r.schemaTable(schema, entity.BranchTableName)
	if err := r.getDB(ctx).
		Table(joinTable+" AS ub").
		Select("ub.user_id, b.id AS branch_id, b.name").
		Joins(fmt.Sprintf("JOIN %s b ON b.id = ub.branch_id", branchTable)).
		Where("ub.user_id IN ?", userIDs).
		Order("b.name ASC").
		Scan(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		user, ok := usersByID[row.UserID]
		if !ok {
			continue
		}
		user.Branches = append(user.Branches, &model.Branch{
			Base: model.Base{ID: row.BranchID},
			Name: row.Name,
		})
	}

	return nil
}

func userBranchIDs(user *model.User) []int64 {
	if user == nil || len(user.Branches) == 0 {
		return nil
	}

	branchIDs := make([]int64, 0, len(user.Branches))
	for _, branch := range user.Branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs
}
