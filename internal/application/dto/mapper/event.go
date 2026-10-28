package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToEventDetailsResponse(event *model.Event) *response.EventDetailsResponse {
	if event == nil {
		return nil
	}

	return &response.EventDetailsResponse{
		ID:          event.ID,
		Title:       event.Title,
		Excerpt:     event.Excerpt,
		Content:     event.Content,
		Thumbnail:   event.Thumbnail,
		IsActive:    event.IsActive,
		PublishedAt: event.PublishedAt,
		Branches:    ToBasicBranchesResponse(event.Branches),
		CreatedAt:   event.CreatedAt,
		UpdatedAt:   event.UpdatedAt,
		CreatedBy:   ToBasicUserResponse(event.CreatedBy),
		UpdatedBy:   ToBasicUserResponse(event.UpdatedBy),
	}
}

func ToEventResponse(event *model.Event) *response.EventResponse {
	if event == nil {
		return nil
	}

	return &response.EventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Content:     event.Content,
		Thumbnail:   event.Thumbnail,
		PublishedAt: event.PublishedAt,
		Branches:    ToBasicBranchesResponse(event.Branches),
	}
}

func ToBasicEventResponse(event *model.Event) *response.BasicEventResponse {
	if event == nil {
		return nil
	}

	return &response.BasicEventResponse{
		ID:          event.ID,
		Slug:        event.Slug,
		Title:       event.Title,
		Excerpt:     event.Excerpt,
		Thumbnail:   event.Thumbnail,
		PublishedAt: event.PublishedAt,
	}
}

func ToBasicEventsResponse(events []*model.Event) []*response.BasicEventResponse {
	if len(events) == 0 {
		return make([]*response.BasicEventResponse, 0)
	}

	eventsRes := make([]*response.BasicEventResponse, 0, len(events))
	for _, event := range events {
		eventsRes = append(eventsRes, ToBasicEventResponse(event))
	}

	return eventsRes
}

func ToSimpleEventResponse(event *model.Event) *response.SimpleEventResponse {
	if event == nil {
		return nil
	}

	return &response.SimpleEventResponse{
		ID:          event.ID,
		Title:       event.Title,
		Thumbnail:   event.Thumbnail,
		IsActive:    event.IsActive,
		PublishedAt: event.PublishedAt,
		Branches:    ToBasicBranchesResponse(event.Branches),
		CreatedAt:   event.CreatedAt,
	}
}

func ToSimpleEventsResponse(events []*model.Event) []*response.SimpleEventResponse {
	if len(events) == 0 {
		return make([]*response.SimpleEventResponse, 0)
	}

	eventsRes := make([]*response.SimpleEventResponse, 0, len(events))
	for _, event := range events {
		eventsRes = append(eventsRes, ToSimpleEventResponse(event))
	}

	return eventsRes
}

func ToEventQuery(query request.EventPaginationQuery) model.EventQuery {
	return model.EventQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Sort:     query.Sort,
		Order:    query.Order,
		BranchID: query.BranchID,
		IsActive: query.IsActive,
		Search:   query.Search,
	}
}
