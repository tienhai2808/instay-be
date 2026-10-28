package bookingroom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

type bookingRoomUseCaseImpl struct {
	logPro          port.LogProvider
	idPro           port.IDProvider
	txPro           port.TxProvider
	cachePro        port.CacheProvider
	tokenPro        port.TokenProvider
	bookingRepo     repository.BookingRepository
	roomRepo        repository.RoomRepository
	roomTypeRepo    repository.RoomTypeRepository
	userRepo        repository.UserRepository
	bookingRoomRepo repository.BookingRoomRepository
	guestRepo       repository.GuestRepository
}

func NewBookingRoomUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	cachePro port.CacheProvider,
	tokenPro port.TokenProvider,
	bookingRepo repository.BookingRepository,
	roomRepo repository.RoomRepository,
	roomTypeRepo repository.RoomTypeRepository,
	userRepo repository.UserRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	guestRepo repository.GuestRepository,
) BookingRoomUseCase {
	return &bookingRoomUseCaseImpl{
		logPro,
		idPro,
		txPro,
		cachePro,
		tokenPro,
		bookingRepo,
		roomRepo,
		roomTypeRepo,
		userRepo,
		bookingRoomRepo,
		guestRepo,
	}
}

func (u *bookingRoomUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateBookingRoomRequest) (int64, error) {
	booking, err := u.bookingRepo.FindByIDWithBookingRooms(ctx, schema, req.BookingID)
	if err != nil {
		u.logPro.Error("Find booking by id failed", port.String("schema", schema), port.Int64("booking_id", req.BookingID), port.Error(err))
		return 0, err
	}
	if booking == nil {
		return 0, exceptions.ErrBookingNotFound
	}
	if booking.Status != model.BookingStatusConfirmed {
		return 0, exceptions.ErrInvalidStatus
	}

	reservedRoomCount := uint32(0)
	for _, bookingRoom := range booking.BookingRooms {
		if bookingRoom.Status == model.BookingRoomStatusReserved {
			reservedRoomCount++
		}
	}
	if reservedRoomCount >= booking.RoomNumber {
		return 0, exceptions.ErrMaxRoomReached
	}

	room, err := u.roomRepo.FindByIDWithFloor(ctx, schema, req.RoomID)
	if err != nil {
		u.logPro.Error("Find room by id with floor failed", port.String("schema", schema), port.Int64("room_id", req.RoomID), port.Error(err))
		return 0, err
	}
	if err := u.validateRoomForBooking(ctx, schema, booking, room, 0); err != nil {
		return 0, err
	}
	if err := u.ensureRoomInUserScope(ctx, schema, userID, room); err != nil {
		return 0, err
	}

	exists, err := u.bookingRoomRepo.ExistsActiveByRoomID(ctx, schema, req.RoomID)
	if err != nil {
		u.logPro.Error("Check active booking rooms by room id failed", port.String("schema", schema), port.Int64("room_id", req.RoomID), port.Error(err))
		return 0, err
	}
	if exists {
		return 0, exceptions.ErrRoomCurrentlyOccupied
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate booking room id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	code := utils.GenerateBase58ID(16)
	bookingRoom := model.NewBookingRoom(id, req.BookingID, req.RoomID, code)

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.bookingRoomRepo.Create(ctx, schema, bookingRoom); err != nil {
			if errors.Is(err, exceptions.ErrBookingSelectedRoom) {
				return err
			}
			if errors.Is(err, exceptions.ErrRoomNotFound) {
				return err
			}
			u.logPro.Error("Create booking room failed", port.String("schema", schema), port.Error(err))
			return err
		}

		if booking.BranchID == nil {
			updateBookingData := map[string]any{
				"branch_id": room.Floor.BranchID,
			}
			if err = u.bookingRepo.Update(ctx, schema, booking.ID, updateBookingData); err != nil {
				u.logPro.Error("Update booking branch failed", port.String("schema", schema), port.Int64("booking_id", booking.ID), port.Error(err))
				return err
			}
		}

		return nil
	}); err != nil {
		return 0, err
	}

	return id, nil
}

func (u *bookingRoomUseCaseImpl) Update(ctx context.Context, schema string, bookingRoomID, userID int64, req request.UpdateBookingRoomRequest) error {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return err
	}
	if bookingRoom == nil {
		return exceptions.ErrBookingRoomNotFound
	}
	if bookingRoom.Status != model.BookingRoomStatusReserved {
		return exceptions.ErrInvalidStatus
	}

	room, err := u.roomRepo.FindByIDWithFloor(ctx, schema, req.RoomID)
	if err != nil {
		u.logPro.Error("Find room by id failed", port.String("schema", schema), port.Int64("room_id", req.RoomID), port.Error(err))
		return err
	}

	booking, err := u.bookingRepo.FindByIDWithBookingRooms(ctx, schema, bookingRoom.BookingID)
	if err != nil {
		u.logPro.Error("Find booking by id failed", port.String("schema", schema), port.Int64("booking_id", bookingRoom.BookingID), port.Error(err))
		return err
	}
	if booking == nil {
		return exceptions.ErrBookingNotFound
	}
	if err := u.validateRoomForBooking(ctx, schema, booking, room, bookingRoomID); err != nil {
		return err
	}
	if err := u.ensureRoomInUserScope(ctx, schema, userID, bookingRoom.Room); err != nil {
		return err
	}
	if err := u.ensureRoomInUserScope(ctx, schema, userID, room); err != nil {
		return err
	}

	exists, err := u.bookingRoomRepo.ExistsActiveByRoomIDExceptID(ctx, schema, bookingRoomID, req.RoomID)
	if err != nil {
		u.logPro.Error("Check active booking rooms by room id except id failed", port.String("schema", schema), port.Int64("room_id", req.RoomID), port.Error(err))
		return err
	}
	if exists {
		return exceptions.ErrRoomCurrentlyOccupied
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		updateBookingRoomData := map[string]any{
			"room_id": req.RoomID,
		}
		if err = u.bookingRoomRepo.Update(ctx, schema, bookingRoomID, updateBookingRoomData); err != nil {
			u.logPro.Error("Update booking room failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
			return err
		}

		updateBookingData := map[string]any{
			"updated_by_id": userID,
		}
		if booking.BranchID == nil {
			updateBookingData["branch_id"] = room.Floor.BranchID
		}
		if err = u.bookingRepo.Update(ctx, schema, bookingRoom.BookingID, updateBookingData); err != nil {
			u.logPro.Error("Update booking failed", port.String("schema", schema), port.Int64("booking_id", bookingRoom.BookingID), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *bookingRoomUseCaseImpl) validateRoomForBooking(ctx context.Context, schema string, booking *model.Booking, room *model.Room, excludeBookingRoomID int64) error {
	if booking == nil {
		return exceptions.ErrBookingNotFound
	}
	if room == nil || room.Floor == nil {
		return exceptions.ErrRoomNotFound
	}
	if room.RoomTypeID != booking.RoomTypeID {
		return exceptions.ErrBookingRoomsNotMatched
	}

	branchID := room.Floor.BranchID
	exists, err := u.roomTypeRepo.ExistsByIDAndBranchID(ctx, schema, room.RoomTypeID, branchID)
	if err != nil {
		u.logPro.Error("Check room type in branch failed", port.String("schema", schema), port.Int64("room_type_id", room.RoomTypeID), port.Int64("branch_id", branchID), port.Error(err))
		return err
	}
	if !exists {
		return exceptions.ErrRoomTypeNotInBranch
	}

	if booking.BranchID != nil && *booking.BranchID != branchID {
		return exceptions.ErrRoomNotInBookingBranch
	}

	for _, bookingRoom := range booking.BookingRooms {
		if bookingRoom == nil || bookingRoom.ID == excludeBookingRoomID {
			continue
		}
		if bookingRoom.Room == nil || bookingRoom.Room.Floor == nil {
			return exceptions.ErrRoomNotFound
		}
		if bookingRoom.Room.Floor.BranchID != branchID {
			return exceptions.ErrRoomNotInBookingBranch
		}
	}

	return nil
}

func (u *bookingRoomUseCaseImpl) ensureRoomInUserScope(ctx context.Context, schema string, userID int64, room *model.Room) error {
	if room == nil || room.Floor == nil {
		return exceptions.ErrRoomNotFound
	}

	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	return scope.Ensure(room.Floor.BranchID)
}

func (u *bookingRoomUseCaseImpl) Verify(ctx context.Context, schema, code string) (*response.SimpleBookingRoomResponse, string, time.Duration, error) {
	cacheKey := fmt.Sprintf("%s:booking_room:%s", schema, code)
	bytes, err := u.cachePro.GetObject(ctx, cacheKey)
	if err != nil {
		u.logPro.Error("Get booking room data failed", port.String("schema", schema), port.Error(err))
		return nil, "", 0, err
	}
	if bytes == nil {
		return nil, "", 0, exceptions.ErrInvalidToken
	}

	var bookingRoomData data.BookingRoomData
	if err = json.Unmarshal(bytes, &bookingRoomData); err != nil {
		u.logPro.Error("Json unmarshal booking room data failed", port.String("schema", schema), port.Error(err))
		return nil, "", 0, err
	}

	bookingRoom, err := u.bookingRoomRepo.FindByIDWithBookingAndRoom(ctx, schema, bookingRoomData.ID)
	if err != nil {
		u.logPro.Error("Find booking room by id with booking and room failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomData.ID), port.Error(err))
		return nil, "", 0, err
	}
	if bookingRoom == nil || bookingRoom.Booking == nil {
		return nil, "", 0, exceptions.ErrInvalidToken
	}
	if bookingRoom.Status != model.BookingRoomStatusReleased {
		return nil, "", 0, exceptions.ErrInvalidStatus
	}
	if bookingRoom.Booking.Status != model.BookingStatusCheckedIn {
		return nil, "", 0, exceptions.ErrInvalidStatus
	}

	ttl := time.Until(bookingRoomData.ExpiredAt)

	guestToken, err := u.tokenPro.GenerateGuestToken(bookingRoomData.ID, ttl)
	if err != nil {
		u.logPro.Error("Generate guest token failed", port.String("schema", schema), port.Error(err))
		return nil, "", 0, err
	}

	return mapper.ToSimpleBookingRoomResponse(bookingRoom), guestToken, ttl, nil
}

func (u *bookingRoomUseCaseImpl) Info(ctx context.Context, schema string, bookingRoomID int64) (*response.SimpleBookingRoomResponse, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithBookingAndRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id with booking and room failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return nil, err
	}
	if bookingRoom == nil {
		return nil, exceptions.ErrBookingRoomNotFound
	}

	return mapper.ToSimpleBookingRoomResponse(bookingRoom), nil
}

func (u *bookingRoomUseCaseImpl) ListGuestsByID(ctx context.Context, schema string, bookingRoomID, userID int64) ([]*response.SimpleGuestResponse, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return nil, err
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return nil, err
	}

	guests, err := u.guestRepo.FindAllByBookingRoomID(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find all guest by booking room id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return nil, err
	}

	return mapper.ToSimpleGuestsResponse(guests), nil
}

func (u *bookingRoomUseCaseImpl) GetByIDFull(ctx context.Context, schema string, bookingRoomID, userID int64) (*response.BookingRoomResponse, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithDetails(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return nil, err
	}
	if bookingRoom == nil {
		return nil, exceptions.ErrBookingRoomNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return nil, err
	}

	return mapper.ToBookingRoomResponse(bookingRoom), nil
}

func (u *bookingRoomUseCaseImpl) GetByIDBasic(ctx context.Context, schema string, bookingRoomID, userID int64) (*response.SimpleBookingRoomResponse, error) {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithBookingAndRoom(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return nil, err
	}
	if bookingRoom == nil {
		return nil, exceptions.ErrBookingRoomNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return nil, err
	}

	return mapper.ToSimpleBookingRoomResponse(bookingRoom), nil
}

func (u *bookingRoomUseCaseImpl) ensureBookingRoomInUserScope(ctx context.Context, schema string, userID int64, bookingRoom *model.BookingRoom) error {
	if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil {
		return exceptions.ErrBookingRoomNotFound
	}

	return u.ensureRoomInUserScope(ctx, schema, userID, bookingRoom.Room)
}
