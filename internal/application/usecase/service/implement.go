package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type serviceUseCaseImpl struct {
	logPro           port.LogProvider
	idPro            port.IDProvider
	txPro            port.TxProvider
	mqPro            port.MessageQueueProvider
	storagePro       port.StorageProvider
	serviceRepo      repository.ServiceRepository
	serviceImageRepo repository.ServiceImageRepository
	branchRepo       repository.BranchRepository
	userRepo         repository.UserRepository
}

func NewServiceUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	mqPro port.MessageQueueProvider,
	storagePro port.StorageProvider,
	serviceRepo repository.ServiceRepository,
	serviceImageRepo repository.ServiceImageRepository,
	branchRepo repository.BranchRepository,
	userRepo repository.UserRepository,
) ServiceUseCase {
	return &serviceUseCaseImpl{
		logPro,
		idPro,
		txPro,
		mqPro,
		storagePro,
		serviceRepo,
		serviceImageRepo,
		branchRepo,
		userRepo,
	}
}

func (u *serviceUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateServiceRequest) (int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return 0, err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate service id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	service := model.NewService(id, req.ServiceTypeID, req.Name, utils.GenerateSlug(req.Name), req.Description, req.Price, *req.IsActive, &userID, &userID)
	if err := validateAndApplyServiceDiscount(service, req.Price, req.DiscountType, req.DiscountValue, req.DiscountStartAt, req.DiscountEndAt); err != nil {
		return 0, err
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return 0, err
	}
	if err = scope.EnsureAll(branchIDs); err != nil {
		return 0, err
	}
	service.Branches = branchesFromIDs(branchIDs)

	serviceImgs := make([]*model.ServiceImage, 0, len(req.Images))
	for _, img := range req.Images {
		serviceImg := &model.ServiceImage{
			ServiceID:   service.ID,
			Key:         img.Key,
			IsThumbnail: *img.IsThumbnail,
			SortOrder:   img.SortOrder,
		}
		serviceImgs = append(serviceImgs, serviceImg)
	}

	if err := u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.serviceRepo.Create(ctx, schema, service); err != nil {
			if errors.Is(err, exceptions.ErrServiceAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrServiceTypeNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrBranchNotFound) {
				return err
			}
			u.logPro.Error("Create service failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if err = u.serviceImageRepo.CreateBatch(ctx, schema, serviceImgs); err != nil {
			u.logPro.Error("Create batch service image failed", port.String("schema", schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return service.ID, nil
}

func (u *serviceUseCaseImpl) GetBySlug(ctx context.Context, schema string, serviceSlug string) (*response.ExtendedServiceResponse, error) {
	service, err := u.serviceRepo.FindBySlugWithServiceTypeAndServiceImages(ctx, schema, serviceSlug)
	if err != nil {
		u.logPro.Error("Find service by slug failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if service == nil {
		return nil, exceptions.ErrServiceNotFound
	}

	return mapper.ToExtendedServiceResponse(service), nil
}

func (u *serviceUseCaseImpl) Update(ctx context.Context, schema string, serviceID, userID int64, req request.UpdateServiceRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	service, err := u.serviceRepo.FindByIDWithServiceImages(ctx, schema, serviceID)
	if err != nil {
		u.logPro.Error("Find service by id with service images failed",
			port.String("schema", schema),
			port.Int64("service_id", serviceID),
			port.Error(err))
		return err
	}
	if service == nil {
		return exceptions.ErrServiceNotFound
	}
	serviceWithDetails, err := u.serviceRepo.FindByIDWithDetails(ctx, schema, serviceID)
	if err != nil {
		u.logPro.Error("Find service by id with details failed", port.String("schema", schema), port.Int64("service_id", serviceID), port.Error(err))
		return err
	}
	if err := scope.EnsureAll(branchIDsFromBranches(serviceWithDetails.Branches)); err != nil {
		return err
	}

	if err := validateAndApplyServiceDiscount(service, req.Price, req.DiscountType, req.DiscountValue, req.DiscountStartAt, req.DiscountEndAt); err != nil {
		return err
	}

	nextImgKeys := make(map[string]struct{}, len(req.Images))
	for _, img := range req.Images {
		if _, ok := nextImgKeys[img.Key]; ok {
			return exceptions.ErrBadRequest
		}
		nextImgKeys[img.Key] = struct{}{}
	}

	deletedImgKeys := make([]string, 0, len(service.ServiceImages))
	for _, img := range service.ServiceImages {
		if _, ok := nextImgKeys[img.Key]; !ok {
			deletedImgKeys = append(deletedImgKeys, img.Key)
		}
	}

	serviceImgs := make([]*model.ServiceImage, 0, len(req.Images))
	for _, img := range req.Images {
		serviceImgs = append(serviceImgs, &model.ServiceImage{
			ServiceID:   serviceID,
			Key:         img.Key,
			IsThumbnail: *img.IsThumbnail,
			SortOrder:   img.SortOrder,
		})
	}

	updateData := map[string]any{
		"name":              req.Name,
		"slug":              utils.GenerateSlug(req.Name),
		"description":       req.Description,
		"price":             req.Price,
		"is_active":         *req.IsActive,
		"discount_type":     req.DiscountType,
		"discount_value":    req.DiscountValue,
		"discount_start_at": req.DiscountStartAt,
		"discount_end_at":   req.DiscountEndAt,
		"service_type_id":   req.ServiceTypeID,
		"updated_by_id":     userID,
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return err
	}
	if err = scope.EnsureAll(branchIDs); err != nil {
		return err
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.serviceRepo.Update(ctx, schema, serviceID, updateData, branchIDs); err != nil {
			if errors.Is(err, exceptions.ErrServiceNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrServiceAlreadyExists) {
				return err
			}
			if errors.Is(err, exceptions.ErrServiceTypeNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrBranchNotFound) {
				return err
			}
			u.logPro.Error("Update service failed", port.String("schema", schema), port.Int64("service_id", serviceID), port.Error(err))
			return err
		}

		if err = u.serviceImageRepo.DeleteByServiceID(ctx, schema, serviceID); err != nil {
			u.logPro.Error("Delete service images by service id failed", port.String("schema", schema), port.Int64("service_id", serviceID), port.Error(err))
			return err
		}

		if err = u.serviceImageRepo.CreateBatch(ctx, schema, serviceImgs); err != nil {
			u.logPro.Error("Create batch service image failed", port.String("schema", schema), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	if len(deletedImgKeys) > 0 {
		u.deleteImages(schema, deletedImgKeys)
	}

	return nil
}

func validateAndApplyServiceDiscount(
	svc *model.Service,
	basePrice float64,
	discountType *model.ServiceDiscountType,
	discountValue *float64,
	discountStartAt *time.Time,
	discountEndAt *time.Time,
) error {
	if svc == nil {
		return exceptions.ErrBadRequest
	}
	if basePrice <= 0 {
		return exceptions.ErrBadRequest
	}

	hasAny := discountType != nil || discountValue != nil || discountStartAt != nil || discountEndAt != nil
	if !hasAny {
		svc.DiscountType = nil
		svc.DiscountValue = nil
		svc.DiscountStartAt = nil
		svc.DiscountEndAt = nil
		return nil
	}

	if discountType == nil || discountValue == nil {
		return exceptions.ErrBadRequest
	}

	if *discountValue <= 0 {
		return exceptions.ErrBadRequest
	}

	switch *discountType {
	case model.ServiceDiscountTypePercent:
		if *discountValue > 100 {
			return exceptions.ErrBadRequest
		}
	case model.ServiceDiscountTypeAmount:
		if *discountValue > basePrice {
			return exceptions.ErrBadRequest
		}
	default:
		return exceptions.ErrBadRequest
	}

	if discountStartAt != nil && discountEndAt != nil && discountStartAt.After(*discountEndAt) {
		return exceptions.ErrBadRequest
	}

	svc.DiscountType = discountType
	svc.DiscountValue = discountValue
	svc.DiscountStartAt = discountStartAt
	svc.DiscountEndAt = discountEndAt

	return nil
}

func (u *serviceUseCaseImpl) GetByID(ctx context.Context, schema string, userID, serviceID int64) (*response.ServiceDetailsResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}

	service, err := u.serviceRepo.FindByIDWithDetails(ctx, schema, serviceID)
	if err != nil {
		u.logPro.Error("Find service by id with details failed",
			port.String("schema", schema),
			port.Int64("service_id", serviceID),
			port.Error(err))
		return nil, err
	}
	if service == nil {
		return nil, exceptions.ErrServiceNotFound
	}
	if err := scope.EnsureAny(branchIDsFromBranches(service.Branches)); err != nil {
		return nil, err
	}

	return mapper.ToServiceDetailsResponse(service), nil
}

func (u *serviceUseCaseImpl) ListForAdmin(ctx context.Context, schema string, userID int64, query request.ServicePaginationQuery) ([]*response.SimpleServiceResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	serviceQuery, err := u.scopedServiceQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	services, total, err := u.serviceRepo.FindAllWithServiceTypePaginated(ctx, schema, serviceQuery)
	if err != nil {
		u.logPro.Error("Find all service paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSimpleServicesResponse(services), meta, nil
}

func (u *serviceUseCaseImpl) ListForGuest(ctx context.Context, schema string, query request.ServicePaginationQuery) ([]*response.ServiceResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	isActive := true
	query.IsActive = &isActive

	services, total, err := u.serviceRepo.FindAllWithServiceTypeAndThumbnailPaginated(ctx, schema, mapper.ToServiceQuery(query))
	if err != nil {
		u.logPro.Error("Find all service with service type and thumbnail paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToServicesResponse(services), meta, nil
}

func (u *serviceUseCaseImpl) Delete(ctx context.Context, schema string, serviceID, userID int64) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	service, err := u.serviceRepo.FindByIDWithServiceImages(ctx, schema, serviceID)
	if err != nil {
		u.logPro.Error("Find service by id with service images failed",
			port.String("schema", schema),
			port.Int64("service_id", serviceID),
			port.Error(err))
		return err
	}
	if service == nil {
		return exceptions.ErrServiceNotFound
	}
	serviceWithDetails, err := u.serviceRepo.FindByIDWithDetails(ctx, schema, serviceID)
	if err != nil {
		u.logPro.Error("Find service by id with details failed", port.String("schema", schema), port.Int64("service_id", serviceID), port.Error(err))
		return err
	}
	if err := scope.EnsureAll(branchIDsFromBranches(serviceWithDetails.Branches)); err != nil {
		return err
	}

	if err = u.serviceRepo.Delete(ctx, schema, serviceID); err != nil {
		if errors.Is(err, exceptions.ErrServiceNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete service failed", port.String("schema", schema), port.Int64("service_id", serviceID), port.Error(err))
		return err
	}

	if len(service.ServiceImages) > 0 {
		imgKeys := make([]string, 0, len(service.ServiceImages))
		for _, img := range service.ServiceImages {
			imgKeys = append(imgKeys, img.Key)
		}

		u.deleteImages(schema, imgKeys)
	}

	return nil
}

func (u *serviceUseCaseImpl) DeleteBatch(ctx context.Context, schema string, userID int64, serviceIDs []int64) (int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return 0, err
	}

	services, err := u.serviceRepo.FindByIDsWithServiceImages(ctx, schema, serviceIDs)
	if err != nil {
		u.logPro.Error("Find service by ids with images failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}
	for _, serviceID := range serviceIDs {
		service, err := u.serviceRepo.FindByIDWithDetails(ctx, schema, serviceID)
		if err != nil {
			u.logPro.Error("Find service by id with details failed", port.String("schema", schema), port.Int64("service_id", serviceID), port.Error(err))
			return 0, err
		}
		if service == nil {
			return 0, exceptions.ErrServiceNotFound
		}
		if err := scope.EnsureAll(branchIDsFromBranches(service.Branches)); err != nil {
			return 0, err
		}
	}

	rowDeleted, err := u.serviceRepo.DeleteByIDs(ctx, schema, serviceIDs)
	if err != nil {
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return 0, err
		}
		u.logPro.Error("Delete service by ids failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	imgKeys := make([]string, 0, 16)
	for _, svc := range services {
		for _, img := range svc.ServiceImages {
			imgKeys = append(imgKeys, img.Key)
		}
	}

	if len(imgKeys) > 0 {
		u.deleteImages(schema, imgKeys)
	}

	return rowDeleted, nil
}

func (u *serviceUseCaseImpl) deleteImages(schema string, imageKeys []string) {
	deleteFilesMsg := data.DeleteFileMessage{
		Schema: schema,
		Keys:   imageKeys,
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

func (u *serviceUseCaseImpl) ensureBranchesExist(ctx context.Context, schema string, branchIDs []int64) error {
	if len(branchIDs) == 0 {
		return exceptions.ErrBranchNotFound
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

func (u *serviceUseCaseImpl) scopedServiceQuery(ctx context.Context, schema string, userID int64, query request.ServicePaginationQuery) (model.ServiceQuery, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return model.ServiceQuery{}, err
	}

	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return model.ServiceQuery{}, err
	}

	serviceQuery := mapper.ToServiceQuery(query)
	serviceQuery.BranchID = branchID
	serviceQuery.BranchIDs = branchIDs

	return serviceQuery, nil
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

func branchIDsFromBranches(branches []*model.Branch) []int64 {
	if len(branches) == 0 {
		return make([]int64, 0)
	}

	branchIDs := make([]int64, 0, len(branches))
	for _, branch := range branches {
		if branch == nil {
			continue
		}
		branchIDs = append(branchIDs, branch.ID)
	}

	return branchIDs
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
