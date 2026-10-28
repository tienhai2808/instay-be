package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToBasicSourceResponse(src *model.Source) *response.BasicSourceResponse {
	if src == nil {
		return nil
	}

	return &response.BasicSourceResponse{
		ID:   src.ID,
		Name: src.Name,
	}
}

func ToBookingRoomResponse(br *model.BookingRoom) *response.BookingRoomResponse {
	if br == nil {
		return nil
	}

	return &response.BookingRoomResponse{
		ID:             br.ID,
		OrderService:   ToBasicOrderServicesResponse(br.OrderServices),
		SupportTickets: ToSimpleSupportTicketsResponse(br.SupportTickets),
		Guests:         ToSimpleGuestsResponse(br.Guests),
		Review:         ToSimpleReviewResponse(br.Review),
	}
}

func ToSimpleReviewResponse(rv *model.Review) *response.SimpleReviewResponse {
	if rv == nil {
		return nil
	}

	return &response.SimpleReviewResponse{
		ID:        rv.ID,
		Star:      rv.Star,
		Content:   rv.Content,
		CreatedAt: rv.CreatedAt,
	}
}

func ToSimpleGuestResponse(guest *model.Guest) *response.SimpleGuestResponse {
	if guest == nil {
		return nil
	}

	return &response.SimpleGuestResponse{
		ID:             guest.ID,
		FullName:       guest.FullName,
		DOB:            guest.DOB,
		National:       guest.National,
		PassportNumber: guest.PassportNumber,
	}
}

func ToSimpleGuestsResponse(guests []*model.Guest) []*response.SimpleGuestResponse {
	if len(guests) == 0 {
		return make([]*response.SimpleGuestResponse, 0)
	}

	guestsRes := make([]*response.SimpleGuestResponse, 0, len(guests))
	for _, guest := range guests {
		guestsRes = append(guestsRes, ToSimpleGuestResponse(guest))
	}

	return guestsRes
}

func ToBasicSourcesResponse(srcs []*model.Source) []*response.BasicSourceResponse {
	if len(srcs) == 0 {
		return make([]*response.BasicSourceResponse, 0)
	}

	srcsRes := make([]*response.BasicSourceResponse, 0, len(srcs))
	for _, src := range srcs {
		srcsRes = append(srcsRes, ToBasicSourceResponse(src))
	}

	return srcsRes
}

func ToSourceResponse(src *model.Source) *response.SourceResponse {
	if src == nil {
		return nil
	}

	return &response.SourceResponse{
		ID:        src.ID,
		Name:      src.Name,
		CreatedAt: src.CreatedAt,
		UpdatedAt: src.UpdatedAt,
		CreatedBy: ToBasicUserResponse(src.CreatedBy),
		UpdatedBy: ToBasicUserResponse(src.UpdatedBy),
	}
}

func ToSourcesResponse(srcs []*model.Source) []*response.SourceResponse {
	if len(srcs) == 0 {
		return make([]*response.SourceResponse, 0)
	}

	srcsRes := make([]*response.SourceResponse, 0, len(srcs))
	for _, src := range srcs {
		srcsRes = append(srcsRes, ToSourceResponse(src))
	}

	return srcsRes
}

func ToBookingQuery(query request.BookingPaginationQuery) model.BookingQuery {
	return model.BookingQuery{
		Page:             query.Page,
		Limit:            query.Limit,
		ToDate:           query.ToDate,
		FromDate:         query.FromDate,
		Sort:             query.Sort,
		Order:            query.Order,
		Status:           query.Status,
		Search:           query.Search,
		SourceID:         query.SourceID,
		BranchID:         query.BranchID,
		UnassignedBranch: query.UnassignedBranch,
	}
}

func ToBookingTimelineQuery(query request.BookingTimelineQuery) model.BookingTimelineQuery {
	return model.BookingTimelineQuery{
		FromDate:   query.FromDate,
		ToDate:     query.ToDate,
		RoomTypeID: query.RoomTypeID,
		BranchID:   query.BranchID,
	}
}

func ToBookingStatusSummaryResponse(summary *model.BookingStatusSummary) *response.BookingStatusSummaryResponse {
	if summary == nil {
		return nil
	}

	return &response.BookingStatusSummaryResponse{
		All:        summary.All,
		CheckedIn:  summary.CheckedIn,
		Confirmed:  summary.Confirmed,
		CheckedOut: summary.CheckedOut,
		Cancelled:  summary.Cancelled,
	}
}

func ToBookingTimelineEventResponse(event *model.BookingTimelineEvent) *response.BookingTimelineEventResponse {
	if event == nil {
		return nil
	}

	return &response.BookingTimelineEventResponse{
		BookingID:     event.BookingID,
		CheckIn:       event.CheckIn,
		CheckOut:      event.CheckOut,
		Title:         event.Title,
		Status:        event.Status,
		BookingNumber: event.BookingNumber,
	}
}

func ToBookingTimelineEventsResponse(events []*model.BookingTimelineEvent) []*response.BookingTimelineEventResponse {
	if len(events) == 0 {
		return make([]*response.BookingTimelineEventResponse, 0)
	}

	eventsRes := make([]*response.BookingTimelineEventResponse, 0, len(events))
	for _, event := range events {
		eventsRes = append(eventsRes, ToBookingTimelineEventResponse(event))
	}

	return eventsRes
}

func ToBookingTimelineRoomResponse(room *model.BookingTimelineRoom) *response.BookingTimelineRoomResponse {
	if room == nil {
		return nil
	}

	return &response.BookingTimelineRoomResponse{
		RoomID:   room.RoomID,
		RoomName: room.RoomName,
		Events:   ToBookingTimelineEventsResponse(room.Events),
	}
}

func ToBookingTimelineRoomsResponse(rooms []*model.BookingTimelineRoom) []*response.BookingTimelineRoomResponse {
	if len(rooms) == 0 {
		return make([]*response.BookingTimelineRoomResponse, 0)
	}

	roomsRes := make([]*response.BookingTimelineRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomsRes = append(roomsRes, ToBookingTimelineRoomResponse(room))
	}

	return roomsRes
}

func ToBookingTimelineRoomTypeResponse(roomType *model.BookingTimelineRoomType) *response.BookingTimelineRoomTypeResponse {
	if roomType == nil {
		return nil
	}

	return &response.BookingTimelineRoomTypeResponse{
		RoomTypeID:   roomType.RoomTypeID,
		RoomTypeName: roomType.RoomTypeName,
		RoomsCount:   roomType.RoomsCount,
		Rooms:        ToBookingTimelineRoomsResponse(roomType.Rooms),
	}
}

func ToBookingTimelineRoomTypesResponse(roomTypes []*model.BookingTimelineRoomType) []*response.BookingTimelineRoomTypeResponse {
	if len(roomTypes) == 0 {
		return make([]*response.BookingTimelineRoomTypeResponse, 0)
	}

	roomTypesRes := make([]*response.BookingTimelineRoomTypeResponse, 0, len(roomTypes))
	for _, roomType := range roomTypes {
		roomTypesRes = append(roomTypesRes, ToBookingTimelineRoomTypeResponse(roomType))
	}

	return roomTypesRes
}

func ToBookingTimelineResponse(timeline *model.BookingTimeline) *response.BookingTimelineResponse {
	if timeline == nil {
		return nil
	}

	return &response.BookingTimelineResponse{
		StartDate: timeline.StartDate,
		EndDate:   timeline.EndDate,
		Status:    ToRoomStatusSummaryResponse(timeline.Status),
		Chart:     ToBookingTimelineRoomTypesResponse(timeline.Data),
	}
}

func ToBookingPredictResponse(predict model.BookingPredict) *response.BookingPredictResponse {
	if predict.RiskLevel == "" && predict.CancellationProbability == 0 && !predict.WillCancel {
		return nil
	}

	return &response.BookingPredictResponse{
		CancellationProbability: predict.CancellationProbability,
		WillCancel:              predict.WillCancel,
		RiskLevel:               predict.RiskLevel,
	}
}

func ToSimpleBookingResponse(booking *model.Booking) *response.SimpleBookingResponse {
	if booking == nil {
		return nil
	}

	return &response.SimpleBookingResponse{
		ID:             booking.ID,
		BookingNumber:  booking.BookingNumber,
		GuestFullName:  booking.GuestFullName,
		GuestEmail:     booking.GuestEmail,
		GuestPhone:     booking.GuestPhone,
		CheckIn:        booking.CheckIn,
		CheckOut:       booking.CheckOut,
		RoomNumber:     booking.RoomNumber,
		TotalSellPrice: booking.TotalSellPrice,
		RoomType:       ToBasicRoomTypeResponse(booking.RoomType),
		Branch:         ToBasicBranchResponse(booking.Branch),
		Source:         ToBasicSourceResponse(booking.Source),
		BookingRooms:   ToBasicBookingRoomsResponse(booking.BookingRooms),
		Status:         booking.Status,
		Predict:        ToBookingPredictResponse(booking.Predict),
	}
}

func ToSimpleBookingsResponse(bks []*model.Booking) []*response.SimpleBookingResponse {
	if len(bks) == 0 {
		return make([]*response.SimpleBookingResponse, 0)
	}

	bksRes := make([]*response.SimpleBookingResponse, 0, len(bks))
	for _, bk := range bks {
		bksRes = append(bksRes, ToSimpleBookingResponse(bk))
	}

	return bksRes
}

func ToBasicBookingRoomResponse(br *model.BookingRoom) *response.BasicBookingRoomResponse {
	if br == nil {
		return nil
	}

	return &response.BasicBookingRoomResponse{
		ID:   br.ID,
		Code: br.Code,
		Room: ToBasicRoomResponse(br.Room),
	}
}

func ToBasicBookingRoomsResponse(brs []*model.BookingRoom) []*response.BasicBookingRoomResponse {
	if len(brs) == 0 {
		return make([]*response.BasicBookingRoomResponse, 0)
	}

	brsRes := make([]*response.BasicBookingRoomResponse, 0, len(brs))
	for _, br := range brs {
		brsRes = append(brsRes, ToBasicBookingRoomResponse(br))
	}

	return brsRes
}

func ToBasicBookingResponse(booking *model.Booking) *response.BasicBookingResponse {
	if booking == nil {
		return nil
	}

	return &response.BasicBookingResponse{
		ID:            booking.ID,
		BookingNumber: booking.BookingNumber,
		GuestFullName: booking.GuestFullName,
		GuestEmail:    booking.GuestEmail,
		GuestPhone:    booking.GuestPhone,
		CheckIn:       booking.CheckIn,
		CheckOut:      booking.CheckOut,
	}
}

func ToSimpleBookingRoomResponse(br *model.BookingRoom) *response.SimpleBookingRoomResponse {
	if br == nil {
		return nil
	}

	return &response.SimpleBookingRoomResponse{
		ID:      br.ID,
		Booking: ToBasicBookingResponse(br.Booking),
		Room:    ToBasicRoomResponse(br.Room),
	}
}

func ToBookingResponse(booking *model.Booking) *response.BookingResponse {
	if booking == nil {
		return nil
	}

	return &response.BookingResponse{
		ID:                 booking.ID,
		BookingNumber:      booking.BookingNumber,
		GuestFullName:      booking.GuestFullName,
		GuestEmail:         booking.GuestEmail,
		GuestPhone:         booking.GuestPhone,
		CheckIn:            booking.CheckIn,
		CheckOut:           booking.CheckOut,
		RoomNumber:         booking.RoomNumber,
		GuestNumber:        booking.GuestNumber,
		BookedOn:           booking.BookedOn,
		TotalNetPrice:      booking.TotalNetPrice,
		TotalSellPrice:     booking.TotalSellPrice,
		PromotionName:      booking.PromotionName,
		MealPlan:           booking.MealPlan,
		BookingPreferences: booking.BookingPreferences,
		BookingConditions:  booking.BookingConditions,
		Note:               booking.Note,
		Status:             booking.Status,
		Branch:             ToBasicBranchResponse(booking.Branch),
		Source:             ToBasicSourceResponse(booking.Source),
		RoomType:           ToBasicRoomTypeResponse(booking.RoomType),
		CreatedAt:          booking.CreatedAt,
		UpdatedAt:          booking.UpdatedAt,
		CreatedBy:          ToBasicUserResponse(booking.CreatedBy),
		UpdatedBy:          ToBasicUserResponse(booking.UpdatedBy),
		BookingRooms:       ToBasicBookingRoomsResponse(booking.BookingRooms),
		Predict:            ToBookingPredictResponse(booking.Predict),
	}
}

func ToSourceQuery(query request.SourcePaginationQuery) model.SourceQuery {
	return model.SourceQuery{
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
	}
}
