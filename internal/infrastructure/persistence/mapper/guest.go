package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToGuestEntity(guest *model.Guest) *entity.GuestEntity {
	if guest == nil {
		return nil
	}

	return &entity.GuestEntity{
		BaseEntity: entity.BaseEntity{
			ID:        guest.ID,
			CreatedAt: guest.CreatedAt,
			UpdatedAt: guest.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: guest.CreatedByID,
			UpdatedByID: guest.UpdatedByID,
		},
		FullName:       guest.FullName,
		National:       guest.National,
		DOB:            guest.DOB,
		PassportNumber: guest.PassportNumber,
		CheckIn:        guest.CheckIn,
		CheckOut:       guest.CheckOut,
		BookingRoomID:  guest.BookingRoomID,
	}
}

func ToGuest(guest *entity.GuestEntity) *model.Guest {
	if guest == nil {
		return nil
	}

	return &model.Guest{
		Base: model.Base{
			ID:        guest.ID,
			CreatedAt: guest.CreatedAt,
			UpdatedAt: guest.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: guest.CreatedByID,
			UpdatedByID: guest.UpdatedByID,
			CreatedBy:   ToUser(guest.CreatedBy),
			UpdatedBy:   ToUser(guest.UpdatedBy),
		},
		FullName:       guest.FullName,
		National:       guest.National,
		DOB:            guest.DOB,
		PassportNumber: guest.PassportNumber,
		CheckIn:        guest.CheckIn,
		CheckOut:       guest.CheckOut,
		BookingRoomID:  guest.BookingRoomID,
		BookingRoom:    ToBookingRoom(guest.BookingRoom),
	}
}

func ToGuests(guests []*entity.GuestEntity) []*model.Guest {
	if len(guests) == 0 {
		return make([]*model.Guest, 0)
	}

	guestsRes := make([]*model.Guest, 0, len(guests))
	for _, guest := range guests {
		guestsRes = append(guestsRes, ToGuest(guest))
	}

	return guestsRes
}
