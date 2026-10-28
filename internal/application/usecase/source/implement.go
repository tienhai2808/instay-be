package source

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

type sourceUseCaseImpl struct {
	logPro     port.LogProvider
	idPro      port.IDProvider
	sourceRepo repository.SourceRepository
}

func NewSourceUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	sourceRepo repository.SourceRepository,
) SourceUseCase {
	return &sourceUseCaseImpl{
		logPro,
		idPro,
		sourceRepo,
	}
}

func (u *sourceUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateSourceRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate source id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	src := model.NewSource(id, req.Name, &userID, &userID)

	if err = u.sourceRepo.Create(ctx, schema, src); err != nil {
		if errors.Is(err, exceptions.ErrSourceAlreadyExists) {
			return err
		}
		u.logPro.Error("Create source failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *sourceUseCaseImpl) Update(ctx context.Context, schema string, sourceID, userID int64, req request.UpdateSourceRequest) error {
	updateData := map[string]any{
		"name":          req.Name,
		"updated_by_id": userID,
	}

	if err := u.sourceRepo.Update(ctx, schema, sourceID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrSourceNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrSourceAlreadyExists) {
			return err
		}
		u.logPro.Error("Update source failed", port.String("schema", schema), port.Int64("source_id", sourceID), port.Error(err))
		return err
	}

	return nil
}

func (u *sourceUseCaseImpl) Delete(ctx context.Context, schema string, sourceID int64) error {
	if err := u.sourceRepo.Delete(ctx, schema, sourceID); err != nil {
		if errors.Is(err, exceptions.ErrSourceNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete source failed", port.String("schema", schema), port.Int64("source_id", sourceID), port.Error(err))
		return err
	}

	return nil
}

func (u *sourceUseCaseImpl) ListBasic(ctx context.Context, schema string) ([]*response.BasicSourceResponse, error) {
	srcs, err := u.sourceRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all source failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicSourcesResponse(srcs), nil
}

func (u *sourceUseCaseImpl) ListFull(ctx context.Context, schema string, query request.SourcePaginationQuery) ([]*response.SourceResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	sources, total, err := u.sourceRepo.FindAllWithDetailsPaginated(ctx, schema, mapper.ToSourceQuery(query))
	if err != nil {
		u.logPro.Error("Find all source paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSourcesResponse(sources), meta, nil
}
