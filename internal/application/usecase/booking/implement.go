package booking

import (
	"context"
	"crypto/rand"
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
	"gorm.io/datatypes"
)

type bookingUseCaseImpl struct {
	logPro          port.LogProvider
	idPro           port.IDProvider
	txPro           port.TxProvider
	cachePro        port.CacheProvider
	aiPro           port.AIProvider
	bookingRepo     repository.BookingRepository
	bookingRoomRepo repository.BookingRoomRepository
	branchRepo      repository.BranchRepository
	userRepo        repository.UserRepository
	roomRepo        repository.RoomRepository
	roomTypeRepo    repository.RoomTypeRepository
	dailyRateRepo   repository.DailyRoomTypeRateRepository
}

func NewBookingUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	txPro port.TxProvider,
	cachePro port.CacheProvider,
	aiPro port.AIProvider,
	bookingRepo repository.BookingRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	branchRepo repository.BranchRepository,
	userRepo repository.UserRepository,
	roomRepo repository.RoomRepository,
	roomTypeRepo repository.RoomTypeRepository,
	dailyRateRepo repository.DailyRoomTypeRateRepository,
) BookingUseCase {
	return &bookingUseCaseImpl{
		logPro,
		idPro,
		txPro,
		cachePro,
		aiPro,
		bookingRepo,
		bookingRoomRepo,
		branchRepo,
		userRepo,
		roomRepo,
		roomTypeRepo,
		dailyRateRepo,
	}
}

func (u *bookingUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateBookingRequest) (int64, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return 0, err
	}
	if err := u.ensureBranchPtrInScope(scope, req.BranchID); err != nil {
		return 0, err
	}

	if req.BookingNumber == nil {
		newBookingNumber, err := u.generateBookingNumber()
		if err != nil {
			u.logPro.Error("Generate booking number failed", port.String("schema", schema), port.Error(err))
			return 0, err
		}

		req.BookingNumber = &newBookingNumber
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate booking id failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	currentTime := time.Now()
	now := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, currentTime.Location())

	if req.TotalNetPrice == nil {
		totalNetPrice, err := u.calculateBookingTotalNetPrice(ctx, schema, req.RoomTypeID, req.CheckIn, req.CheckOut, req.RoomNumber)
		if err != nil {
			return 0, err
		}

		req.TotalNetPrice = &totalNetPrice
		if req.TotalSellPrice == nil {
			req.TotalSellPrice = &totalNetPrice
		}
	}

	totalNetPrice := 0.0
	if req.TotalNetPrice != nil {
		totalNetPrice = *req.TotalNetPrice
	}

	totalSellPrice := 0.0
	if req.TotalSellPrice != nil {
		totalSellPrice = *req.TotalSellPrice
	}

	if err := u.validateBookingBranch(ctx, schema, req.RoomTypeID, req.BranchID); err != nil {
		return 0, err
	}

	booking := model.NewBooking(id, *req.BookingNumber, req.GuestFullName, req.GuestNumber, &req.GuestEmail, &req.GuestPhone, req.CheckIn, req.CheckOut, now, req.RoomNumber, totalNetPrice, totalSellPrice, req.PromotionName, req.MealPlan, req.BookingPreferences, req.BookingConditions, req.Note, model.BookingStatusConfirmed, req.RoomTypeID, req.SourceID, req.BranchID, &userID, &userID)

	if err = u.bookingRepo.Create(ctx, schema, booking); err != nil {
		if errors.Is(err, exceptions.ErrBookingAlreadyExists) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrSourceNotFound) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			return 0, err
		}
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return 0, err
		}
		u.logPro.Error("Create booking failed", port.String("schema", schema), port.Error(err))
		return 0, err
	}

	u.predictBooking(schema, id)

	return id, nil
}

func (u *bookingUseCaseImpl) Update(ctx context.Context, schema string, bookingID, userID int64, req request.UpdateBookingRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	booking, err := u.bookingRepo.FindByIDWithBookingRooms(ctx, schema, bookingID)
	if err != nil {
		u.logPro.Error("Find booking by id failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
		return err
	}
	if booking == nil {
		return exceptions.ErrBookingNotFound
	}

	if booking.Status != model.BookingStatusConfirmed {
		return exceptions.ErrInvalidStatus
	}
	if err := u.ensureBranchPtrInScope(scope, booking.BranchID); err != nil {
		return err
	}

	targetBranchID := booking.BranchID
	if req.BranchID != nil {
		targetBranchID = req.BranchID
	}
	if err := u.ensureBranchPtrInScope(scope, targetBranchID); err != nil {
		return err
	}
	if err := u.validateBookingBranch(ctx, schema, req.RoomTypeID, targetBranchID); err != nil {
		return err
	}
	if err := u.validateBookingRoomsInBranch(booking, targetBranchID); err != nil {
		return err
	}

	updateData := map[string]any{
		"guest_full_name":     req.GuestFullName,
		"guest_email":         req.GuestEmail,
		"guest_phone":         req.GuestPhone,
		"check_in":            req.CheckIn,
		"check_out":           req.CheckOut,
		"room_number":         req.RoomNumber,
		"room_type_id":        req.RoomTypeID,
		"guest_number":        req.GuestNumber,
		"promotion_name":      req.PromotionName,
		"meal_plan":           req.MealPlan,
		"booking_preferences": req.BookingPreferences,
		"booking_conditions":  req.BookingConditions,
		"note":                req.Note,
		"updated_by_id":       userID,
	}
	if req.BranchID != nil {
		updateData["branch_id"] = req.BranchID
	}

	isCheckInChanged := !req.CheckIn.Equal(booking.CheckIn)
	isCheckOutChanged := !req.CheckOut.Equal(booking.CheckOut)
	isRoomTypeChanged := req.RoomTypeID != booking.RoomTypeID
	isRoomNumberChanged := req.RoomNumber != booking.RoomNumber
	if isCheckInChanged || isCheckOutChanged || isRoomTypeChanged || isRoomNumberChanged {
		totalNetPrice, err := u.calculateBookingTotalNetPrice(ctx, schema, req.RoomTypeID, req.CheckIn, req.CheckOut, req.RoomNumber)
		if err != nil {
			return err
		}

		updateData["total_net_price"] = totalNetPrice
		updateData["total_sell_price"] = totalNetPrice
	}

	if err = u.bookingRepo.Update(ctx, schema, bookingID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrBookingNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrBookingAlreadyExists) {
			return err
		}
		if errors.Is(err, exceptions.ErrSourceNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomTypeNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrBranchNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomTypeNotInBranch) {
			return err
		}
		if errors.Is(err, exceptions.ErrRoomNotInBookingBranch) {
			return err
		}
		u.logPro.Error("Update booking failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
		return err
	}

	u.predictBooking(schema, bookingID)

	return nil
}

func (u *bookingUseCaseImpl) UpdateStatus(ctx context.Context, schema string, bookingID, userID int64, req request.UpdateBookingStatusRequest) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	booking, err := u.bookingRepo.FindByIDWithBookingRooms(ctx, schema, bookingID)
	if err != nil {
		u.logPro.Error("Find booking by id failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
		return err
	}
	if booking == nil {
		return exceptions.ErrBookingNotFound
	}
	if err := u.ensureBranchPtrInScope(scope, booking.BranchID); err != nil {
		return err
	}

	switch booking.Status {
	case model.BookingStatusConfirmed:
		if req.Status != model.BookingStatusCancelled && req.Status != model.BookingStatusCheckedIn {
			return exceptions.ErrInvalidStatus
		}
	case model.BookingStatusCheckedIn:
		if req.Status != model.BookingStatusCheckedOut {
			return exceptions.ErrInvalidStatus
		}
	default:
		return exceptions.ErrInvalidStatus
	}

	updateData := map[string]any{
		"status":        req.Status,
		"updated_by_id": userID,
	}

	if req.Status == model.BookingStatusCheckedIn {
		reservedBookingRooms := make([]*model.BookingRoom, 0, len(booking.BookingRooms))
		for _, bookingRoom := range booking.BookingRooms {
			if bookingRoom.Status == model.BookingRoomStatusReserved {
				reservedBookingRooms = append(reservedBookingRooms, bookingRoom)
			}
		}

		if len(reservedBookingRooms) != int(booking.RoomNumber) {
			return exceptions.ErrBookingRoomsNotMatched
		}

		if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
			bookingRoomIDs := make([]int64, 0, 8)
			for _, bookingRoom := range reservedBookingRooms {
				bookingRoomIDs = append(bookingRoomIDs, bookingRoom.ID)
			}

			updateBookingRoomData := map[string]any{
				"status": model.BookingRoomStatusReleased,
			}
			if err = u.bookingRoomRepo.Updates(ctx, schema, bookingRoomIDs, updateBookingRoomData); err != nil {
				u.logPro.Error("Update booking rooms status failed", port.String("schema", schema), port.Error(err))
				return err
			}

			if err = u.bookingRepo.Update(ctx, schema, bookingID, updateData); err != nil {
				if errors.Is(err, exceptions.ErrBookingNotFound) {
					return err
				}
				u.logPro.Error("Update booking status failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		ttl := max(time.Until(booking.CheckOut), 0)
		for _, bookingRoom := range reservedBookingRooms {
			bookingRoomData := data.BookingRoomData{
				ID:        bookingRoom.ID,
				ExpiredAt: booking.CheckOut,
			}
			bytes, err := json.Marshal(bookingRoomData)
			if err != nil {
				u.logPro.Error("Json marshal booking room data failed", port.String("schema", schema), port.Error(err))
				return err
			}

			cacheKey := fmt.Sprintf("%s:booking_room:%s", schema, bookingRoom.Code)
			if err = u.cachePro.SetObject(ctx, cacheKey, bytes, ttl); err != nil {
				u.logPro.Error("Save booking room data failed", port.String("schema", schema), port.String("cache_key", cacheKey), port.Error(err))
				return err
			}
		}

		return nil
	}

	if req.Status == model.BookingStatusCancelled {
		if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
			if err = u.bookingRoomRepo.DeleteAllByBookingID(ctx, schema, bookingID); err != nil {
				if errors.Is(err, exceptions.ErrProtectedRecord) {
					return err
				}
				u.logPro.Error("Delete booking rooms by booking id failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
				return err
			}

			if err = u.bookingRepo.Update(ctx, schema, bookingID, updateData); err != nil {
				if errors.Is(err, exceptions.ErrBookingNotFound) {
					return err
				}
				u.logPro.Error("Update booking status failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
				return err
			}

			return nil
		}); err != nil {
			return err
		}

		return nil
	}

	if err = u.bookingRepo.Update(ctx, schema, bookingID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrBookingNotFound) {
			return err
		}
		u.logPro.Error("Update booking status failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
		return err
	}

	return nil
}

func (u *bookingUseCaseImpl) Delete(ctx context.Context, schema string, bookingID, userID int64) error {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	booking, err := u.bookingRepo.FindByIDWithBookingRooms(ctx, schema, bookingID)
	if err != nil {
		u.logPro.Error("Find booking by id failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
		return err
	}
	if booking == nil {
		return exceptions.ErrBookingNotFound
	}
	if err := u.ensureBranchPtrInScope(scope, booking.BranchID); err != nil {
		return err
	}

	if booking.Status != model.BookingStatusConfirmed {
		return exceptions.ErrInvalidStatus
	}

	if err = u.txPro.WithinTx(ctx, func(ctx context.Context) error {
		if err = u.bookingRoomRepo.DeleteAllByBookingID(ctx, schema, bookingID); err != nil {
			if errors.Is(err, exceptions.ErrProtectedRecord) {
				return err
			}
			u.logPro.Error("Delete booking rooms by booking id failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
			return err
		}

		if err = u.bookingRepo.Delete(ctx, schema, bookingID); err != nil {
			if errors.Is(err, exceptions.ErrBookingNotFound) {
				return err
			}
			if errors.Is(err, exceptions.ErrProtectedRecord) {
				return err
			}
			u.logPro.Error("Delete booking failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (u *bookingUseCaseImpl) GetByID(ctx context.Context, schema string, bookingID, userID int64) (*response.BookingResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}

	booking, err := u.bookingRepo.FindByIDWithDetails(ctx, schema, bookingID)
	if err != nil {
		u.logPro.Error("Find booking by id with details failed", port.String("schema", schema), port.Int64("booking_id", bookingID), port.Error(err))
		return nil, err
	}
	if booking == nil {
		return nil, exceptions.ErrBookingNotFound
	}
	if err := u.ensureBranchPtrInScope(scope, booking.BranchID); err != nil {
		return nil, err
	}

	return mapper.ToBookingResponse(booking), nil
}

func (u *bookingUseCaseImpl) StatusSummary(ctx context.Context, schema string, userID int64) (*response.BookingStatusSummaryResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, err
	}

	var summary *model.BookingStatusSummary
	if scope.IsGlobal() {
		summary, err = u.bookingRepo.CountStatusSummary(ctx, schema)
	} else {
		summary, err = u.bookingRepo.CountStatusSummary(ctx, schema, scope.BranchIDs)
	}
	if err != nil {
		u.logPro.Error("Count booking status summary failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToBookingStatusSummaryResponse(summary), nil
}

func (u *bookingUseCaseImpl) Timeline(ctx context.Context, schema string, userID int64, query request.BookingTimelineQuery) (*response.BookingTimelineResponse, error) {
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

	query.FromDate = startOfDay(query.FromDate)

	if query.ToDate.IsZero() {
		query.ToDate = query.FromDate
	}
	query.ToDate = endOfDay(query.ToDate)

	minEndDate := endOfDay(query.FromDate.AddDate(0, 0, 5))
	if query.ToDate.Before(minEndDate) {
		query.ToDate = minEndDate
	}

	timelineQuery := mapper.ToBookingTimelineQuery(query)
	timelineQuery.BranchIDs = branchIDs
	timeline, err := u.bookingRepo.FindTimeline(ctx, schema, timelineQuery)
	if err != nil {
		u.logPro.Error("Find booking timeline failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	statusSummary, err := u.roomRepo.CountStatusSummary(ctx, schema, model.RoomStatusQuery{
		RoomTypeID: query.RoomTypeID,
		BranchID:   query.BranchID,
		BranchIDs:  branchIDs,
		SnapshotAt: time.Now(),
	})
	if err != nil {
		u.logPro.Error("Count room status summary failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	timeline.Status = statusSummary

	return mapper.ToBookingTimelineResponse(timeline), nil
}

func (u *bookingUseCaseImpl) List(ctx context.Context, schema string, userID int64, query request.BookingPaginationQuery) ([]*response.SimpleBookingResponse, *response.MetaResponse, error) {
	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return nil, nil, err
	}
	branchID, branchIDs, err := scope.Apply(query.BranchID)
	if err != nil {
		return nil, nil, err
	}
	query.BranchID = branchID

	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	bookingQuery := mapper.ToBookingQuery(query)
	bookingQuery.BranchIDs = branchIDs
	bookingRooms, total, err := u.bookingRepo.FindAllWithDetailsPaginated(ctx, schema, bookingQuery)
	if err != nil {
		u.logPro.Error("Find all booking paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToSimpleBookingsResponse(bookingRooms), meta, nil
}

func (u *bookingUseCaseImpl) validateBookingBranch(ctx context.Context, schema string, roomTypeID int64, branchID *int64) error {
	if branchID == nil {
		return nil
	}

	count, err := u.branchRepo.CountByIDs(ctx, schema, []int64{*branchID})
	if err != nil {
		u.logPro.Error("Count branch by id failed", port.String("schema", schema), port.Int64("branch_id", *branchID), port.Error(err))
		return err
	}
	if count == 0 {
		return exceptions.ErrBranchNotFound
	}

	exists, err := u.roomTypeRepo.ExistsByIDAndBranchID(ctx, schema, roomTypeID, *branchID)
	if err != nil {
		u.logPro.Error("Check room type in branch failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Int64("branch_id", *branchID), port.Error(err))
		return err
	}
	if !exists {
		return exceptions.ErrRoomTypeNotInBranch
	}

	return nil
}

func (u *bookingUseCaseImpl) ensureBranchPtrInScope(scope branchscope.Scope, branchID *int64) error {
	if scope.IsGlobal() {
		return nil
	}
	if branchID == nil {
		return exceptions.ErrForbidden
	}

	return scope.Ensure(*branchID)
}

func (u *bookingUseCaseImpl) validateBookingRoomsInBranch(booking *model.Booking, branchID *int64) error {
	if booking == nil || branchID == nil {
		return nil
	}

	for _, bookingRoom := range booking.BookingRooms {
		if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil {
			return exceptions.ErrRoomNotFound
		}
		if bookingRoom.Room.Floor.BranchID != *branchID {
			return exceptions.ErrRoomNotInBookingBranch
		}
	}

	return nil
}

func (u *bookingUseCaseImpl) calculateBookingTotalNetPrice(ctx context.Context, schema string, roomTypeID int64, checkIn, checkOut time.Time, roomNumber uint32) (float64, error) {
	roomType, err := u.roomTypeRepo.FindByID(ctx, schema, roomTypeID)
	if err != nil {
		u.logPro.Error("Find room type by id failed", port.String("schema", schema), port.Int64("room_type_id", roomTypeID), port.Error(err))
		return 0, err
	}
	if roomType == nil {
		return 0, exceptions.ErrRoomTypeNotFound
	}

	stayDays := utils.CalculateStayDays(checkIn, checkOut)
	checkInDate := utils.NormalizeDateOnly(checkIn)
	checkOutDate := checkInDate.AddDate(0, 0, int(stayDays))

	dailyTotalPrice, coveredDays, err := u.dailyRateRepo.SumFinalPriceByRoomTypeIDAndDateRange(ctx, schema, roomTypeID, checkInDate, checkOutDate)
	if err != nil {
		u.logPro.Error(
			"sum daily room type rate by date range failed",
			port.String("schema", schema),
			port.Int64("room_type_id", roomTypeID),
			port.Time("check_in", checkInDate),
			port.Time("check_out", checkOutDate),
			port.Error(err),
		)
		return 0, err
	}

	if coveredDays < 0 {
		coveredDays = 0
	}
	if coveredDays > stayDays {
		coveredDays = stayDays
	}

	missingDays := stayDays - coveredDays
	pricePerRoom := dailyTotalPrice + (float64(missingDays) * roomType.Price)

	roomCount := float64(roomNumber)
	if roomCount < 1 {
		roomCount = 1
	}

	return pricePerRoom * roomCount, nil
}

func (u *bookingUseCaseImpl) generateBookingNumber() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, 17)
	charsetLen := byte(len(charset))

	randomBytes := make([]byte, 17)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	for i := range 17 {
		result[i] = charset[randomBytes[i]%charsetLen]
	}

	return fmt.Sprintf("OFF%s", string(result)), nil
}

func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func endOfDay(t time.Time) time.Time {
	dayStart := startOfDay(t)
	return dayStart.Add(24*time.Hour - time.Second)
}

func (u *bookingUseCaseImpl) predictBooking(schema string, bookingID int64) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		booking, err := u.bookingRepo.FindByIDWithDetails(ctx, schema, bookingID)
		if err != nil {
			u.logPro.Error(
				"find booking by id for prediction failed",
				port.String("schema", schema),
				port.Int64("booking_id", bookingID),
				port.Error(err),
			)
			return
		}
		if booking == nil {
			return
		}

		predictRes, err := u.aiPro.PredictBooking(ctx, booking)
		if err != nil {
			u.logPro.Error("Predict booking from gRPC failed", port.String("schema", schema), port.Error(err))
			return
		}

		predict := model.BookingPredict{
			CancellationProbability: predictRes.CancellationProbability,
			WillCancel:              predictRes.WillCancel,
			RiskLevel:               predictRes.RiskLevel,
		}

		predictJSON, err := json.Marshal(predict)
		if err != nil {
			u.logPro.Error("Marshal booking prediction failed", port.String("schema", schema), port.Error(err))
			return
		}

		if err := u.bookingRepo.Update(ctx, schema, bookingID, map[string]any{
			"predict": datatypes.JSON(predictJSON),
		}); err != nil {
			u.logPro.Error("Update booking prediction failed", port.String("schema", schema), port.Error(err))
			return
		}
	}()
}
