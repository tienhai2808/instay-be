package guest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

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

type guestUseCaseImpl struct {
	logPro          port.LogProvider
	idPro           port.IDProvider
	kycPro          port.KYCProvider
	guestRepo       repository.GuestRepository
	bookingRoomRepo repository.BookingRoomRepository
	userRepo        repository.UserRepository
}

type kycVendorData struct {
	Schema        string `json:"schema"`
	BookingRoomID int64  `json:"booking_room_id"`
}

type verifiedGuestIdentity struct {
	FullName       string
	PassportNumber string
	National       string
	DOB            time.Time
}

func NewGuestUseCase(
	logPro port.LogProvider,
	idPro port.IDProvider,
	kycPro port.KYCProvider,
	guestRepo repository.GuestRepository,
	bookingRoomRepo repository.BookingRoomRepository,
	userRepo repository.UserRepository,
) GuestUseCase {
	return &guestUseCaseImpl{
		logPro,
		idPro,
		kycPro,
		guestRepo,
		bookingRoomRepo,
		userRepo,
	}
}

func (u *guestUseCaseImpl) Create(ctx context.Context, schema string, userID int64, req request.CreateGuestRequest) error {
	bookingRoom, err := u.bookingRoomRepo.FindByIDWithBookingAndRoom(ctx, schema, req.BookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id with booking failed", port.String("schema", schema), port.Int64("booking_room_id", req.BookingRoomID), port.Error(err))
		return err
	}
	if bookingRoom == nil || bookingRoom.Booking == nil {
		return exceptions.ErrBookingRoomNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return err
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate guest id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	guest := model.NewGuest(
		id,
		req.BookingRoomID,
		req.FullName,
		req.National,
		req.PassportNumber,
		req.DOB,
		bookingRoom.Booking.CheckIn,
		bookingRoom.Booking.CheckOut,
		&userID,
		&userID,
	)

	if err = u.guestRepo.Create(ctx, schema, guest); err != nil {
		if errors.Is(err, exceptions.ErrBookingRoomNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrGuestInBookingRoomAlreadyExists) {
			return err
		}
		u.logPro.Error("Create guest failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *guestUseCaseImpl) CreateVerification(ctx context.Context, schema string, bookingRoomID int64, callbackURL string) (string, string, error) {
	vendorDataBytes, err := json.Marshal(kycVendorData{
		Schema:        schema,
		BookingRoomID: bookingRoomID,
	})
	if err != nil {
		u.logPro.Error("Marshal kyc vendor data failed", port.String("schema", schema), port.Error(err))
		return "", "", err
	}

	session, err := u.kycPro.CreateVerificationSession(ctx, string(vendorDataBytes), callbackURL)
	if err != nil {
		u.logPro.Error("Create verification session failed", port.String("schema", schema), port.Error(err))
		return "", "", err
	}

	return session.SessionID, session.VerificationURL, nil
}

func (u *guestUseCaseImpl) Verify(ctx context.Context, req request.VerifyGuestRequest) error {
	if !strings.EqualFold(strings.TrimSpace(req.Status), "approved") {
		return nil
	}

	schema, bookingRoomID, err := parseWebhookVendorData(req)
	if err != nil {
		return err
	}

	identity, err := parseVerifiedGuestIdentity(req.Decision)
	if err != nil {
		return err
	}

	bookingRoom, err := u.bookingRoomRepo.FindByIDWithBooking(ctx, schema, bookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id with booking failed", port.String("schema", schema), port.Int64("booking_room_id", bookingRoomID), port.Error(err))
		return err
	}
	if bookingRoom == nil || bookingRoom.Booking == nil {
		return fmt.Errorf("booking room or booking not found")
	}
	if bookingRoom.Booking.Status != model.BookingStatusCheckedIn {
		return fmt.Errorf("booking status must be checked_in")
	}

	id, err := u.idPro.NextID()
	if err != nil {
		u.logPro.Error("Generate guest id failed", port.String("schema", schema), port.Error(err))
		return err
	}

	guest := model.NewGuest(
		id,
		bookingRoomID,
		identity.FullName,
		identity.National,
		identity.PassportNumber,
		identity.DOB,
		bookingRoom.Booking.CheckIn,
		bookingRoom.Booking.CheckOut,
		nil,
		nil,
	)

	if err = u.guestRepo.Create(ctx, schema, guest); err != nil {
		if errors.Is(err, exceptions.ErrBookingRoomNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrGuestInBookingRoomAlreadyExists) {
			return err
		}
		u.logPro.Error("Create guest from kyc webhook failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *guestUseCaseImpl) Update(ctx context.Context, schema string, guestID, userID int64, req request.UpdateGuestRequest) error {
	guest, err := u.guestRepo.FindByIDWithBookingRoomDetails(ctx, schema, guestID)
	if err != nil {
		u.logPro.Error("Find guest by id failed", port.String("schema", schema), port.Int64("guest_id", guestID), port.Error(err))
		return err
	}
	if guest == nil {
		return exceptions.ErrGuestNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, guest.BookingRoom); err != nil {
		return err
	}

	bookingRoom, err := u.bookingRoomRepo.FindByIDWithRoom(ctx, schema, req.BookingRoomID)
	if err != nil {
		u.logPro.Error("Find booking room by id with room failed", port.String("schema", schema), port.Int64("booking_room_id", req.BookingRoomID), port.Error(err))
		return err
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, bookingRoom); err != nil {
		return err
	}

	updateData := map[string]any{
		"full_name":       req.FullName,
		"dob":             req.DOB,
		"check_in":        req.CheckIn,
		"check_out":       req.CheckOut,
		"passport_number": req.PassportNumber,
		"national":        req.National,
		"booking_room_id": req.BookingRoomID,
		"updated_by_id":   userID,
	}

	if err := u.guestRepo.Update(ctx, schema, guestID, updateData); err != nil {
		if errors.Is(err, exceptions.ErrGuestNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrBookingRoomNotFound) {
			return err
		}
		u.logPro.Error("Updated guest failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *guestUseCaseImpl) Delete(ctx context.Context, schema string, guestID, userID int64) error {
	guest, err := u.guestRepo.FindByIDWithBookingRoomDetails(ctx, schema, guestID)
	if err != nil {
		u.logPro.Error("Find guest by id failed", port.String("schema", schema), port.Int64("guest_id", guestID), port.Error(err))
		return err
	}
	if guest == nil {
		return exceptions.ErrGuestNotFound
	}
	if err := u.ensureBookingRoomInUserScope(ctx, schema, userID, guest.BookingRoom); err != nil {
		return err
	}

	if err := u.guestRepo.Delete(ctx, schema, guestID); err != nil {
		if errors.Is(err, exceptions.ErrGuestNotFound) {
			return err
		}
		if errors.Is(err, exceptions.ErrProtectedRecord) {
			return err
		}
		u.logPro.Error("Delete guest failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *guestUseCaseImpl) List(ctx context.Context, schema string, userID int64, query request.GuestPaginationQuery) ([]*response.GuestResponse, *response.MetaResponse, error) {
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 10
	}

	guestQuery, err := u.scopedGuestQuery(ctx, schema, userID, query)
	if err != nil {
		return nil, nil, err
	}

	guests, total, err := u.guestRepo.FindAllWithDetailsPaginated(ctx, schema, guestQuery)
	if err != nil {
		u.logPro.Error("Find all guest paginated failed", port.String("schema", schema), port.Error(err))
		return nil, nil, err
	}

	meta := mapper.ToMetaResponse(total, query.Page, query.Limit)

	return mapper.ToGuestsResponse(guests), meta, nil
}

func (u *guestUseCaseImpl) Export(ctx context.Context, schema string, userID int64, req request.ExportListGuestsRequest) ([]*response.GuestResponse, error) {
	query, err := toGuestExportQuery(req)
	if err != nil {
		return nil, err
	}
	branchIDs, err := u.guestBranchIDs(ctx, schema, userID)
	if err != nil {
		return nil, err
	}
	query.BranchIDs = branchIDs

	guests, err := u.guestRepo.FindAllWithDetails(ctx, schema, query)
	if err != nil {
		u.logPro.Error("Find all guest for export failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}

	return mapper.ToGuestsResponse(guests), nil
}

func (u *guestUseCaseImpl) scopedGuestQuery(ctx context.Context, schema string, userID int64, query request.GuestPaginationQuery) (model.GuestQuery, error) {
	branchIDs, err := u.guestBranchIDs(ctx, schema, userID)
	if err != nil {
		return model.GuestQuery{}, err
	}

	guestQuery := mapper.ToGuestQuery(query)
	guestQuery.BranchIDs = branchIDs

	return guestQuery, nil
}

func (u *guestUseCaseImpl) guestBranchIDs(ctx context.Context, schema string, userID int64) ([]int64, error) {
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

func (u *guestUseCaseImpl) ensureBookingRoomInUserScope(ctx context.Context, schema string, userID int64, bookingRoom *model.BookingRoom) error {
	if bookingRoom == nil || bookingRoom.Room == nil || bookingRoom.Room.Floor == nil {
		return exceptions.ErrBookingRoomNotFound
	}

	scope, err := branchscope.Load(ctx, schema, userID, u.userRepo)
	if err != nil {
		u.logPro.Error("Load user branch scope failed", port.String("schema", schema), port.Int64("user_id", userID), port.Error(err))
		return err
	}

	return scope.Ensure(bookingRoom.Room.Floor.BranchID)
}

func parseWebhookVendorData(req request.VerifyGuestRequest) (string, int64, error) {
	vendorData := strings.TrimSpace(req.VendorData)
	if vendorData == "" && req.Decision != nil {
		vendorData = strings.TrimSpace(req.Decision.VendorData)
	}
	if vendorData == "" {
		return "", 0, fmt.Errorf("vendor_data is required")
	}

	schema, bookingRoomID, ok := parseVendorDataPayload(vendorData)
	if ok && schema != "" && bookingRoomID > 0 {
		return schema, bookingRoomID, nil
	}

	schema, bookingRoomID, ok = parseVendorDataPipe(vendorData)
	if ok && schema != "" && bookingRoomID > 0 {
		return schema, bookingRoomID, nil
	}

	return "", 0, fmt.Errorf("invalid vendor_data")
}

func parseVendorDataPayload(vendorData string) (string, int64, bool) {
	var payload map[string]any
	if err := json.Unmarshal([]byte(vendorData), &payload); err != nil {
		return "", 0, false
	}

	schema := strings.TrimSpace(utils.ParseString(payload["schema"]))
	bookingRoomID := utils.ParseInt64(payload["booking_room_id"])

	if schema == "" || bookingRoomID <= 0 {
		return "", 0, false
	}

	return schema, bookingRoomID, true
}

func parseVendorDataPipe(vendorData string) (string, int64, bool) {
	parts := strings.Split(vendorData, "|")
	if len(parts) != 2 {
		return "", 0, false
	}

	schema := strings.TrimSpace(parts[0])
	bookingRoomID, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64)
	if err != nil || schema == "" || bookingRoomID <= 0 {
		return "", 0, false
	}

	return schema, bookingRoomID, true
}

func parseVerifiedGuestIdentity(decision *request.VerifyGuestDecisionRequest) (*verifiedGuestIdentity, error) {
	if decision == nil || len(decision.IDVerifications) == 0 {
		return nil, fmt.Errorf("decision.id_verifications is required")
	}

	var idVerification *request.VerifyGuestIDVerificationRequest
	for i := range decision.IDVerifications {
		if strings.EqualFold(strings.TrimSpace(decision.IDVerifications[i].Status), "approved") {
			idVerification = &decision.IDVerifications[i]
			break
		}
	}
	if idVerification == nil {
		idVerification = &decision.IDVerifications[0]
	}

	fullName := strings.TrimSpace(idVerification.FullName)
	if fullName == "" {
		fullName = strings.TrimSpace(strings.TrimSpace(idVerification.FirstName) + " " + strings.TrimSpace(idVerification.LastName))
	}
	if fullName == "" {
		return nil, fmt.Errorf("full_name is required")
	}

	normalizedDocumentType := strings.ToLower(strings.NewReplacer(" ", "", "_", "", "-", "").Replace(strings.TrimSpace(idVerification.DocumentType)))
	preferPersonalNumber := normalizedDocumentType == "identitycard" || normalizedDocumentType == "identitynumber"
	passportNumber := strings.TrimSpace(idVerification.DocumentNumber)
	numberField := "document_number"

	personalNumber := strings.TrimSpace(idVerification.PersonalNumber)
	if preferPersonalNumber && personalNumber != "" {
		passportNumber = personalNumber
		numberField = "personal_number"
	}

	if passportNumber == "" || len(passportNumber) > 20 {
		return nil, fmt.Errorf("%s is invalid", numberField)
	}

	national := strings.TrimSpace(idVerification.Nationality)
	if national == "" || len(national) > 10 {
		return nil, fmt.Errorf("nationality is invalid")
	}

	dob, err := parseDOB(idVerification.DateOfBirth)
	if err != nil {
		return nil, err
	}

	return &verifiedGuestIdentity{
		FullName:       fullName,
		PassportNumber: passportNumber,
		National:       national,
		DOB:            dob,
	}, nil
}

func toGuestExportQuery(req request.ExportListGuestsRequest) (model.GuestQuery, error) {
	query := model.GuestQuery{
		Sort:  "check_in",
		Order: "desc",
	}

	if req.Sort != nil && strings.TrimSpace(*req.Sort) != "" {
		query.Sort = strings.TrimSpace(*req.Sort)
	}

	sortOrder := ""
	if req.Order != nil && strings.TrimSpace(*req.Order) != "" {
		sortOrder = strings.TrimSpace(*req.Order)
	}
	if sortOrder != "" {
		query.Order = sortOrder
	}

	fromProvided := req.FromDate != nil && !req.FromDate.IsZero()
	toProvided := req.ToDate != nil && !req.ToDate.IsZero()

	switch {
	case fromProvided && !toProvided:
		return query, exceptions.ErrBadRequest.WithData(map[string]string{
			"field": "to_date",
			"tag":   "required_with",
			"param": "from_date",
		})
	case !fromProvided && toProvided:
		return query, exceptions.ErrBadRequest.WithData(map[string]string{
			"field": "from_date",
			"tag":   "required_with",
			"param": "to_date",
		})
	case !fromProvided && !toProvided:
		return query, nil
	}

	fromDate := utils.NormalizeDateOnly(*req.FromDate)
	toDate := utils.NormalizeDateOnly(*req.ToDate)

	if toDate.Before(fromDate) {
		return query, exceptions.ErrBadRequest.WithData(map[string]string{
			"field": "to_date",
			"tag":   "gtefield",
			"param": "from_date",
		})
	}

	query.FromDate = fromDate
	query.ToDate = toDate

	return query, nil
}

func parseDOB(value string) (time.Time, error) {
	dateOfBirth := strings.TrimSpace(value)
	if dateOfBirth == "" {
		return time.Time{}, fmt.Errorf("date_of_birth is required")
	}

	layouts := []string{
		"2006-01-02",
		time.RFC3339,
		"2006-01-02T15:04:05.999999Z07:00",
	}

	for _, layout := range layouts {
		dob, err := time.Parse(layout, dateOfBirth)
		if err == nil {
			return dob, nil
		}
	}

	return time.Time{}, fmt.Errorf("date_of_birth is invalid")
}
