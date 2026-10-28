package pricingtag

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

type pricingTagUseCaseImpl struct {
	logPro         port.LogProvider
	idPro          port.IDProvider
	pricingTagRepo repository.PricingTagRepository
}

func NewPricingTagUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	pricingTagRepo repository.PricingTagRepository,
) PricingTagUseCase {
	return &pricingTagUseCaseImpl{
		logPro,
		idPro,
		pricingTagRepo,
	}
}

func (u *pricingTagUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreatePricingTagRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate pricing tag id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	pricingTag := model.NewPricingTag(id, req.Name, &userID, &userID)

	if err := u.pricingTagRepo.Create(ctx, schema, pricingTag); err != nil {
		if errors.Is(err, exceptions.ErrPricingTagAlreadyExists) {
			return err
		}
		u.logPro.Error("Create pricing tag failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *pricingTagUseCaseImpl) ListBasic(ctx context.Context, schema string) ([]*response.BasicPricingTagResponse, error) {
	pricingTags, err := u.pricingTagRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all pricing tags failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicPricingTagsResponse(pricingTags), nil
}

func (u *pricingTagUseCaseImpl) ListFull(ctx context.Context, schema string, query request.PricingTagPaginationQuery) ([]*response.PricingTagResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	pricingTags, total, err := u.pricingTagRepo.FindAllWithDetailsPaginated(ctx, schema, mapper.ToPricingTagQuery(query))
	if err != nil {
		u.logPro.Error("Find all pricing tags paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToPricingTagsResponse(pricingTags), meta, nil
}

func (u *pricingTagUseCaseImpl) Update(ctx context.Context, schema string, pricingTagID, userID int64, req request.UpdatePricingTagRequest) error {
	updateData := map[string]any{
		"name":          req.Name,
		"updated_by_id": userID,
	}

	if err := u.pricingTagRepo.Update(ctx, schema, pricingTagID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrPricingTagNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrPricingTagAlreadyExists) {
			return err
		}
		u.logPro.Error("Update pricing tag failed", port.String("schema", schema), port.Int64("pricing_tag_id", pricingTagID), port.Error(err))
		return err
	}

	return nil
}

func (u *pricingTagUseCaseImpl) Delete(ctx context.Context, schema string, pricingTagID int64) error {
	if err := u.pricingTagRepo.Delete(ctx, schema, pricingTagID); err != nil {
		if errors.Is(err, exceptions.ErrPricingTagNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete pricing tag failed", port.String("schema", schema), port.Int64("pricing_tag_id", pricingTagID), port.Error(err))
		return err
	}

	return nil
}
