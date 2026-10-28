package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"slices"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type userUseCaseImpl struct {
	txPro            port.TxProvider
	mqPro            port.MessageQueueProvider
	logPro           port.LogProvider
	idPro            port.IDProvider
	cachePro         port.CacheProvider
	userRepo         repository.UserRepository
	deptRepo         repository.DepartmentRepository
	branchRepo       repository.BranchRepository
	refreshTokenRepo repository.RefreshTokenRepository
}

func NewUserUseCase(
	txPro port.TxProvider,
	mqPro port.MessageQueueProvider,
	logPro port.LogProvider,
	idPro port.IDProvider,
	cachePro port.CacheProvider,
	userRepo repository.UserRepository,
	deptRepo repository.DepartmentRepository,
	branchRepo repository.BranchRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
) UserUseCase {
	return &userUseCaseImpl{
		txPro,
		mqPro,
		logPro,
		idPro,
		cachePro,
		userRepo,
		deptRepo,
		branchRepo,
		refreshTokenRepo,
	}
}

func (u *userUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateUserRequest) (int64, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		u.logPro.Error("Hash password failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate user id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	user := model.NewUser(id, req.Username, req.Email, req.Role, req.FirstName, req.LastName, req.Phone, hashedPassword, *req.IsActive, req.Avatar, req.DepartmentID, &userID, &userID)
	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return 0, err
	}
	user.Branches = branchesFromIDs(branchIDs)

	if err = u.userRepo.Create(ctx, schema, user); err != nil {
		if errors.Is(err, exceptions.ErrEmailAlreadyExists) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrUsernameAlreadyExists) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrPhoneAlreadyExists) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrDepartmentNotFound) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return 0, err
		}
		u.logPro.Error("Create user failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	return id, nil
}

func (u *userUseCaseImpl) GetByID(ctx context.Context, schema string, userID int64) (*response.UserDetailsResponse, error) {
	user, err := u.userRepo.FindByIDWithDetails(ctx, schema, userID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}
	if user == nil {
		return nil, exceptions.ErrUserNotFound
	}

	return mapper.ToUserDetailsResponse(user), nil
}

func (u *userUseCaseImpl) List(ctx context.Context, schema string, query request.UserPaginationQuery) ([]*response.SimpleUserResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	users, total, err := u.userRepo.FindAllWithDepartmentPaginated(ctx, schema, mapper.ToUserQuery(query))
	if err != nil {
		u.logPro.Error("Find all user paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSimpleUsersResponse(users), meta, nil
}

func (u *userUseCaseImpl) Update(ctx context.Context, schema string, userID, currentUserID int64, req request.UpdateUserRequest) error {
	if userID == currentUserID && (!*req.IsActive || req.Role == model.RoleStaff) {
		exists, err := u.userRepo.ExistsActiveAdminExceptID(ctx, schema, userID)
		if err != nil {
			u.logPro.Error("Check active admin except id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
			return err
		}
		if !exists {
			return exceptions.ErrNeedAdmin
		}
	}

	user, err := u.userRepo.FindByID(ctx, schema, userID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}
	if user == nil {
		return exceptions.ErrUserNotFound
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return err
	}

	updateData := map[string]any{
		"username":      req.Username,
		"email":         req.Email,
		"phone":         req.Phone,
		"first_name":    req.FirstName,
		"last_name":     req.LastName,
		"avatar":        req.Avatar,
		"role":          req.Role,
		"is_active":     *req.IsActive,
		"department_id": req.DepartmentID,
		"updated_by_id": currentUserID,
	}

	departmentChanged := req.DepartmentID != nil && user.DepartmentID != nil && *req.DepartmentID != *user.DepartmentID ||
		(req.DepartmentID == nil) != (user.DepartmentID == nil)
	branchChanged := !hasSameBranchIDs(user.Branches, branchIDs)

	if err := u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err := u.userRepo.Update(ctx, schema, userID, updateData, branchIDs); err != nil {
			if errors.Is(err, exceptions.ErrUserNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrEmailAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrUsernameAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrPhoneAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrDepartmentNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrBranchNotFound) {
				return err
			}
			u.logPro.Error("Update user failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
			return err
		}

		if !*req.IsActive || departmentChanged || branchChanged || req.Role != user.Role {
			if err := u.refreshTokenRepo.DeleteAllByUserID(ctx, schema, userID); err != nil {
				u.logPro.Error("Delete all refresh token by user id failed", port.String("schema", schema), port.Error(err))
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	if !*req.IsActive || departmentChanged || branchChanged || req.Role != user.Role {
		cacheKey := fmt.Sprintf("%s:user_version:%d", schema, userID)
		if err := u.cachePro.Increment(ctx, cacheKey); err != nil {
			u.logPro.Error("Increase token version failed", port.String("schema", schema), port.Error(err))
		}
	}

	oldAvatar := user.Avatar
	newAvatar := req.Avatar
	shouldDelete := false
	if oldAvatar != nil {
		if newAvatar == nil {
			shouldDelete = true
		} else if *oldAvatar != *newAvatar {
			shouldDelete = true
		}
	}

	if shouldDelete {
		u.deleteAvatars(schema, []string{*oldAvatar})
	}

	return nil
}

func (u *userUseCaseImpl) UpdatePassword(ctx context.Context, schema string, userID, currentUserID int64, newPassword string) error {
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		u.logPro.Error("Hash password failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateData := map[string]any{
			"password":      hashedPassword,
			"updated_by_id": currentUserID,
		}

		if err = u.userRepo.Update(ctx, schema, userID, updateData); err != nil {
			if errors.Is(err, exceptions.ErrUserNotFound) {
				return err
			}
			u.logPro.Error("Update password failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
			return err
		}

		if err := u.refreshTokenRepo.DeleteAllByUserID(ctx, schema, userID); err != nil {
			u.logPro.Error("Delete all refresh token by user id failed", port.String("schema", schema), port.Error(err))
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	cacheKey := fmt.Sprintf("%s:user_version:%d", schema, userID)
	if err = u.cachePro.Increment(ctx, cacheKey); err != nil {
		u.logPro.Error("Increase token version failed", port.String("schema", schema), port.Error(err))
	}

	return nil
}

func (u *userUseCaseImpl) Delete(ctx context.Context, schema string, userID, currentUserID int64) error {
	if userID == currentUserID {
		exists, err := u.userRepo.ExistsActiveAdminExceptID(ctx, schema, userID)
		if err != nil {
			u.logPro.Error("Check active admin except id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
			return err
		}
		if !exists {
			return exceptions.ErrNeedAdmin
		}
	}

	user, err := u.userRepo.FindByID(ctx, schema, userID)
	if err != nil {
		u.logPro.Error("Find user by id failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}
	if user == nil {
		return exceptions.ErrUserNotFound
	}

	if err := u.userRepo.Delete(ctx, schema, userID); err != nil {
		if errors.Is(err, exceptions.ErrUserNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete user failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	cacheKey := fmt.Sprintf("%s:user_version:%d", schema, userID)
	if err := u.cachePro.Del(ctx, cacheKey); err != nil {
		u.logPro.Error("Delete user version failed", port.String("schema", schema), port.Error(err))
	}

	if user.Avatar != nil {
		u.deleteAvatars(schema, []string{*user.Avatar})
	}

	return nil
}

func (u *userUseCaseImpl) DeleteBatch(ctx context.Context, schema string, currentUserID int64, userIDs []int64) (int64, error) {
	if slices.Contains(userIDs, currentUserID) {
		exists, err := u.userRepo.ExistsActiveAdminExceptID(ctx, schema, currentUserID)
		if err != nil {
			u.logPro.Error("Check active admin except id failed", port.String("schema", schema), port.Int64("current_user_id", currentUserID), port.Error(err))
			return 0, err
		}
		if !exists {
			return 0, exceptions.ErrNeedAdmin
		}
	}

	users, err := u.userRepo.FindByIDs(ctx, schema, userIDs)
	if err != nil {
		u.logPro.Error("Find user by ids failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	rowsDeleted, err := u.userRepo.DeleteByIDs(ctx, schema, userIDs)
	if err != nil {
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return 0, err
		}
		u.logPro.Error("Delete users failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	avtKeys := make([]string, 0, len(users))

	for _, user := range users {
		cacheKey := fmt.Sprintf("%s:user_version:%d", schema, user.ID)
		if err := u.cachePro.Del(ctx, cacheKey); err != nil {
			u.logPro.Error("Delete user version failed", port.String("schema", schema), port.Error(err))
		}

		avtKeys = append(avtKeys, *user.Avatar)
	}

	if len(avtKeys) > 0 {
		u.deleteAvatars(schema, avtKeys)
	}

	return rowsDeleted, nil
}

func (u *userUseCaseImpl) deleteAvatars(schema string, avatarKeys []string) {
	deleteFilesMsg := data.DeleteFileMessage{
		Schema: schema,
		Keys:   avatarKeys,
	}

	go func(msg data.DeleteFileMessage) {
		body, err := json.Marshal(msg)
		if err != nil {
			u.logPro.Error("Json marshal failed", port.String("schema", schema), port.Error(err))
		}

		if err = u.mqPro.PublishMessage(port.ExchangeFile, port.RoutingKeyDeleteFile, body); err != nil {
			u.logPro.Error("Publish delete files message failed", port.String("schema", schema), port.Error(err))
		}
	}(deleteFilesMsg)
}

func (u *userUseCaseImpl) ensureBranchesExist(ctx context.Context, schema string, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return nil
	}

	count, err := u.branchRepo.CountByIDs(ctx, schema, branchIDs)
	if err != nil {
		u.logPro.Error("Count branch by ids failed", port.String("schema", schema), port.Any("branch_ids", branchIDs), port.Error(err))
		return err
	}
	if count != int64(len(branchIDs)) {
		return exceptions.ErrBranchNotFound
	}

	return nil
}

func branchesFromIDs(branchIDs []int64) []*model.Branch {
	if len(branchIDs) == 0 {
		return make([]*model.Branch, 0)
	}

	branches := make([]*model.Branch, 0, len(branchIDs))
	for _, branchID := range branchIDs {
		branches = append(branches, &model.Branch{
			Base: model.Base{ID: branchID},
		})
	}

	return branches
}

func uniqueInt64s(values []int64) []int64 {
	if len(values) == 0 {
		return make([]int64, 0)
	}

	seen := make(map[int64]struct{}, len(values))
	res := make([]int64, 0, len(values))
	for _, value := range values {
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		res = append(res, value)
	}

	return res
}

func hasSameBranchIDs(branches []*model.Branch, branchIDs []int64) bool {
	if len(branches) != len(branchIDs) {
		return false
	}

	seen := make(map[int64]struct{}, len(branches))
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		seen[branch.ID] = struct{}{}
	}

	for _, branchID := range branchIDs {
		if _, ok := seen[branchID]; !ok {
			return false
		}
	}

	return true
}
