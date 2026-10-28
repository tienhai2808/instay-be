package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToGuestResponse(guest *model.Guest) *response.GuestResponse {
	if guest == nil {
		return nil
	}

	return &response.GuestResponse{
		ID:             guest.ID,
		FullName:       guest.FullName,
		DOB:            guest.DOB,
		CheckIn:        guest.CheckIn,
		CheckOut:       guest.CheckOut,
		PassportNumber: guest.PassportNumber,
		National:       guest.National,
		BookingRoom:    ToBasicBookingRoomResponse(guest.BookingRoom),
		CreatedAt:      guest.CreatedAt,
		UpdatedAt:      guest.UpdatedAt,
		CreatedBy:      ToBasicUserResponse(guest.CreatedBy),
		UpdatedBy:      ToBasicUserResponse(guest.UpdatedBy),
	}
}

func ToGuestQuery(query request.GuestPaginationQuery) model.GuestQuery {
	return model.GuestQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Sort:     query.Sort,
		Order:    query.Order,
		Search:   query.Search,
		FromDate: query.FromDate,
		ToDate:   query.ToDate,
	}
}

func ToGuestsResponse(guests []*model.Guest) []*response.GuestResponse {
	if len(guests) == 0 {
		return make([]*response.GuestResponse, 0)
	}

	guestsRes := make([]*response.GuestResponse, 0, len(guests))
	for _, guest := range guests {
		guestsRes = append(guestsRes, ToGuestResponse(guest))
	}

	return guestsRes
}
