package department

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type departmentUseCaseImpl struct {
	logPro                   port.LogProvider
	idPro                    port.IDProvider
	txPro                    port.TxProvider
	cachePro                 port.CacheProvider
	departmentRepo           repository.DepartmentRepository
	refreshTokenRepo         repository.RefreshTokenRepository
	permissionRepo           repository.PermissionRepository
	userRepo                 repository.UserRepository
	departmentPermissionRepo repository.DepartmentPermissionRepository
}

func NewDepartmentUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	cachePro port.CacheProvider,
	departmentRepo repository.DepartmentRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	permissionRepo repository.PermissionRepository,
	userRepo repository.UserRepository,
	departmentPermissionRepo repository.DepartmentPermissionRepository,
) DepartmentUseCase {
	return &departmentUseCaseImpl{
		logPro,
		idPro,
		txPro,
		cachePro,
		departmentRepo,
		refreshTokenRepo,
		permissionRepo,
		userRepo,
		departmentPermissionRepo,
	}
}

func (u *departmentUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateDepartmentRequest) error {
	permissions, err := u.permissionRepo.FindByIDs(ctx, schema, req.PermissionIDs)
	if err != nil {
		u.logPro.Error("Find permissions by ids failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if len(permissions) != len(req.PermissionIDs) {
		return exceptions.ErrHasPermissionNotFound
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate department id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	dept := model.NewDepartment(id, req.Name, req.Phone, req.Description, &userID, &userID, permissions)

	dps := make([]*model.DepartmentPermission, 0, len(permissions))
	deptPerms := make([]string, 0, len(permissions))
	for _, p := range permissions {
		dps = append(dps, &model.DepartmentPermission{
			DepartmentID: id,
			PermissionID: p.ID,
		})
		deptPerms = append(deptPerms, string(p.Name))
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.departmentRepo.Create(ctx, schema, dept); err != nil {
			if errors.Is(err, exceptions.ErrNameAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrPhoneAlreadyExists) {
				return err
			}
			u.logPro.Error("Create department failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if len(dps) > 0 {
			if err = u.departmentPermissionRepo.CreateBatch(ctx, schema, dps); err != nil {
				u.logPro.Error("Create batch department permission failed", port.String("schema", schema), port.Error(err))
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("%s:dept_perms:%d", schema, id)

	bytes, err := json.Marshal(deptPerms)
	if err != nil {
		u.logPro.Error("Json marshal permission data failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if err := u.cachePro.SetObject(ctx, cacheKey, bytes, 0); err != nil {
		u.logPro.Error("Save permission data failed", port.String("schema", schema), port.Error(err))
	}

	return nil
}

func (u *departmentUseCaseImpl) ListBasic(ctx context.Context, schema string) ([]*response.BasicDepartmentResponse, error) {
	depts, err := u.departmentRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all department failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicDepartmentsResponse(depts), nil
}

func (u *departmentUseCaseImpl) GetByName(ctx context.Context, schema, departmentName string) (*response.BasicDepartmentResponse, error) {
	dept, err := u.departmentRepo.FindByName(ctx, schema, departmentName)
	if err != nil {
		u.logPro.Error("Find department by name failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if dept == nil {
		return nil, exceptions.ErrDepartmentNotFound
	}

	return mapper.ToBasicDepartmentResponse(dept), nil
}

func (u *departmentUseCaseImpl) GetByID(ctx context.Context, schema string, departmentID int64) (*response.BasicDepartmentResponse, error) {
	dept, err := u.departmentRepo.FindByID(ctx, schema, departmentID)
	if err != nil {
		u.logPro.Error("Find department by id failed", port.String("schema", schema), port.Int64("department_id", departmentID), port.Error(err))
		return nil, err
	}
	if dept == nil {
		return nil, exceptions.ErrDepartmentNotFound
	}

	return mapper.ToBasicDepartmentResponse(dept), nil
}

func (u *departmentUseCaseImpl) ListFull(ctx context.Context, schema string, query request.DepartmentPaginationQuery) ([]*response.DepartmentResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	depts, total, err := u.departmentRepo.FindAllWithDetailsPaginated(ctx, schema, mapper.ToDepartmentQuery(query))
	if err != nil {
		u.logPro.Error("Find all department paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	if len(depts) == 0 {
		return make([]*response.DepartmentResponse, 0), meta, nil
	}

	deptIDs := make([]int64, len(depts))
	for i, dept := range depts {
		deptIDs[i] = dept.ID
	}

	usersCounts, err := u.departmentRepo.CountUserByIDs(ctx, schema, deptIDs)
	if err != nil {
		u.logPro.Error("Count user by department id failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	for _, dept := range depts {
		dept.UsersCount = usersCounts[dept.ID]
	}

	return mapper.ToDepartmentsResponse(depts), meta, nil
}

func (u *departmentUseCaseImpl) Update(ctx context.Context, schema string, departmentID, userID int64, req request.UpdateDepartmentRequest) error {
	dept, err := u.departmentRepo.FindByIDWithPermissions(ctx, schema, departmentID)
	if err != nil {
		u.logPro.Error("Find department by id failed", port.String("schema", schema), port.Int64("department_id", departmentID), port.Error(err))
		return err
	}
	if dept == nil {
		return exceptions.ErrDepartmentNotFound
	}

	currentPermIDs := make([]int64, 0, len(dept.Permissions))
	for _, perm := range dept.Permissions {
		currentPermIDs = append(currentPermIDs, perm.ID)
	}

	var permissions []*model.Permission
	if u.isPermissionsChanged(currentPermIDs, req.PermissionIDs) {
		permissions, err = u.permissionRepo.FindByIDs(ctx, schema, req.PermissionIDs)
		if err != nil {
			u.logPro.Error("Find permissions by ids failed", port.String("schema", schema), port.Error(err))
			return err
		}
		if len(permissions) != len(req.PermissionIDs) {
			return exceptions.ErrHasPermissionNotFound
		}
	}

	users, err := u.userRepo.FindAllActiveByDepartmentID(ctx, schema, departmentID)
	if err != nil {
		u.logPro.Error("Find all active user by department id failed", port.String("schema", schema), port.Int64("department_id", departmentID), port.Error(err))
		return err
	}

	userIDs := make([]int64, 0, 8)
	if len(users) > 0 {
		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
	}

	updateData := map[string]any{
		"name":          req.Name,
		"phone":         req.Phone,
		"description":   req.Description,
		"updated_by_id": userID,
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.departmentRepo.Update(ctx, schema, departmentID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrDepartmentNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrNameAlreadyExists) {
				return err
			}
			u.logPro.Error("Update department failed", port.String("schema", schema), port.Int64("department_id", departmentID), port.Error(err))
			return err
		}

		if u.isPermissionsChanged(currentPermIDs, req.PermissionIDs) {
			if err = u.departmentRepo.UpdatePermissions(ctx, schema, dept, permissions); err != nil {
				u.logPro.Error("Update permissions for department failed", port.String("schema", schema), port.String("department", dept.Name), port.Error(err))
				return err
			}
			if len(userIDs) > 0 {
				if err := u.refreshTokenRepo.DeleteAllByUserIDs(ctx, schema, userIDs); err != nil {
					u.logPro.Error("Delete all refresh token by user ids failed", port.String("schema", schema), port.Error(err))
					return err
				}
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if len(userIDs) > 0 {
		for _, id := range userIDs {
			cacheKey := fmt.Sprintf("%s:user_version:%d", schema, id)
			if err = u.cachePro.Increment(ctx, cacheKey); err != nil {
				u.logPro.Error("Increase token version failed", port.String("schema", schema), port.Error(err))
			}
		}
	}

	if u.isPermissionsChanged(currentPermIDs, req.PermissionIDs) {
		deptPerms := make([]string, 0, len(permissions))
		for _, perm := range permissions {
			deptPerms = append(deptPerms, string(perm.Name))
		}

		cacheKey := fmt.Sprintf("%s:dept_perms:%d", schema, departmentID)

		bytes, err := json.Marshal(deptPerms)
		if err != nil {
			u.logPro.Error("Json marshal permission data failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if err := u.cachePro.SetObject(ctx, cacheKey, bytes, 0); err != nil {
			u.logPro.Error("Save permission data failed", port.String("schema", schema), port.Error(err))
		}
	}

	return nil
}

func (u *departmentUseCaseImpl) Delete(ctx context.Context, schema string, departmentID int64) error {
	if err := u.departmentRepo.Delete(ctx, schema, departmentID); err != nil {
		if errors.Is(err, exceptions.ErrDepartmentNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete department failed", port.String("schema", schema), port.Int64("department_id", departmentID), port.Error(err))
		return err
	}

	return nil
}

func (u *departmentUseCaseImpl) isPermissionsChanged(oldPermIDs, newPermIDs []int64) bool {
	if len(oldPermIDs) != len(newPermIDs) {
		return true
	}

	m := make(map[int64]struct{}, len(oldPermIDs))

	for _, id := range oldPermIDs {
		m[id] = struct{}{}
	}

	for _, id := range newPermIDs {
		if _, ok := m[id]; !ok {
			return true
		}
	}

	return false
}
