package supportcategory

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

type supportCategoryUseCaseImpl struct {
	logPro              port.LogProvider
	idPro               port.IDProvider
	supportCategoryRepo repository.SupportCategoryRepository
}

func NewSupportCategoryUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	supportCategoryRepo repository.SupportCategoryRepository,
) SupportCategoryUseCase {
	return &supportCategoryUseCaseImpl{
		logPro,
		idPro,
		supportCategoryRepo,
	}
}

func (u *supportCategoryUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateSupportCategoryRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate support category id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	supCat := model.NewSupportCategory(id, req.DepartmentID, req.Name, &userID, &userID)

	if err = u.supportCategoryRepo.Create(ctx, schema, supCat); err != nil {
		if errors.Is(err, exceptions.ErrSupportCategoryAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrDepartmentNotFound) {
			return err
		}
		u.logPro.Error("Create support category failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *supportCategoryUseCaseImpl) Update(ctx context.Context, schema string, supportCategoryID, userID int64, req request.UpdateSupportCategoryRequest) error {
	updateData := map[string]any{
		"name":          req.Name,
		"department_id": req.DepartmentID,
		"updated_by_id": userID,
	}

	if err := u.supportCategoryRepo.Update(ctx, schema, supportCategoryID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrSupportCategoryNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrSupportCategoryAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrDepartmentNotFound) {
			return err
		}
		u.logPro.Error("Update support category failed", port.String("schema", schema), port.Int64("support_category_id", supportCategoryID), port.Error(err))
		return err
	}

	return nil
}

func (u *supportCategoryUseCaseImpl) ListFullForAdmin(ctx context.Context, schema string, query request.SupportCategoryPaginationQuery) ([]*response.SupportCategoryResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	supCats, total, err := u.supportCategoryRepo.FindAllWithDetailsPaginated(ctx, schema, mapper.ToSupportCategoryQuery(query))
	if err != nil {
		u.logPro.Error("Find all support category for admin failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	if len(supCats) == 0 {
		return mapper.ToSupportCategoriesResponse(supCats), meta, nil
	}

	supCatIDs := make([]int64, len(supCats))
	for i, supCat := range supCats {
		supCatIDs[i] = supCat.ID
	}

	supTicksCounts, err := u.supportCategoryRepo.CountSupportTicketByIDs(ctx, schema, supCatIDs)
	if err != nil {
		u.logPro.Error("Count support ticket by support category id failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	for _, supCat := range supCats {
		supCat.SupportTicketsCount = supTicksCounts[supCat.ID]
	}

	return mapper.ToSupportCategoriesResponse(supCats), meta, nil
}

func (u *supportCategoryUseCaseImpl) ListBasicForAdmin(ctx context.Context, schema string) ([]*response.BasicSupportCategoryResponse, error) {
	supCats, err := u.supportCategoryRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all support category basic for admin failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicSupportCategoriesResponse(supCats), nil
}

func (u *supportCategoryUseCaseImpl) Delete(ctx context.Context, schema string, supportCategoryID int64) error {
	if err := u.supportCategoryRepo.Delete(ctx, schema, supportCategoryID); err != nil {
		if errors.Is(err, exceptions.ErrSupportCategoryNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete support category failed", port.String("schema", schema), port.Int64("support_category_id", supportCategoryID), port.Error(err))
		return err
	}

	return nil
}

func (u *supportCategoryUseCaseImpl) ListForGuest(ctx context.Context, schema string) ([]*response.BasicSupportCategoryResponse, error) {
	supCats, err := u.supportCategoryRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all support category for guest failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicSupportCategoriesResponse(supCats), nil
}
