package customer

import (
	"context"
	"errors"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type customerUseCaseImpl struct {
	logPro          port.LogProvider
	idPro           port.IDProvider
	customerRepo    repository.CustomerRepository
	bookingRoomRepo repository.BookingRoomRepository
	userRepo        repository.UserRepository
}

func NewCustomerUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	customerRepo repository.CustomerRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	userRepo repository.UserRepository,
) CustomerUseCase {
	return &customerUseCaseImpl{
		logPro,
		idPro,
		customerRepo,
		bookingRoomRepo,
		userRepo,
	}
}

func (u *customerUseCaseImpl) Create(ctx context.Context, schema string, bookingRoomID int64, req request.CreateCustomerRequest) error {
	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate customer id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	customer := model.NewCustomer(
		id,
		bookingRoomID,
		req.BirthYear,
		req.Gender,
		req.Phone,
	)

	if err = u.customerRepo.Create(ctx, schema, customer); err != nil {
		if errors.Is(err, exceptions.ErrBookingRoomNotFound) {
			return err
		}
		u.logPro.Error("Create customer failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *customerUseCaseImpl) List(ctx context.Context, schema string, userID int64, query request.CustomerPaginationQuery) ([]*response.CustomerResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	customerQuery, err := u.scopedCustomerQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	customers, total, err := u.customerRepo.FindAllPaginated(ctx, schema, customerQuery)
	if err != nil {
		u.logPro.Error("Find all customer paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToCustomersResponse(customers), meta, nil
}

func (u *customerUseCaseImpl) Export(ctx context.Context, schema string, userID int64, req request.ExportListCustomersRequest) ([]*response.CustomerResponse, error) {
	query := toCustomerExportQuery(req)
	branchIDs, err := u.customerBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}
	query.BranchIDs = branchIDs

	customers, err := u.customerRepo.FindAll(ctx, schema, query)
	if err != nil {
		u.logPro.Error("Find all customer for export failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToCustomersResponse(customers), nil
}

func (u *customerUseCaseImpl) scopedCustomerQuery(ctx context.Context, schema string, userID int64, query request.CustomerPaginationQuery) (model.CustomerQuery, error) {
	branchIDs, err := u.customerBranchIDs(ctx, schema, userID)
	if err != nil {
		return model.CustomerQuery{}, err
	}

	customerQuery := mapper.ToCustomerQuery(query)
	customerQuery.BranchIDs = branchIDs

	return customerQuery, nil
}

func (u *customerUseCaseImpl) customerBranchIDs(ctx context.Context, schema string, userID int64) ([]int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}
	if scope.IsGlobal() {
		return nil, nil
	}

	return scope.BranchIDs, nil
}

func toCustomerExportQuery(req request.ExportListCustomersRequest) model.CustomerQuery {
	query := model.CustomerQuery{
		Sort:  "created_at",
		Order: "desc",
	}

	if req.Sort != nil && strings.TrimSpace(*req.Sort) != "" {
		query.Sort = strings.TrimSpace(*req.Sort)
	}

	if req.Order != nil && strings.TrimSpace(*req.Order) != "" {
		query.Order = strings.TrimSpace(*req.Order)
	}

	return query
}
