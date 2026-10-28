package roomtype

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	recalculateUC "github.com/InstaySystem/is_v2-be/internal/application/usecase/recalculate"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
	"github.com/InstaySystem/is_v2-be/pkg/utils"
)

type roomTypeUseCaseImpl struct {
	logPro        port.LogProvider
	idPro         port.IDProvider
	recalculateUC recalculateUC.RecalculateUseCase
	roomTypeRepo  repository.RoomTypeRepository
	branchRepo    repository.BranchRepository
	userRepo      repository.UserRepository
	dailyRateRepo repository.DailyRoomTypeRateRepository
	roomRepo      repository.RoomRepository
}

func NewRoomTypeUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	recalculateUC recalculateUC.RecalculateUseCase,
	roomTypeRepo repository.RoomTypeRepository,
	branchRepo repository.BranchRepository,
	userRepo repository.UserRepository,
	dailyRateRepo repository.DailyRoomTypeRateRepository,
	roomRepo repository.RoomRepository,
) RoomTypeUseCase {
	return &roomTypeUseCaseImpl{
		logPro,
		idPro,
		recalculateUC,
		roomTypeRepo,
		branchRepo,
		userRepo,
		dailyRateRepo,
		roomRepo,
	}
}

func (u *roomTypeUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateRoomTypeRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate room type id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return err
	}
	if err = scope.EnsureAll(branchIDs); err != nil {
		return err
	}

	roomType := model.NewRoomType(id, req.Name, req.Price, req.SurfaceArea, req.MaxPeople, &userID, &userID)
	roomType.Branches = branchesFromIDs(branchIDs)

	if err = u.roomTypeRepo.Create(ctx, schema, roomType); err != nil {
		if errors.Is(err, exceptions.ErrRoomTypeAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return err
		}
		u.logPro.Error("Create room type failed", port.String("schema", schema), port.Error(err))
		return err
	}

	u.recalculateDailyRatesAsync(schema, id)

	return nil
}

func (u *roomTypeUseCaseImpl) ListFull(ctx context.Context, schema string, userID int64, query request.RoomTypePaginationQuery) ([]*response.RoomTypeResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	roomTypeQuery, err := u.scopedRoomTypeQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	roomTypes, total, err := u.roomTypeRepo.FindAllWithDetailsPaginated(ctx, schema, roomTypeQuery)
	if err != nil {
		u.logPro.Error("Find all room type paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	if len(roomTypes) == 0 {
		return make([]*response.RoomTypeResponse, 0), meta, nil
	}

	roomTypeIDs := make([]int64, len(roomTypes))
	for i, roomType := range roomTypes {
		roomTypeIDs[i] = roomType.ID
	}

	roomsCounts, err := u.roomTypeRepo.CountRoomByIDs(ctx, schema, roomTypeIDs)
	if err != nil {
		u.logPro.Error("Count room by room type id failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	for _, roomType := range roomTypes {
		roomType.RoomsCount = roomsCounts[roomType.ID]
	}

	return mapper.ToRoomTypesResponse(roomTypes), meta, nil
}

func (u *roomTypeUseCaseImpl) ListBasic(ctx context.Context, schema string, userID int64, query request.RoomTypePaginationQuery) ([]*response.BasicRoomTypeResponse, error) {
	roomTypeQuery, err := u.scopedRoomTypeQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, err
	}

	roomTypes, err := u.roomTypeRepo.FindAll(ctx, schema, roomTypeQuery)
	if err != nil {
		u.logPro.Error("Find all room type failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBasicRoomTypesResponse(roomTypes), nil
}

func (u *roomTypeUseCaseImpl) ListRoomsByID(ctx context.Context, schema string, userID, roomTypeID int64, includeRoomType bool) (*response.BasicRoomTypeResponse, []*response.SimpleRoomResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, nil, err
	}

	var roomTypeRes *response.BasicRoomTypeResponse
	roomType, err := u.roomTypeRepo.FindByID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find room type by id failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return nil, nil, err
	}
	if roomType == nil {
		return nil, nil, exceptions.ErrRoomTypeNotFound
	}
	if err := scope.EnsureAll(branchIDsFromBranches(roomType.Branches)); err != nil {
		return nil, nil, err
	}

	if includeRoomType {
		roomTypeRes = mapper.ToBasicRoomTypeResponse(roomType)
	}

	rooms, err := u.roomRepo.FindAllByRoomTypeID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find all room by room type id failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	return roomTypeRes, mapper.ToSimpleRoomsResponse(rooms), nil
}

func (u *roomTypeUseCaseImpl) CalculatePrice(ctx context.Context, schema string, req request.CalculateRoomTypeRateRequest) (*response.RoomTypePriceCalculationResponse, error) {
	roomType, err := u.roomTypeRepo.FindByID(ctx, schema, req.RoomTypeID)
	if err != nil {
		u.logPro.Error("Find room type by id failed", port.String("schema", schema), port.Int64("room_type_id", req.RoomTypeID), port.Error(err))
		return nil, err
	}
	if roomType == nil {
		return nil, exceptions.ErrRoomTypeNotFound
	}

	totalDays := utils.CalculateStayDays(req.FromDate, req.ToDate)
	fromDate := utils.NormalizeDateOnly(req.FromDate)
	toDate := fromDate.AddDate(0, 0, int(totalDays))

	dailyTotalPrice, coveredDays, err := u.dailyRateRepo.SumFinalPriceByRoomTypeIDAndDateRange(ctx, schema, req.RoomTypeID, fromDate, toDate)
	if err != nil {
		u.logPro.Error(
			"sum daily room type rate by date range failed",
			port.String("schema", schema),
			port.Int64("room_type_id", req.RoomTypeID),
			port.Time("from_date", fromDate),
			port.Time("to_date", toDate),
			port.Error(err),
		)
		return nil, err
	}

	if coveredDays < 0 {
		coveredDays = 0
	}
	if coveredDays > totalDays {
		coveredDays = totalDays
	}

	missingDays := totalDays - coveredDays
	pricePerRoom := dailyTotalPrice + (float64(missingDays) * roomType.Price)
	totalPrice := pricePerRoom * float64(req.Quantity)
	averagePrice := totalPrice / float64(totalDays)

	return &response.RoomTypePriceCalculationResponse{
		TotalDays:    totalDays,
		TotalPrice:   totalPrice,
		AveragePrice: averagePrice,
	}, nil
}

func (u *roomTypeUseCaseImpl) Update(ctx context.Context, schema string, roomTypeID, userID int64, req request.UpdateRoomTypeRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	roomType, err := u.roomTypeRepo.FindByID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find room type by id failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}
	if roomType == nil {
		return exceptions.ErrRoomTypeNotFound
	}
	if err := scope.EnsureAll(branchIDsFromBranches(roomType.Branches)); err != nil {
		return err
	}

	branchIDs := uniqueInt64s(req.BranchIDs)
	if err = u.ensureBranchesExist(ctx, schema, branchIDs); err != nil {
		return err
	}
	if err = scope.EnsureAll(branchIDs); err != nil {
		return err
	}

	var isChangePrice bool
	if roomType.Price != req.Price {
		isChangePrice = true
	}

	updateData := map[string]any{
		"name":          req.Name,
		"price":         req.Price,
		"surface_area":  req.SurfaceArea,
		"max_people":    req.MaxPeople,
		"updated_by_id": userID,
	}

	if err := u.roomTypeRepo.Update(ctx, schema, roomTypeID, updateData, branchIDs); err != nil {
		if errors.Is(err, exceptions.ErrRoomTypeAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return err
		}
		u.logPro.Error("Update room type failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if isChangePrice {
		u.recalculateDailyRatesAsync(schema, roomTypeID)
	}

	return nil
}

func (u *roomTypeUseCaseImpl) Delete(ctx context.Context, schema string, roomTypeID, userID int64) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}
	roomType, err := u.roomTypeRepo.FindByID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find room type by id failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return err
	}
	if roomType == nil {
		return exceptions.ErrRoomTypeNotFound
	}
	if err := scope.EnsureAll(branchIDsFromBranches(roomType.Branches)); err != nil {
		return err
	}

	if err := u.roomTypeRepo.Delete(ctx, schema, roomTypeID); err != nil {
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete room type failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *roomTypeUseCaseImpl) recalculateDailyRatesAsync(schema string, roomTypeID int64) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := u.recalculateUC.RecalculateDailyRoomTypeRate(ctx, schema, roomTypeID); err != nil {
			u.logPro.Error(
				"recalculate daily room type rates failed",
				port.String("schema", schema),
				port.Int64("room_type_id", roomTypeID),
				port.Error(err),
			)
		}
	}()
}

func (u *roomTypeUseCaseImpl) ensureBranchesExist(ctx context.Context, schema string, branchIDs []int64) error {
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

func (u *roomTypeUseCaseImpl) scopedRoomTypeQuery(ctx context.Context, schema string, userID int64, query request.RoomTypePaginationQuery) (model.RoomTypeQuery, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return model.RoomTypeQuery{}, err
	}

	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return model.RoomTypeQuery{}, err
	}

	roomTypeQuery := mapper.ToRoomTypeQuery(query)
	roomTypeQuery.BranchID = branchID
	roomTypeQuery.BranchIDs = branchIDs

	return roomTypeQuery, nil
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
