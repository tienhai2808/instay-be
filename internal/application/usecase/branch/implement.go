package branch

import (
	"context"
	"errors"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type branchUseCaseImpl struct {
	logPro     port.LogProvider
	idPro      port.IDProvider
	branchRepo repository.BranchRepository
}

func NewBranchUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	branchRepo repository.BranchRepository,
) BranchUseCase {
	return &branchUseCaseImpl{
		logPro,
		idPro,
		branchRepo,
	}
}

func (u *branchUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateBranchRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate branch id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	branch := model.NewBranch(id, req.Name, req.Address, req.Ward, req.City, *req.IsActive, &userID, &userID)

	if err = u.branchRepo.Create(ctx, schema, branch); err != nil {
		if errors.Is(err, exceptions.ErrBranchAlreadyExists) {
			return err
		}
		u.logPro.Error("Create branch failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *branchUseCaseImpl) ListBasic(ctx context.Context, schema string) ([]*response.BasicBranchResponse, error) {
	branches, err := u.branchRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all branch failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicBranchesResponse(branches), nil
}

func (u *branchUseCaseImpl) ListFull(ctx context.Context, schema string, query request.BranchPaginationQuery) ([]*response.BranchResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	branches, total, err := u.branchRepo.FindAllWithDetailsPaginated(ctx, schema, mapper.ToBranchQuery(query))
	if err != nil {
		u.logPro.Error("Find all branch paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToBranchesResponse(branches), meta, nil
}

func (u *branchUseCaseImpl) Update(ctx context.Context, schema string, branchID, userID int64, req request.UpdateBranchRequest) error {
	updateData := map[string]any{
		"name":          req.Name,
		"address":       req.Address,
		"ward":          req.Ward,
		"city":          req.City,
		"is_active":     *req.IsActive,
		"updated_by_id": userID,
	}

	if err := u.branchRepo.Update(ctx, schema, branchID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrBranchAlreadyExists) {
			return err
		}
		u.logPro.Error("Update branch failed", port.String("schema", schema), port.Int64("branch_id", branchID), port.Error(err))
		return err
	}

	return nil
}

func (u *branchUseCaseImpl) Delete(ctx context.Context, schema string, branchID int64) error {
	if err := u.branchRepo.Delete(ctx, schema, branchID); err != nil {
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete branch failed", port.String("schema", schema), port.Int64("branch_id", branchID), port.Error(err))
		return err
	}

	return nil
}
