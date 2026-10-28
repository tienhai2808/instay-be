package tenant

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/InstaySystem/is_v2-be/internal/application/dto/data"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/mapper"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/application/port"
	"github.com/InstaySystem/is_v2-be/internal/domain/exceptions"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/domain/repository"
)

type tenantUseCaseImpl struct {
	logPro     port.LogProvider
	mqPro      port.MessageQueueProvider
	tenantRepo repository.TenantRepository
}

func NewTenantUseCase(
	logPro port.LogProvider,
	mqPro port.MessageQueueProvider,
	tenantRepo repository.TenantRepository,
) TenantUseCase {
	return &tenantUseCaseImpl{
		logPro,
		mqPro,
		tenantRepo,
	}
}

func (u *tenantUseCaseImpl) GetBySchemaName(ctx context.Context, schema string) (*response.BasicTenantResponse, error) {
	tenant, err := u.tenantRepo.FindBySchemaName(ctx, schema)
	if err != nil {
		u.logPro.Error("Find tenant by schema failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if tenant == nil {
		return nil, exceptions.ErrTenantNotFound
	}

	return mapper.ToBasicTenantResponse(tenant), nil
}

func (u *tenantUseCaseImpl) GetSettingsBySchemaName(ctx context.Context, schema string) (*response.TenantSettingsResponse, error) {
	tenant, err := u.tenantRepo.FindBySchemaName(ctx, schema)
	if err != nil {
		u.logPro.Error("Find tenant by schema failed", port.String("schema", schema), port.Error(err))
		return nil, err
	}
	if tenant == nil {
		return nil, exceptions.ErrTenantNotFound
	}

	return mapper.ToTenantSettingsResponse(tenant), nil
}

func (u *tenantUseCaseImpl) UpdateSettingsBySchemaName(ctx context.Context, schema string, req request.UpdateTenantSettingsRequest) error {
	tenant, err := u.tenantRepo.FindBySchemaName(ctx, schema)
	if err != nil {
		u.logPro.Error("Find tenant by schema failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if tenant == nil {
		return exceptions.ErrTenantNotFound
	}

	settings, err := validateAndBuildTenantSettings(req, tenant.Settings)
	if err != nil {
		return err
	}

	settingsJSON, err := json.Marshal(settings)
	if err != nil {
		u.logPro.Error("Marshal tenant settings failed", port.String("schema", schema), port.Error(err))
		return err
	}

	if err := u.tenantRepo.UpdateBySchemaName(ctx, schema, map[string]any{
		"settings": settingsJSON,
	}); err != nil {
		if errors.Is(err, exceptions.ErrTenantNotFound) {
			return err
		}
		u.logPro.Error("Update tenant settings by schema failed", port.String("schema", schema), port.Error(err))
		return err
	}

	return nil
}

func (u *tenantUseCaseImpl) UpdateBySchemaName(ctx context.Context, schema string, req request.UpdateTenantRequest) error {
	tenant, err := u.tenantRepo.FindBySchemaName(ctx, schema)
	if err != nil {
		u.logPro.Error("Find tenant by schema failed", port.String("schema", schema), port.Error(err))
		return err
	}
	if tenant == nil {
		return exceptions.ErrTenantNotFound
	}

	updateData := map[string]any{
		"name":  req.Name,
		"logo":  req.Logo,
		"phone": req.Phone,
		"city":  req.City,
	}

	if err := u.tenantRepo.UpdateBySchemaName(ctx, schema, updateData); err != nil {
		if errors.Is(err, exceptions.ErrTenantNotFound) {
			return err
		}
		u.logPro.Error("Update tenant by schema failed", port.String("schema", schema), port.Error(err))
		return err
	}

	oldLogo := tenant.Logo
	newLogo := req.Logo
	shouldDelete := false
	if oldLogo != nil {
		if newLogo == nil {
			shouldDelete = true
		} else if *oldLogo != *newLogo {
			shouldDelete = true
		}
	}

	if shouldDelete {
		deleteFilesMsg := data.DeleteFileMessage{
			Schema: schema,
			Keys:   []string{*oldLogo},
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

	return nil
}

func validateAndBuildTenantSettings(
	req request.UpdateTenantSettingsRequest,
	existing model.TenantSettings,
) (model.TenantSettings, error) {
	var settings model.TenantSettings

	host := strings.TrimSpace(req.IMAP.Host)
	user := strings.TrimSpace(req.IMAP.User)
	password := strings.TrimSpace(req.IMAP.Password)
	if password == "" {
		password = strings.TrimSpace(existing.IMAP.Password)
	}
	portValue := req.IMAP.Port

	hasAnyIMAP := host != "" || user != "" || password != "" || portValue != 0
	hasAllIMAP := host != "" && user != "" && password != "" && portValue != 0

	if hasAnyIMAP {
		switch {
		case host == "":
			return settings, badTenantSettingsField("imap.host", "required_with", "imap")
		case portValue == 0:
			return settings, badTenantSettingsField("imap.port", "required_with", "imap")
		case portValue < 0:
			return settings, badTenantSettingsField("imap.port", "gt", "0")
		case user == "":
			return settings, badTenantSettingsField("imap.user", "required_with", "imap")
		case password == "":
			return settings, badTenantSettingsField("imap.password", "required_with", "imap")
		}
	} else if portValue < 0 {
		return settings, badTenantSettingsField("imap.port", "gt", "0")
	}

	emailFields := map[string]model.PatternList{
		"email_booking.subject":             req.EmailBooking.Subject,
		"email_booking.booking_number":      req.EmailBooking.BookingNumber,
		"email_booking.guest":               req.EmailBooking.Guest,
		"email_booking.check_in":            req.EmailBooking.CheckIn,
		"email_booking.check_out":           req.EmailBooking.CheckOut,
		"email_booking.booked_on":           req.EmailBooking.BookedOn,
		"email_booking.rooms_booked":        req.EmailBooking.RoomsBooked,
		"email_booking.source":              req.EmailBooking.Source,
		"email_booking.total_net_price":     req.EmailBooking.TotalNetPrice,
		"email_booking.total_sell_price":    req.EmailBooking.TotalSellPrice,
		"email_booking.guest_number":        req.EmailBooking.GuestNumber,
		"email_booking.room_number":         req.EmailBooking.RoomNumber,
		"email_booking.promotion_name":      req.EmailBooking.PromotionName,
		"email_booking.meal_plan":           req.EmailBooking.MealPlan,
		"email_booking.booking_conditions":  req.EmailBooking.BookingConditions,
		"email_booking.booking_preferences": req.EmailBooking.BookingPreferences,
		"email_cancel.subject":              req.EmailCancel.Subject,
		"email_cancel.bodies":               req.EmailCancel.Bodies,
	}

	hasAnyEmailValue := false
	for field, patterns := range emailFields {
		if patterns == nil {
			return settings, badTenantSettingsField(field, "required", "")
		}
		if len(patterns) > 0 {
			hasAnyEmailValue = true
		}
	}

	if hasAnyEmailValue && !hasAllIMAP {
		return settings, badTenantSettingsField("imap", "required_with", "email_settings")
	}

	settings.IMAP.Host = host
	settings.IMAP.Port = portValue
	settings.IMAP.User = user
	settings.IMAP.Password = password

	settings.EmailBooking.Subject = req.EmailBooking.Subject
	settings.EmailBooking.BookingNumber = req.EmailBooking.BookingNumber
	settings.EmailBooking.Guest = req.EmailBooking.Guest
	settings.EmailBooking.CheckIn = req.EmailBooking.CheckIn
	settings.EmailBooking.CheckOut = req.EmailBooking.CheckOut
	settings.EmailBooking.BookedOn = req.EmailBooking.BookedOn
	settings.EmailBooking.RoomsBooked = req.EmailBooking.RoomsBooked
	settings.EmailBooking.Source = req.EmailBooking.Source
	settings.EmailBooking.TotalNetPrice = req.EmailBooking.TotalNetPrice
	settings.EmailBooking.TotalSellPrice = req.EmailBooking.TotalSellPrice
	settings.EmailBooking.GuestNumber = req.EmailBooking.GuestNumber
	settings.EmailBooking.RoomNumber = req.EmailBooking.RoomNumber
	settings.EmailBooking.PromotionName = req.EmailBooking.PromotionName
	settings.EmailBooking.MealPlan = req.EmailBooking.MealPlan
	settings.EmailBooking.BookingConditions = req.EmailBooking.BookingConditions
	settings.EmailBooking.BookingPreferences = req.EmailBooking.BookingPreferences

	settings.EmailCancel.Subject = req.EmailCancel.Subject
	settings.EmailCancel.Bodies = req.EmailCancel.Bodies

	return settings, nil
}

func badTenantSettingsField(field, tag, param string) error {
	return exceptions.ErrBadRequest.WithData(map[string]string{
		"field": field,
		"tag":   tag,
		"param": param,
	})
}
