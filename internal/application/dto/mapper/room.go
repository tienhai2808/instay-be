package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToPricingCalendarResponse(pc *model.PricingCalendar) *response.PricingCalendarResponse {
	if pc == nil {
		return nil
	}

	return &response.PricingCalendarResponse{
		Date:         pc.Date,
		PricingTagID: pc.PricingTagID,
		PricingTag:   ToBasicPricingTagResponse(pc.PricingTag),
	}
}

func ToPricingTagQuery(query request.PricingTagPaginationQuery) model.PricingTagQuery {
	return model.PricingTagQuery{
		Page:   query.Page,
		Limit:  query.Limit,
		Search: query.Search,
	}
}

func ToBasicPricingTagResponse(pricingTag *model.PricingTag) *response.BasicPricingTagResponse {
	if pricingTag == nil {
		return nil
	}

	return &response.BasicPricingTagResponse{
		ID:   pricingTag.ID,
		Name: pricingTag.Name,
	}
}

func ToBasicPricingTagsResponse(pricingTags []*model.PricingTag) []*response.BasicPricingTagResponse {
	if len(pricingTags) == 0 {
		return make([]*response.BasicPricingTagResponse, 0)
	}

	pricingTagsRes := make([]*response.BasicPricingTagResponse, 0, len(pricingTags))
	for _, pricingTag := range pricingTags {
		pricingTagsRes = append(pricingTagsRes, ToBasicPricingTagResponse(pricingTag))
	}

	return pricingTagsRes
}

func ToPricingTagResponse(pricingTag *model.PricingTag) *response.PricingTagResponse {
	if pricingTag == nil {
		return nil
	}

	return &response.PricingTagResponse{
		ID:        pricingTag.ID,
		Name:      pricingTag.Name,
		CreatedAt: pricingTag.CreatedAt,
		UpdatedAt: pricingTag.UpdatedAt,
		CreatedBy: ToBasicUserResponse(pricingTag.CreatedBy),
		UpdatedBy: ToBasicUserResponse(pricingTag.UpdatedBy),
	}
}

func ToPricingTagsResponse(pricingTags []*model.PricingTag) []*response.PricingTagResponse {
	if len(pricingTags) == 0 {
		return make([]*response.PricingTagResponse, 0)
	}

	pricingTagsRes := make([]*response.PricingTagResponse, 0, len(pricingTags))
	for _, pricingTag := range pricingTags {
		pricingTagsRes = append(pricingTagsRes, ToPricingTagResponse(pricingTag))
	}

	return pricingTagsRes
}

func ToPricingCalendarsResponse(pcs []*model.PricingCalendar) []*response.PricingCalendarResponse {
	if len(pcs) == 0 {
		return make([]*response.PricingCalendarResponse, 0)
	}

	res := make([]*response.PricingCalendarResponse, 0, len(pcs))
	for _, pc := range pcs {
		res = append(res, ToPricingCalendarResponse(pc))
	}

	return res
}

func ToBasicRoomTypeResponse(roomType *model.RoomType) *response.BasicRoomTypeResponse {
	if roomType == nil {
		return nil
	}

	return &response.BasicRoomTypeResponse{
		ID:   roomType.ID,
		Name: roomType.Name,
	}
}

func ToRoomTypePriceResponse(rtp *model.RoomTypePrice) *response.RoomTypePriceResponse {
	if rtp == nil {
		return nil
	}

	return &response.RoomTypePriceResponse{
		ID:           rtp.ID,
		RoomTypeID:   rtp.RoomTypeID,
		PricingTagID: rtp.PricingTagID,
		PricingType:  string(rtp.PricingType),
		PricingValue: rtp.PricingValue,
		PricingTag:   ToBasicPricingTagResponse(rtp.PricingTag),
	}
}

func ToRoomTypePricesResponse(rtps []*model.RoomTypePrice) []*response.RoomTypePriceResponse {
	if len(rtps) == 0 {
		return make([]*response.RoomTypePriceResponse, 0)
	}

	res := make([]*response.RoomTypePriceResponse, 0, len(rtps))
	for _, rtp := range rtps {
		res = append(res, ToRoomTypePriceResponse(rtp))
	}
	return res
}

func ToBasicRoomResponse(room *model.Room) *response.BasicRoomResponse {
	if room == nil {
		return nil
	}

	return &response.BasicRoomResponse{
		ID:   room.ID,
		Name: room.Name,
	}
}

func ToBasicRoomTypesResponse(roomTypes []*model.RoomType) []*response.BasicRoomTypeResponse {
	if len(roomTypes) == 0 {
		return make([]*response.BasicRoomTypeResponse, 0)
	}

	roomTypesRes := make([]*response.BasicRoomTypeResponse, 0, len(roomTypes))
	for _, roomType := range roomTypes {
		roomTypesRes = append(roomTypesRes, ToBasicRoomTypeResponse(roomType))
	}

	return roomTypesRes
}

func ToRoomTypeResponse(roomType *model.RoomType) *response.RoomTypeResponse {
	if roomType == nil {
		return nil
	}

	return &response.RoomTypeResponse{
		ID:             roomType.ID,
		Name:           roomType.Name,
		Price:          roomType.Price,
		SurfaceArea:    roomType.SurfaceArea,
		MaxPeople:      roomType.MaxPeople,
		CreatedAt:      roomType.CreatedAt,
		UpdatedAt:      roomType.UpdatedAt,
		CreatedBy:      ToBasicUserResponse(roomType.CreatedBy),
		UpdatedBy:      ToBasicUserResponse(roomType.UpdatedBy),
		RoomsCount:     roomType.RoomsCount,
		RoomTypePrices: ToRoomTypePricesResponse(roomType.RoomTypePrices),
		Branches:       ToBasicBranchesResponse(roomType.Branches),
	}
}

func ToSimpleRoomResponse(room *model.Room) *response.SimpleRoomResponse {
	if room == nil {
		return nil
	}

	return &response.SimpleRoomResponse{
		ID:          room.ID,
		Name:        room.Name,
		IsAvailable: room.IsAvailable,
	}
}

func ToSimpleRoomsResponse(rooms []*model.Room) []*response.SimpleRoomResponse {
	if len(rooms) == 0 {
		return make([]*response.SimpleRoomResponse, 0)
	}

	roomsRes := make([]*response.SimpleRoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomsRes = append(roomsRes, ToSimpleRoomResponse(room))
	}

	return roomsRes
}

func ToRoomTypesResponse(roomTypes []*model.RoomType) []*response.RoomTypeResponse {
	if len(roomTypes) == 0 {
		return make([]*response.RoomTypeResponse, 0)
	}

	roomTypesRes := make([]*response.RoomTypeResponse, 0, len(roomTypes))
	for _, roomType := range roomTypes {
		roomTypesRes = append(roomTypesRes, ToRoomTypeResponse(roomType))
	}

	return roomTypesRes
}

func ToBasicFloorResponse(floor *model.Floor) *response.BasicFloorResponse {
	if floor == nil {
		return nil
	}

	return &response.BasicFloorResponse{
		ID:     floor.ID,
		Name:   floor.Name,
		Branch: ToBasicBranchResponse(floor.Branch),
	}
}

func ToBasicFloorsResponse(floors []*model.Floor) []*response.BasicFloorResponse {
	if len(floors) == 0 {
		return make([]*response.BasicFloorResponse, 0)
	}

	floorsRes := make([]*response.BasicFloorResponse, 0, len(floors))
	for _, floor := range floors {
		floorsRes = append(floorsRes, ToBasicFloorResponse(floor))
	}

	return floorsRes
}

func ToRoomResponse(room *model.Room) *response.RoomResponse {
	if room == nil {
		return nil
	}

	return &response.RoomResponse{
		ID:        room.ID,
		Name:      room.Name,
		RoomType:  ToBasicRoomTypeResponse(room.RoomType),
		Floor:     ToBasicFloorResponse(room.Floor),
		CreatedAt: room.CreatedAt,
		UpdatedAt: room.UpdatedAt,
		CreatedBy: ToBasicUserResponse(room.CreatedBy),
		UpdatedBy: ToBasicUserResponse(room.UpdatedBy),
	}
}

func ToRoomsResponse(rooms []*model.Room) []*response.RoomResponse {
	if len(rooms) == 0 {
		return make([]*response.RoomResponse, 0)
	}

	roomsRes := make([]*response.RoomResponse, 0, len(rooms))
	for _, room := range rooms {
		roomsRes = append(roomsRes, ToRoomResponse(room))
	}

	return roomsRes
}

func ToFloorQuery(query request.FloorPaginationQuery) model.FloorQuery {
	return model.FloorQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Search:   query.Search,
		BranchID: query.BranchID,
	}
}

func ToRoomQuery(query request.RoomPaginationQuery) model.RoomQuery {
	return model.RoomQuery{
		Page:       query.Page,
		Limit:      query.Limit,
		Sort:       query.Sort,
		Order:      query.Order,
		BranchID:   query.BranchID,
		FloorID:    query.FloorID,
		RoomTypeID: query.RoomTypeID,
		Search:     query.Search,
	}
}

func ToRoomTimelineQuery(query request.RoomTimelineQuery) model.RoomTimelineQuery {
	return model.RoomTimelineQuery{
		Search:     query.Search,
		BranchID:   query.BranchID,
		RoomTypeID: query.RoomTypeID,
		FloorID:    query.FloorID,
	}
}

func ToRoomTypeQuery(query request.RoomTypePaginationQuery) model.RoomTypeQuery {
	return model.RoomTypeQuery{
		Page:     query.Page,
		Limit:    query.Limit,
		Search:   query.Search,
		BranchID: query.BranchID,
	}
}

func ToFloorResponse(floor *model.Floor) *response.FloorResponse {
	if floor == nil {
		return nil
	}

	return &response.FloorResponse{
		ID:        floor.ID,
		Name:      floor.Name,
		Branch:    ToBasicBranchResponse(floor.Branch),
		CreatedAt: floor.CreatedAt,
		UpdatedAt: floor.UpdatedAt,
		CreatedBy: ToBasicUserResponse(floor.CreatedBy),
		UpdatedBy: ToBasicUserResponse(floor.UpdatedBy),
	}
}

func ToFloorsResponse(floors []*model.Floor) []*response.FloorResponse {
	if len(floors) == 0 {
		return make([]*response.FloorResponse, 0)
	}

	floorsRes := make([]*response.FloorResponse, 0, len(floors))
	for _, floor := range floors {
		floorsRes = append(floorsRes, ToFloorResponse(floor))
	}

	return floorsRes
}

func ToRoomStatusSummaryResponse(summary *model.RoomStatusSummary) *response.RoomStatusSummaryResponse {
	if summary == nil {
		return nil
	}

	return &response.RoomStatusSummaryResponse{
		All:       summary.All,
		Available: summary.Available,
		Occupied:  summary.Occupied,
		Booked:    summary.Booked,
	}
}

func ToRoomTimelineBookingResponse(booking *model.RoomTimelineBooking) *response.RoomTimelineBookingResponse {
	if booking == nil {
		return nil
	}

	return &response.RoomTimelineBookingResponse{
		BookingID:     booking.BookingID,
		BookingNumber: booking.BookingNumber,
		Status:        booking.Status,
		SourceName:    booking.SourceName,
		GuestName:     booking.GuestName,
		CheckIn:       booking.CheckIn,
		CheckOut:      booking.CheckOut,
	}
}

func ToRoomTimelineRoomResponse(room *model.RoomTimelineRoom) *response.RoomTimelineRoomResponse {
	if room == nil {
		return nil
	}

	return &response.RoomTimelineRoomResponse{
		RoomID:   room.RoomID,
		RoomName: room.RoomName,
		Price:    room.Price,
		Booking:  ToRoomTimelineBookingResponse(room.Booking),
	}
}

func ToRoomTimelineRoomTypeResponse(roomType *model.RoomTimelineRoomType) *response.RoomTimelineRoomTypeResponse {
	if roomType == nil {
		return nil
	}

	rooms := make([]*response.RoomTimelineRoomResponse, 0, len(roomType.Rooms))
	for _, room := range roomType.Rooms {
		rooms = append(rooms, ToRoomTimelineRoomResponse(room))
	}

	return &response.RoomTimelineRoomTypeResponse{
		RoomTypeID:   roomType.RoomTypeID,
		RoomTypeName: roomType.RoomTypeName,
		RoomsCount:   roomType.RoomsCount,
		Rooms:        rooms,
	}
}

func ToRoomTimelineResponse(timeline *model.RoomTimeline) *response.RoomTimelineResponse {
	if timeline == nil {
		return nil
	}

	roomTypes := make([]*response.RoomTimelineRoomTypeResponse, 0, len(timeline.Data))
	for _, roomType := range timeline.Data {
		roomTypes = append(roomTypes, ToRoomTimelineRoomTypeResponse(roomType))
	}

	return &response.RoomTimelineResponse{
		ServerNow: timeline.ServerNow,
		Status:    ToRoomTimelineStatusSummaryResponse(timeline.Status),
		Timeline:  roomTypes,
	}
}

func ToRoomTimelineStatusSummaryResponse(summary *model.RoomTimelineStatusSummary) *response.RoomTimelineStatusSummaryResponse {
	if summary == nil {
		return nil
	}

	return &response.RoomTimelineStatusSummaryResponse{
		All:             summary.All,
		Available:       summary.Available,
		Occupied:        summary.Occupied,
		UpcomingCheckIn: summary.UpcomingCheckIn,
		OverdueCheckout: summary.OverdueCheckout,
		OverdueCheckIn:  summary.OverdueCheckIn,
	}
}
