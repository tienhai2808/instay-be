package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToSupportCategoryResponse(supCat *model.SupportCategory) *response.SupportCategoryResponse {
	if supCat == nil {
		return nil
	}

	return &response.SupportCategoryResponse{
		ID:                  supCat.ID,
		Name:                supCat.Name,
		CreatedAt:           supCat.CreatedAt,
		UpdatedAt:           supCat.UpdatedAt,
		CreatedBy:           ToBasicUserResponse(supCat.CreatedBy),
		UpdatedBy:           ToBasicUserResponse(supCat.UpdatedBy),
		Department:          ToBasicDepartmentResponse(supCat.Department),
		SupportTicketsCount: supCat.SupportTicketsCount,
	}
}

func ToSupportCategoriesResponse(supCats []*model.SupportCategory) []*response.SupportCategoryResponse {
	if len(supCats) == 0 {
		return make([]*response.SupportCategoryResponse, 0)
	}

	supCatsRes := make([]*response.SupportCategoryResponse, 0, len(supCats))
	for _, supCat := range supCats {
		supCatsRes = append(supCatsRes, ToSupportCategoryResponse(supCat))
	}

	return supCatsRes
}

func ToBasicSupportCategoryResponse(supCat *model.SupportCategory) *response.BasicSupportCategoryResponse {
	if supCat == nil {
		return nil
	}

	return &response.BasicSupportCategoryResponse{
		ID:   supCat.ID,
		Name: supCat.Name,
	}
}

func ToBasicSupportCategoriesResponse(supCats []*model.SupportCategory) []*response.BasicSupportCategoryResponse {
	if len(supCats) == 0 {
		return make([]*response.BasicSupportCategoryResponse, 0)
	}

	supCatsRes := make([]*response.BasicSupportCategoryResponse, 0, len(supCats))
	for _, supCat := range supCats {
		supCatsRes = append(supCatsRes, ToBasicSupportCategoryResponse(supCat))
	}

	return supCatsRes
}

func ToSupportTicketPrioritiesResponse() []response.SupportTicketPriorityResponse {
	return []response.SupportTicketPriorityResponse{
		{
			Name:          model.SupportTicketPriorityLow,
			DisplayNameVI: model.SupportTicketPriorityLowDisplayNameVI,
			DisplayNameEN: model.SupportTicketPriorityLowDisplayNameEN,
		},
		{
			Name:          model.SupportTicketPriorityMedium,
			DisplayNameVI: model.SupportTicketPriorityMediumDisplayNameVI,
			DisplayNameEN: model.SupportTicketPriorityMediumDisplayNameEN,
		},
		{
			Name:          model.SupportTicketPriorityHigh,
			DisplayNameVI: model.SupportTicketPriorityHighDisplayNameVI,
			DisplayNameEN: model.SupportTicketPriorityHighDisplayNameEN,
		},
		{
			Name:          model.SupportTicketPriorityUrgent,
			DisplayNameVI: model.SupportTicketPriorityUrgentDisplayNameVI,
			DisplayNameEN: model.SupportTicketPriorityUrgentDisplayNameEN,
		},
	}
}

func ToSimpleSupportTicketResponse(supTicket *model.SupportTicket) *response.SimpleSupportTicketResponse {
	if supTicket == nil {
		return nil
	}

	return &response.SimpleSupportTicketResponse{
		ID:              supTicket.ID,
		SupportCategory: ToBasicSupportCategoryResponse(supTicket.SupportCategory),
		BookingRoom:     ToBasicBookingRoomResponse(supTicket.BookingRoom),
		Content:         supTicket.Content,
		Status:          supTicket.Status,
		Priority:        supTicket.Priority,
		CreatedAt:       supTicket.CreatedAt,
	}
}

func ToSimpleSupportTicketsResponse(supTickets []*model.SupportTicket) []*response.SimpleSupportTicketResponse {
	if len(supTickets) == 0 {
		return make([]*response.SimpleSupportTicketResponse, 0)
	}

	supTicketsRes := make([]*response.SimpleSupportTicketResponse, 0, len(supTickets))
	for _, supTicket := range supTickets {
		supTicketsRes = append(supTicketsRes, ToSimpleSupportTicketResponse(supTicket))
	}

	return supTicketsRes
}

func ToSupportTicketResponse(supTicket *model.SupportTicket) *response.SupportTicketResponse {
	if supTicket == nil {
		return nil
	}

	return &response.SupportTicketResponse{
		ID:                supTicket.ID,
		SupportCategory:   ToBasicSupportCategoryResponse(supTicket.SupportCategory),
		BookingRoom:       ToBasicBookingRoomResponse(supTicket.BookingRoom),
		Content:           supTicket.Content,
		Status:            supTicket.Status,
		Priority:          supTicket.Priority,
		ResolutionSeconds: supTicket.ResolutionSeconds,
		CreatedAt:         supTicket.CreatedAt,
		UpdatedAt:         supTicket.UpdatedAt,
		UpdatedBy:         ToBasicUserResponse(supTicket.UpdatedBy),
	}
}

func ToSupportTicketsResponse(supTickets []*model.SupportTicket) []*response.SupportTicketResponse {
	if len(supTickets) == 0 {
		return make([]*response.SupportTicketResponse, 0)
	}

	supTicketsRes := make([]*response.SupportTicketResponse, 0, len(supTickets))
	for _, supTicket := range supTickets {
		supTicketsRes = append(supTicketsRes, ToSupportTicketResponse(supTicket))
	}

	return supTicketsRes
}

func ToSupportTicketQuery(query request.SupportTicketPaginationQuery) model.SupportTicketQuery {
	return model.SupportTicketQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Sort:     query.Sort,
		Order:    query.Order,
		Status:   query.Status,
		Priority: query.Priority,
		FromDate: query.FromDate,
		ToDate:   query.ToDate,
	}
}

func ToSupportCategoryQuery(query request.SupportCategoryPaginationQuery) model.SupportCategoryQuery {
	return model.SupportCategoryQuery{
		Page:         query.Page,
		Limit:        query.Limit,
		DepartmentID: query.DepartmentID,
		Search:       query.Search,
	}
}
