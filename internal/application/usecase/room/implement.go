package room

import (
	"context"
	"errors"
	"time"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/application/usecase/branchscope"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type roomUseCaseImpl struct {
	logPro       port.LogProvider
	idPro        port.IDProvider
	roomRepo     repository.RoomRepository
	floorRepo    repository.FloorRepository
	roomTypeRepo repository.RoomTypeRepository
	userRepo     repository.UserRepository
}

func NewRoomUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	roomRepo repository.RoomRepository,
	floorRepo repository.FloorRepository,
	roomTypeRepo repository.RoomTypeRepository,
	userRepo repository.UserRepository,
) RoomUseCase {
	return &roomUseCaseImpl{
		logPro,
		idPro,
		roomRepo,
		floorRepo,
		roomTypeRepo,
		userRepo,
	}
}

func (u *roomUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateRoomRequest) error {
	branchID, err := u.ensureFloorInScope(ctx, schema, userID, req.FloorID)
	if err != nil {
		return err
	}
	if err := u.ensureRoomTypeInBranch(ctx, schema, req.RoomTypeID, branchID); err != nil {
		return err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate room id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	room := model.NewRoom(id, req.Name, req.FloorID, req.RoomTypeID, &userID, &userID)

	if err = u.roomRepo.Create(ctx, schema, room); err != nil {
		if errors.Is(err, exceptions.ErrRoomAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrFloorNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			return err
		}
		u.logPro.Error("Create room failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *roomUseCaseImpl) List(ctx context.Context, schema string, userID int64, query request.RoomPaginationQuery) ([]*response.RoomResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	roomQuery, err := u.scopedRoomQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	rooms, total, err := u.roomRepo.FindAllWithDetailsPaginated(ctx, schema, roomQuery)
	if err != nil {
		u.logPro.Error("Find all room paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToRoomsResponse(rooms), meta, nil
}

func (u *roomUseCaseImpl) Timeline(ctx context.Context, schema string, userID int64, query request.RoomTimelineQuery) (*response.RoomTimelineResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}
	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return nil, err
	}
	query.BranchID = branchID

	timelineQuery := mapper.ToRoomTimelineQuery(query)
	timelineQuery.BranchIDs = branchIDs
	timelineQuery.SnapshotAt = time.Now()
	timelineQuery.UpcomingThreshold = time.Hour

	timeline, err := u.roomRepo.FindTimeline(ctx, schema, timelineQuery)
	if err != nil {
		u.logPro.Error("Find room timeline failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToRoomTimelineResponse(timeline), nil
}

func (u *roomUseCaseImpl) Update(ctx context.Context, schema string, roomID, userID int64, req request.UpdateRoomRequest) error {
	if err := u.ensureRoomInScope(ctx, schema, userID, roomID); err != nil {
		return err
	}
	branchID, err := u.ensureFloorInScope(ctx, schema, userID, req.FloorID)
	if err != nil {
		return err
	}
	if err := u.ensureRoomTypeInBranch(ctx, schema, req.RoomTypeID, branchID); err != nil {
		return err
	}

	updateData := map[string]any{
		"name":          req.Name,
		"floor_id":      req.FloorID,
		"room_type_id":  req.RoomTypeID,
		"updated_by_id": userID,
	}

	if err := u.roomRepo.Update(ctx, schema, roomID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrRoomNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrFloorNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			return err
		}
		u.logPro.Error("Update room failed", port.String("schema", schema), port.Int64("room_id", roomID), port.Error(err))
		return err
	}

	return nil
}

func (u *roomUseCaseImpl) Delete(ctx context.Context, schema string, roomID, userID int64) error {
	if err := u.ensureRoomInScope(ctx, schema, userID, roomID); err != nil {
		return err
	}

	if err := u.roomRepo.Delete(ctx, schema, roomID); err != nil {
		if errors.Is(err, exceptions.ErrRoomNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete room type failed", port.String("schema", schema), port.Int64("room_id", roomID), port.Error(err))
		return err
	}

	return nil
}

func (u *roomUseCaseImpl) DeleteBatch(ctx context.Context, schema string, userID int64, roomIDs []int64) (int64, error) {
	for _, roomID := range roomIDs {
		if err := u.ensureRoomInScope(ctx, schema, userID, roomID); err != nil {
			return 0, err
		}
	}

	rowsDeleted, err := u.roomRepo.DeleteByIDs(ctx, schema, roomIDs)
	if err != nil {
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return 0, err
		}
		u.logPro.Error("Delete rooms failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	return rowsDeleted, nil
}

func (u *roomUseCaseImpl) scopedRoomQuery(ctx context.Context, schema string, userID int64, query request.RoomPaginationQuery) (model.RoomQuery, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return model.RoomQuery{}, err
	}

	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return model.RoomQuery{}, err
	}

	roomQuery := mapper.ToRoomQuery(query)
	roomQuery.BranchID = branchID
	roomQuery.BranchIDs = branchIDs

	return roomQuery, nil
}

func (u *roomUseCaseImpl) ensureFloorInScope(ctx context.Context, schema string, userID, floorID int64) (int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return 0, err
	}

	floor, err := u.floorRepo.FindByID(ctx, schema, floorID)
	if err != nil {
		u.logPro.Error("Find floor by id failed", port.String("schema", schema), port.Int64("floor_id", floorID), port.Error(err))
		return 0, err
	}
	if floor == nil {
		return 0, exceptions.ErrFloorNotFound
	}
	if err := scope.Ensure(floor.BranchID); err != nil {
		return 0, err
	}

	return floor.BranchID, nil
}

func (u *roomUseCaseImpl) ensureRoomInScope(ctx context.Context, schema string, userID, roomID int64) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	room, err := u.roomRepo.FindByIDWithFloor(ctx, schema, roomID)
	if err != nil {
		u.logPro.Error("Find room by id failed", port.String("schema", schema), port.Int64("room_id", roomID), port.Error(err))
		return err
	}
	if room == nil {
		return exceptions.ErrRoomNotFound
	}
	if room.Floor == nil {
		return exceptions.ErrFloorNotFound
	}

	return scope.Ensure(room.Floor.BranchID)
}

func (u *roomUseCaseImpl) ensureRoomTypeInBranch(ctx context.Context, schema string, roomTypeID, branchID int64) error {
	exists, err := u.roomTypeRepo.ExistsByIDAndBranchID(ctx, schema, roomTypeID, branchID)
	if err != nil {
		u.logPro.Error("Check room type in branch failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Int64("branch_id", branchID), port.Error(err))
		return err
	}
	if !exists {
		return exceptions.ErrRoomTypeNotInBranch
	}

	return nil
}
