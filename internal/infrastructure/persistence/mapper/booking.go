package mapper

import (
	"encoding/json"

	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
	"gorm.io/datatypes"
)

func ToBookingEntity(bk *model.Booking) *entity.BookingEntity {
	if bk == nil {
		return nil
	}

	var predictJSON datatypes.JSON
	if bk.Predict.RiskLevel != "" || bk.Predict.CancellationProbability != 0 || bk.Predict.WillCancel {
		bytes, _ := json.Marshal(bk.Predict)
		predictJSON = datatypes.JSON(bytes)
	}

	return &entity.BookingEntity{
		BaseEntity: entity.BaseEntity{
			ID:        bk.ID,
			CreatedAt: bk.CreatedAt,
			UpdatedAt: bk.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: bk.CreatedByID,
			UpdatedByID: bk.UpdatedByID,
		},
		BookingNumber:      bk.BookingNumber,
		GuestFullName:      bk.GuestFullName,
		GuestEmail:         bk.GuestEmail,
		GuestPhone:         bk.GuestPhone,
		CheckIn:            bk.CheckIn,
		CheckOut:           bk.CheckOut,
		RoomNumber:         bk.RoomNumber,
		GuestNumber:        bk.GuestNumber,
		BookedOn:           bk.BookedOn,
		TotalNetPrice:      bk.TotalNetPrice,
		TotalSellPrice:     bk.TotalSellPrice,
		PromotionName:      bk.PromotionName,
		MealPlan:           bk.MealPlan,
		BookingPreferences: bk.BookingPreferences,
		BookingConditions:  bk.BookingConditions,
		Note:               bk.Note,
		Status:             string(bk.Status),
		BranchID:           bk.BranchID,
		SourceID:           bk.SourceID,
		RoomTypeID:         bk.RoomTypeID,
		Predict:            predictJSON,
	}
}

func ToBooking(bk *entity.BookingEntity) *model.Booking {
	if bk == nil {
		return nil
	}

	var predict model.BookingPredict
	if len(bk.Predict) > 0 {
		_ = json.Unmarshal(bk.Predict, &predict)
	}

	return &model.Booking{
		Base: model.Base{
			ID:        bk.ID,
			CreatedAt: bk.CreatedAt,
			UpdatedAt: bk.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: bk.CreatedByID,
			UpdatedByID: bk.UpdatedByID,
			CreatedBy:   ToUser(bk.CreatedBy),
			UpdatedBy:   ToUser(bk.UpdatedBy),
		},
		BookingNumber:      bk.BookingNumber,
		GuestFullName:      bk.GuestFullName,
		GuestEmail:         bk.GuestEmail,
		GuestPhone:         bk.GuestPhone,
		CheckIn:            bk.CheckIn,
		CheckOut:           bk.CheckOut,
		RoomNumber:         bk.RoomNumber,
		GuestNumber:        bk.GuestNumber,
		BookedOn:           bk.BookedOn,
		TotalNetPrice:      bk.TotalNetPrice,
		TotalSellPrice:     bk.TotalSellPrice,
		PromotionName:      bk.PromotionName,
		MealPlan:           bk.MealPlan,
		BookingPreferences: bk.BookingPreferences,
		BookingConditions:  bk.BookingConditions,
		Note:               bk.Note,
		Status:             model.BookingStatus(bk.Status),
		BranchID:           bk.BranchID,
		SourceID:           bk.SourceID,
		RoomTypeID:         bk.RoomTypeID,
		Predict:            predict,
		Branch:             ToBranch(bk.Branch),
		Source:             ToSource(bk.Source),
		RoomType:           ToRoomType(bk.RoomType),
		BookingRooms:       ToBookingRooms(bk.BookingRooms),
	}
}

func ToSourceEntity(src *model.Source) *entity.SourceEntity {
	if src == nil {
		return nil
	}

	return &entity.SourceEntity{
		BaseEntity: entity.BaseEntity{
			ID:        src.ID,
			CreatedAt: src.CreatedAt,
			UpdatedAt: src.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: src.CreatedByID,
			UpdatedByID: src.UpdatedByID,
		},
		Name: src.Name,
	}
}

func ToSource(src *entity.SourceEntity) *model.Source {
	if src == nil {
		return nil
	}

	return &model.Source{
		Base: model.Base{
			ID:        src.ID,
			CreatedAt: src.CreatedAt,
			UpdatedAt: src.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: src.CreatedByID,
			UpdatedByID: src.UpdatedByID,
			CreatedBy:   ToUser(src.CreatedBy),
			UpdatedBy:   ToUser(src.UpdatedBy),
		},
		Name: src.Name,
	}
}

func ToSources(srcs []*entity.SourceEntity) []*model.Source {
	if len(srcs) == 0 {
		return make([]*model.Source, 0)
	}

	srcsRes := make([]*model.Source, 0, len(srcs))
	for _, src := range srcs {
		srcsRes = append(srcsRes, ToSource(src))
	}

	return srcsRes
}

func ToBookingRoom(br *entity.BookingRoomEntity) *model.BookingRoom {
	if br == nil {
		return nil
	}

	return &model.BookingRoom{
		ID:             br.ID,
		BookingID:      br.BookingID,
		RoomID:         br.RoomID,
		Code:           br.Code,
		Status:         model.BookingRoomStatus(br.Status),
		Room:           ToRoom(br.Room),
		Booking:        ToBooking(br.Booking),
		OrderServices:  ToOrderServices(br.OrderServices),
		SupportTickets: ToSupportTickets(br.SupportTickets),
		Guests:         ToGuests(br.Guests),
		Customers:      ToCustomers(br.Customers),
		Review:         ToReview(br.Review),
		Chats:          ToChats(br.Chats),
	}
}

func ToBookingRooms(brs []*entity.BookingRoomEntity) []*model.BookingRoom {
	if len(brs) == 0 {
		return make([]*model.BookingRoom, 0)
	}

	brsRes := make([]*model.BookingRoom, 0, len(brs))
	for _, br := range brs {
		brsRes = append(brsRes, ToBookingRoom(br))
	}

	return brsRes
}

func ToBookings(bks []*entity.BookingEntity) []*model.Booking {
	if len(bks) == 0 {
		return make([]*model.Booking, 0)
	}

	bksRes := make([]*model.Booking, 0, len(bks))
	for _, bk := range bks {
		bksRes = append(bksRes, ToBooking(bk))
	}

	return bksRes
}

func ToBookingRoomEntity(br *model.BookingRoom) *entity.BookingRoomEntity {
	if br == nil {
		return nil
	}

	return &entity.BookingRoomEntity{
		ID:        br.ID,
		BookingID: br.BookingID,
		RoomID:    br.RoomID,
		Code:      br.Code,
		Status:    string(br.Status),
	}
}

func ToBookingRoomEntities(brs []*model.BookingRoom) []*entity.BookingRoomEntity {
	if len(brs) == 0 {
		return make([]*entity.BookingRoomEntity, 0)
	}

	brsRes := make([]*entity.BookingRoomEntity, 0, len(brs))
	for _, br := range brs {
		brsRes = append(brsRes, ToBookingRoomEntity(br))
	}

	return brsRes
}
