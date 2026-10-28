package servicetype

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
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type serviceTypeUseCaseImpl struct {
	logPro          port.LogProvider
	idPro           port.IDProvider
	serviceTypeRepo repository.ServiceTypeRepository
	serviceRepo     repository.ServiceRepository
}

func NewServiceTypeUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	serviceTypeRepo repository.ServiceTypeRepository,
	serviceRepo repository.ServiceRepository,
) ServiceTypeUseCase {
	return &serviceTypeUseCaseImpl{
		logPro,
		idPro,
		serviceTypeRepo,
		serviceRepo,
	}
}

func (u *serviceTypeUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateServiceTypeRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate service type id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	svcType := model.NewServiceType(id, req.DepartmentID, req.Name, req.Icon, utils.GenerateSlug(req.Name), &userID, &userID)

	if err := u.serviceTypeRepo.Create(ctx, schema, svcType); err != nil {
		if errors.Is(err, exceptions.ErrServiceTypeAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrDepartmentNotFound) {
			return err
		}
		u.logPro.Error("Create service type failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *serviceTypeUseCaseImpl) ListFullForAdmin(ctx context.Context, schema string, query request.ServiceTypePaginationQuery) ([]*response.ServiceTypeResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	serviceTypes, total, err := u.serviceTypeRepo.FindAllWithDetailsPaginated(ctx, schema, mapper.ToServiceTypeQuery(query))
	if err != nil {
		u.logPro.Error("Find all service type for admin failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	if len(serviceTypes) == 0 {
		return mapper.ToServiceTypesResponse(serviceTypes), meta, nil
	}

	serviceTypeIDs := make([]int64, len(serviceTypes))
	for i, serviceType := range serviceTypes {
		serviceTypeIDs[i] = serviceType.ID
	}

	servicesCounts, err := u.serviceTypeRepo.CountServiceByIDs(ctx, schema, serviceTypeIDs)
	if err != nil {
		u.logPro.Error("Count service by service type id failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	for _, serviceType := range serviceTypes {
		serviceType.ServicesCount = servicesCounts[serviceType.ID]
	}

	return mapper.ToServiceTypesResponse(serviceTypes), meta, nil
}

func (u *serviceTypeUseCaseImpl) ListBasicForAdmin(ctx context.Context, schema string) ([]*response.BasicServiceTypeResponse, error) {
	serviceTypes, err := u.serviceTypeRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all service type failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicServiceTypesResponse(serviceTypes), nil
}

func (u *serviceTypeUseCaseImpl) ListWithActiveServices(ctx context.Context, schema string, query request.ServiceTypeWithActiveServicesQuery) ([]*response.SimpleServiceTypeResponse, error) {
	serviceTypes, err := u.serviceTypeRepo.FindAllWithActiveServicesDetailsFiltered(ctx, schema, mapper.ToServiceTypeWithServicesQuery(query))
	if err != nil {
		u.logPro.Error("Find all service type with active services failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToSimpleServiceTypesResponse(serviceTypes), nil
}

func (u *serviceTypeUseCaseImpl) ListServicesBySlug(ctx context.Context, schema, svcTypeSlug string, includeServiceType bool, branchID int64) (*response.BasicServiceTypeResponse, []*response.ServiceResponse, error) {
	var svcTypeRes *response.BasicServiceTypeResponse
	if includeServiceType {
		svcType, err := u.serviceTypeRepo.FindBySlug(ctx, schema, svcTypeSlug)
		if err != nil {
			u.logPro.Error("Find service type by slug failed", port.String("schema", schema), port.Error(err))
			return nil, nil, err
		}
		if svcType == nil {
			return nil, nil, exceptions.ErrServiceTypeNotFound
		}

		svcTypeRes = mapper.ToBasicServiceTypeResponse(svcType)
	}

	services, err := u.serviceRepo.FindAllActiveByServiceTypeSlugWithThumbnail(ctx, schema, svcTypeSlug, branchID)
	if err != nil {
		u.logPro.Error("Find all service by service type slug failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	return svcTypeRes, mapper.ToServicesResponse(services), nil
}

func (u *serviceTypeUseCaseImpl) Update(ctx context.Context, schema string, svcTypeID, userID int64, req request.UpdateServiceTypeRequest) error {
	updateData := map[string]any{
		"name":          req.Name,
		"icon":          req.Icon,
		"slug":          utils.GenerateSlug(req.Name),
		"department_id": req.DepartmentID,
		"updated_by_id": userID,
	}

	if err := u.serviceTypeRepo.Update(ctx, schema, svcTypeID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrServiceTypeNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrServiceTypeAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrDepartmentNotFound) {
			return err
		}
		u.logPro.Error("Update service type failed", port.String("schema", schema), port.Int64("service_type_id", svcTypeID), port.Error(err))
		return err
	}

	return nil
}

func (u *serviceTypeUseCaseImpl) Delete(ctx context.Context, schema string, svcTypeID int64) error {
	if err := u.serviceTypeRepo.Delete(ctx, schema, svcTypeID); err != nil {
		if errors.Is(err, exceptions.ErrServiceTypeNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete service type failed", port.String("schema", schema), port.Int64("service_type_id", svcTypeID), port.Error(err))
		return err
	}

	return nil
}

func (u *serviceTypeUseCaseImpl) ListForGuest(ctx context.Context, schema string) ([]*response.BasicServiceTypeResponse, error) {
	serviceTypes, err := u.serviceTypeRepo.FindAll(ctx, schema)
	if err != nil {
		u.logPro.Error("Find all service type for guest failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicServiceTypesResponse(serviceTypes), nil
}
