package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBasicTenantResponse(tenant *model.Tenant) *response.BasicTenantResponse {
	if tenant == nil {
		return nil
	}

	return &response.BasicTenantResponse{
		Name:  tenant.Name,
		Logo:  tenant.Logo,
		Phone: tenant.Phone,
	}
}

func ToTenantSettingsResponse(tenant *model.Tenant) *response.TenantSettingsResponse {
	if tenant == nil {
		return nil
	}

	res := &response.TenantSettingsResponse{}
	res.IMAP.Host = tenant.Settings.IMAP.Host
	res.IMAP.Port = tenant.Settings.IMAP.Port
	res.IMAP.User = tenant.Settings.IMAP.User
	res.IMAP.HasPassword = tenant.Settings.IMAP.Password != ""

	res.EmailBooking.Subject = tenant.Settings.EmailBooking.Subject
	res.EmailBooking.BookingNumber = tenant.Settings.EmailBooking.BookingNumber
	res.EmailBooking.Guest = tenant.Settings.EmailBooking.Guest
	res.EmailBooking.CheckIn = tenant.Settings.EmailBooking.CheckIn
	res.EmailBooking.CheckOut = tenant.Settings.EmailBooking.CheckOut
	res.EmailBooking.BookedOn = tenant.Settings.EmailBooking.BookedOn
	res.EmailBooking.RoomsBooked = tenant.Settings.EmailBooking.RoomsBooked
	res.EmailBooking.Source = tenant.Settings.EmailBooking.Source
	res.EmailBooking.TotalNetPrice = tenant.Settings.EmailBooking.TotalNetPrice
	res.EmailBooking.TotalSellPrice = tenant.Settings.EmailBooking.TotalSellPrice
	res.EmailBooking.GuestNumber = tenant.Settings.EmailBooking.GuestNumber
	res.EmailBooking.RoomNumber = tenant.Settings.EmailBooking.RoomNumber
	res.EmailBooking.PromotionName = tenant.Settings.EmailBooking.PromotionName
	res.EmailBooking.MealPlan = tenant.Settings.EmailBooking.MealPlan
	res.EmailBooking.BookingConditions = tenant.Settings.EmailBooking.BookingConditions
	res.EmailBooking.BookingPreferences = tenant.Settings.EmailBooking.BookingPreferences

	res.EmailCancel.Subject = tenant.Settings.EmailCancel.Subject
	res.EmailCancel.Bodies = tenant.Settings.EmailCancel.Bodies

	return res
}
