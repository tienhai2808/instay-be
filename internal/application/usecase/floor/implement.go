package floor

import (
	"context"
	"errors"

	"github.com/InstaySystem/is_v2-be/internal/application/port"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type floorUseCaseImpl struct {
	logPro    port.LogProvider
	idPro     port.IDProvider
	floorRepo repository.FloorRepository
	userRepo  repository.UserRepository
}

func NewFloorUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	floorRepo repository.FloorRepository,
	userRepo repository.UserRepository,
) FloorUseCase {
	return &floorUseCaseImpl{
		logPro,
		idPro,
		floorRepo,
		userRepo,
	}
}

func (u *floorUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateFloorRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}
	if err := scope.Ensure(req.BranchID); err != nil {
		return err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate floor id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	floor := model.NewFloor(id, req.Name, req.BranchID, &userID, &userID)

	if err = u.floorRepo.Create(ctx, schema, floor); err != nil {
		if errors.Is(err, exceptions.ErrFloorAlreadyExists) {
			return err
		}
		u.logPro.Error("Create floor failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *floorUseCaseImpl) ListBasic(ctx context.Context, schema string, userID int64, query request.FloorPaginationQuery) ([]*response.BasicFloorResponse, error) {
	floorQuery, err := u.scopedFloorQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, err
	}

	floors, err := u.floorRepo.FindAll(ctx, schema, floorQuery)
	if err != nil {
		u.logPro.Error("Find all floor failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicFloorsResponse(floors), nil
}

func (u *floorUseCaseImpl) ListFull(ctx context.Context, schema string, userID int64, query request.FloorPaginationQuery) ([]*response.FloorResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	floorQuery, err := u.scopedFloorQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	floors, total, err := u.floorRepo.FindAllWithDetailsPaginated(ctx, schema, floorQuery)
	if err != nil {
		u.logPro.Error("Find all floor paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToFloorsResponse(floors), meta, nil
}

func (u *floorUseCaseImpl) Update(ctx context.Context, schema string, floorID, userID int64, req request.UpdateFloorRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	floor, err := u.floorRepo.FindByID(ctx, schema, floorID)
	if err != nil {
		u.logPro.Error("Find floor by id failed", port.String("schema", schema), port.Int64("floor_id", floorID), port.Error(err))
		return err
	}
	if floor == nil {
		return exceptions.ErrFloorNotFound
	}
	if err := scope.Ensure(floor.BranchID); err != nil {
		return err
	}
	if err := scope.Ensure(req.BranchID); err != nil {
		return err
	}

	updateData := map[string]any{
		"name":          req.Name,
		"branch_id":     req.BranchID,
		"updated_by_id": userID,
	}

	if err := u.floorRepo.Update(ctx, schema, floorID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrFloorNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrFloorAlreadyExists) {
			return err
		}
		u.logPro.Error("Update floor failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *floorUseCaseImpl) Delete(ctx context.Context, schema string, floorID, userID int64) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	floor, err := u.floorRepo.FindByID(ctx, schema, floorID)
	if err != nil {
		u.logPro.Error("Find floor by id failed", port.String("schema", schema), port.Int64("floor_id", floorID), port.Error(err))
		return err
	}
	if floor == nil {
		return exceptions.ErrFloorNotFound
	}
	if err := scope.Ensure(floor.BranchID); err != nil {
		return err
	}

	if err := u.floorRepo.Delete(ctx, schema, floorID); err != nil {
		if errors.Is(err, exceptions.ErrFloorNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete floor failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *floorUseCaseImpl) scopedFloorQuery(ctx context.Context, schema string, userID int64, query request.FloorPaginationQuery) (model.FloorQuery, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return model.FloorQuery{}, err
	}

	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return model.FloorQuery{}, err
	}

	floorQuery := mapper.ToFloorQuery(query)
	floorQuery.BranchID = branchID
	floorQuery.BranchIDs = branchIDs

	return floorQuery, nil
}
