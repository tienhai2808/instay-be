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

type UsersCountResult struct {
	DepartmentID int64 `gorm:"column:department_id"`
	UsersCount   int64 `gorm:"column:users_count"`
}

type departmentRepositoryImpl struct {
	BaseRepository
}

func NewDepartmentRepository(baseRepo BaseRepository) repository.DepartmentRepository {
	return &departmentRepositoryImpl{baseRepo}
}

func (r *departmentRepositoryImpl) Create(ctx context.Context, schema string, dept *model.Department) error {
	deptEntity := mapper.ToDepartmentEntity(dept)
	deptEntity.Permissions = nil
	return r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Create(deptEntity).Error
}

func (r *departmentRepositoryImpl) FindAll(ctx context.Context, schema string) ([]*model.Department, error) {
	var departmentEntities []*entity.DepartmentEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Select("id", "name").
		Find(&departmentEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToDepartments(departmentEntities), nil
}

func (r *departmentRepositoryImpl) FindByName(ctx context.Context, schema, name string) (*model.Department, error) {
	var departmentEntity entity.DepartmentEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Where("name = ?", name).
		First(&departmentEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToDepartment(&departmentEntity), nil
}

func (r *departmentRepositoryImpl) FindByID(ctx context.Context, schema string, id int64) (*model.Department, error) {
	var departmentEntity entity.DepartmentEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Where("id = ?", id).
		First(&departmentEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return mapper.ToDepartment(&departmentEntity), nil
}

func (r *departmentRepositoryImpl) FindAllWithDetailsPaginated(ctx context.Context, schema string, query model.DepartmentQuery) ([]*model.Department, int64, error) {
	var departmentEntities []*entity.DepartmentEntity
	var total int64

	db := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Model(&entity.DepartmentEntity{})

	db = r.applyFilters(db, query)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Department{}, 0, nil
	}

	db = db.Session(&gorm.Session{})

	db = db.Preload("CreatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))
	db = db.Preload("UpdatedBy", r.preloadWithSchema(schema, entity.UserTableName, "id", "username", "first_name", "last_name"))

	offset := (query.Page - 1) * query.Limit

	if err := db.Offset(int(offset)).
		Limit(int(query.Limit)).
		Find(&departmentEntities).Error; err != nil {
		return nil, 0, err
	}

	deptIDs := make([]int64, len(departmentEntities))
	for i, d := range departmentEntities {
		deptIDs[i] = d.ID
	}

	type deptPerm struct {
		DepartmentID int64
		entity.PermissionEntity
	}
	var deptPerms []deptPerm

	joinTable := r.schemaTable(schema, entity.DepartmentPermissionTableName)
	permTable := r.schemaTable(schema, entity.PermissionTableName)

	if err := r.getDB(ctx).
		Table(permTable+" AS p").
		Select("dp.department_id, p.*").
		Joins(fmt.Sprintf(`JOIN %s dp ON dp.permission_id = p.id`, joinTable)).
		Where("dp.department_id IN ?", deptIDs).
		Find(&deptPerms).Error; err != nil {
		return nil, 0, err
	}

	permMap := make(map[int64][]*entity.PermissionEntity)
	for i := range deptPerms {
		dp := &deptPerms[i]
		permMap[dp.DepartmentID] = append(permMap[dp.DepartmentID], &dp.PermissionEntity)
	}
	for _, d := range departmentEntities {
		d.Permissions = permMap[d.ID]
	}

	return mapper.ToDepartments(departmentEntities), total, nil
}

func (r *departmentRepositoryImpl) CountUserByIDs(ctx context.Context, schema string, ids []int64) (map[int64]int64, error) {
	var counts []UsersCountResult
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.UserTableName)).
		Select("department_id, COUNT(*) as users_count").
		Where("department_id IN ?", ids).
		Group("department_id").
		Scan(&counts).Error; err != nil {
		return nil, err
	}

	countMap := make(map[int64]int64, len(counts))
	for _, c := range counts {
		countMap[c.DepartmentID] = c.UsersCount
	}

	return countMap, nil
}

func (r *departmentRepositoryImpl) FindByIDWithPermissions(ctx context.Context, schema string, id int64) (*model.Department, error) {
	var deptEntity entity.DepartmentEntity
	if err := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Where("id = ?", id).
		First(&deptEntity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	joinTable := r.schemaTable(schema, entity.DepartmentPermissionTableName)
	permTable := r.schemaTable(schema, entity.PermissionTableName)

	var permissionEntities []*entity.PermissionEntity
	if err := r.getDB(ctx).
		Table(permTable+" AS p").
		Joins(fmt.Sprintf(`JOIN %s dp ON dp.permission_id = p.id`, joinTable)).
		Where("dp.department_id = ?", id).
		Find(&permissionEntities).Error; err != nil {
		return nil, err
	}
	deptEntity.Permissions = permissionEntities

	return mapper.ToDepartment(&deptEntity), nil
}

func (r *departmentRepositoryImpl) UpdatePermissions(ctx context.Context, schema string, dept *model.Department, perms []*model.Permission) error {
	joinTable := r.schemaTable(schema, entity.DepartmentPermissionTableName)
	db := r.getDB(ctx)

	if err := db.Exec(
		fmt.Sprintf(`DELETE FROM %s WHERE department_id = ?`, joinTable),
		dept.ID,
	).Error; err != nil {
		return err
	}
	if len(perms) == 0 {
		return nil
	}

	rows := make([]entity.DepartmentPermissionEntity, 0, len(perms))
	for _, p := range perms {
		rows = append(rows, entity.DepartmentPermissionEntity{
			DepartmentID: dept.ID,
			PermissionID: p.ID,
		})
	}

	placeholders := make([]string, 0, len(rows))
	args := make([]any, 0, len(rows)*2)
	for i, row := range rows {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		args = append(args, row.DepartmentID, row.PermissionID)
	}

	query := fmt.Sprintf(
		`INSERT INTO %s (department_id, permission_id) VALUES %s ON CONFLICT DO NOTHING`,
		joinTable,
		strings.Join(placeholders, ", "),
	)

	return db.Exec(query, args...).Error
}

func (r *departmentRepositoryImpl) Update(ctx context.Context, schema string, id int64, updateData map[string]any) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Model(&entity.DepartmentEntity{}).
		Where("id = ?", id).
		Updates(updateData)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrDepartmentNotFound
	}

	return nil
}

func (r *departmentRepositoryImpl) Delete(ctx context.Context, schema string, id int64) error {
	result := r.getDB(ctx).
		Table(r.schemaTable(schema, entity.DepartmentTableName)).
		Where("id = ?", id).
		Delete(&entity.DepartmentEntity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return exceptions.ErrDepartmentNotFound
	}

	return nil
}

func (r *departmentRepositoryImpl) FindAllByPermissionNameWithActiveUsers(ctx context.Context, schema string, permName model.PermissionName) ([]*model.Department, error) {
	deptTable := r.schemaTable(schema, entity.DepartmentTableName)
	joinTable := r.schemaTable(schema, entity.DepartmentPermissionTableName)
	permTable := r.schemaTable(schema, entity.PermissionTableName)
	userTable := r.schemaTable(schema, entity.UserTableName)

	var departmentEntities []*entity.DepartmentEntity
	if err := r.getDB(ctx).
		Table(deptTable+" AS d").
		Joins(fmt.Sprintf(`JOIN %s dp ON dp.department_id = d.id`, joinTable)).
		Joins(fmt.Sprintf(`JOIN %s p ON p.id = dp.permission_id`, permTable)).
		Where("p.name = ?", permName).
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.
				Table(userTable).
				Where("is_active = true")
		}).
		Find(&departmentEntities).Error; err != nil {
		return nil, err
	}

	return mapper.ToDepartments(departmentEntities), nil
}

func (r *departmentRepositoryImpl) applyFilters(db *gorm.DB, query model.DepartmentQuery) *gorm.DB {
	if query.Search != "" {
		term := "%" + query.Search + "%"
		db = db.Where("name ILIKE ?", term)
	}

	return db
}
