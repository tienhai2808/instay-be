package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToRoomEntity(room *model.Room) *entity.RoomEntity {
	if room == nil {
		return nil
	}

	return &entity.RoomEntity{
		BaseEntity: entity.BaseEntity{
			ID:        room.ID,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: room.CreatedByID,
			UpdatedByID: room.UpdatedByID,
		},
		Name:       room.Name,
		FloorID:    room.FloorID,
		RoomTypeID: room.RoomTypeID,
	}
}

func ToRoomTypePrice(rtp *entity.RoomTypePriceEntity) *model.RoomTypePrice {
	if rtp == nil {
		return nil
	}

	return &model.RoomTypePrice{
		ID:           rtp.ID,
		RoomTypeID:   rtp.RoomTypeID,
		PricingTagID: rtp.PricingTagID,
		PricingType:  model.RoomTypePricePricingType(rtp.PricingType),
		PricingValue: rtp.PricingValue,
		PricingTag:   ToPricingTag(rtp.PricingTag),
	}
}

func ToRoomTypePriceEntity(rtp *model.RoomTypePrice) *entity.RoomTypePriceEntity {
	if rtp == nil {
		return nil
	}

	return &entity.RoomTypePriceEntity{
		ID:           rtp.ID,
		RoomTypeID:   rtp.RoomTypeID,
		PricingTagID: rtp.PricingTagID,
		PricingType:  string(rtp.PricingType),
		PricingValue: rtp.PricingValue,
	}
}

func ToPricingCalendar(pc *entity.PricingCalendarEntity) *model.PricingCalendar {
	if pc == nil {
		return nil
	}

	return &model.PricingCalendar{
		Date:         pc.Date,
		PricingTagID: pc.PricingTagID,
		PricingTag:   ToPricingTag(pc.PricingTag),
	}
}

func ToPricingCalendarEntity(pc *model.PricingCalendar) *entity.PricingCalendarEntity {
	if pc == nil {
		return nil
	}

	return &entity.PricingCalendarEntity{
		Date:         pc.Date,
		PricingTagID: pc.PricingTagID,
	}
}

func ToPricingTagEntity(pt *model.PricingTag) *entity.PricingTagEntity {
	if pt == nil {
		return nil
	}

	return &entity.PricingTagEntity{
		BaseEntity: entity.BaseEntity{
			ID:        pt.ID,
			CreatedAt: pt.CreatedAt,
			UpdatedAt: pt.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: pt.CreatedByID,
			UpdatedByID: pt.UpdatedByID,
		},
		Name: pt.Name,
	}
}

func ToPricingTag(pt *entity.PricingTagEntity) *model.PricingTag {
	if pt == nil {
		return nil
	}

	return &model.PricingTag{
		Base: model.Base{
			ID:        pt.ID,
			CreatedAt: pt.CreatedAt,
			UpdatedAt: pt.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: pt.CreatedByID,
			UpdatedByID: pt.UpdatedByID,
			CreatedBy:   ToUser(pt.CreatedBy),
			UpdatedBy:   ToUser(pt.UpdatedBy),
		},
		Name: pt.Name,
	}
}

func ToPricingTags(pts []*entity.PricingTagEntity) []*model.PricingTag {
	if len(pts) == 0 {
		return make([]*model.PricingTag, 0)
	}

	ptsRes := make([]*model.PricingTag, 0, len(pts))
	for _, pt := range pts {
		ptsRes = append(ptsRes, ToPricingTag(pt))
	}

	return ptsRes
}

func ToPricingCalendars(pcs []*entity.PricingCalendarEntity) []*model.PricingCalendar {
	if len(pcs) == 0 {
		return make([]*model.PricingCalendar, 0)
	}

	res := make([]*model.PricingCalendar, 0, len(pcs))
	for _, pc := range pcs {
		res = append(res, ToPricingCalendar(pc))
	}

	return res
}

func ToPricingCalendarEntities(pcs []*model.PricingCalendar) []*entity.PricingCalendarEntity {
	if len(pcs) == 0 {
		return make([]*entity.PricingCalendarEntity, 0)
	}

	res := make([]*entity.PricingCalendarEntity, 0, len(pcs))
	for _, pc := range pcs {
		if pc == nil {
			continue
		}
		res = append(res, ToPricingCalendarEntity(pc))
	}

	return res
}

func ToDailyRoomTypeRateEntity(rate *model.DailyRoomTypeRate) *entity.DailyRoomTypeRateEntity {
	if rate == nil {
		return nil
	}

	return &entity.DailyRoomTypeRateEntity{
		RoomTypeID:          rate.RoomTypeID,
		PricingCalendarDate: rate.PricingCalendarDate,
		FinalPrice:          rate.FinalPrice,
	}
}

func ToDailyRoomTypeRateEntities(rates []*model.DailyRoomTypeRate) []*entity.DailyRoomTypeRateEntity {
	if len(rates) == 0 {
		return make([]*entity.DailyRoomTypeRateEntity, 0)
	}

	res := make([]*entity.DailyRoomTypeRateEntity, 0, len(rates))
	for _, rate := range rates {
		if rate == nil {
			continue
		}
		res = append(res, ToDailyRoomTypeRateEntity(rate))
	}

	return res
}

func ToRoomTypePrices(rtps []*entity.RoomTypePriceEntity) []*model.RoomTypePrice {
	if len(rtps) == 0 {
		return make([]*model.RoomTypePrice, 0)
	}

	res := make([]*model.RoomTypePrice, 0, len(rtps))
	for _, rtp := range rtps {
		res = append(res, ToRoomTypePrice(rtp))
	}

	return res
}

func ToRoomTypeEntity(rt *model.RoomType) *entity.RoomTypeEntity {
	if rt == nil {
		return nil
	}

	return &entity.RoomTypeEntity{
		BaseEntity: entity.BaseEntity{
			ID:        rt.ID,
			CreatedAt: rt.CreatedAt,
			UpdatedAt: rt.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: rt.CreatedByID,
			UpdatedByID: rt.UpdatedByID,
		},
		Name:        rt.Name,
		Price:       rt.Price,
		SurfaceArea: rt.SurfaceArea,
		MaxPeople:   rt.MaxPeople,
		Branches:    ToBranchEntities(rt.Branches),
	}
}

func ToFloorEntity(fl *model.Floor) *entity.FloorEntity {
	if fl == nil {
		return nil
	}

	return &entity.FloorEntity{
		BaseEntity: entity.BaseEntity{
			ID:        fl.ID,
			CreatedAt: fl.CreatedAt,
			UpdatedAt: fl.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: fl.CreatedByID,
			UpdatedByID: fl.UpdatedByID,
		},
		Name:     fl.Name,
		BranchID: fl.BranchID,
	}
}

func ToRoom(room *entity.RoomEntity) *model.Room {
	if room == nil {
		return nil
	}

	return &model.Room{
		Base: model.Base{
			ID:        room.ID,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: room.CreatedByID,
			UpdatedByID: room.UpdatedByID,
			CreatedBy:   ToUser(room.CreatedBy),
			UpdatedBy:   ToUser(room.UpdatedBy),
		},
		Name:        room.Name,
		FloorID:     room.FloorID,
		RoomTypeID:  room.RoomTypeID,
		IsAvailable: room.IsAvailable,
		RoomType:    ToRoomType(room.RoomType),
		Floor:       ToFloor(room.Floor),
	}
}

func ToFloor(fl *entity.FloorEntity) *model.Floor {
	if fl == nil {
		return nil
	}

	return &model.Floor{
		Base: model.Base{
			ID:        fl.ID,
			CreatedAt: fl.CreatedAt,
			UpdatedAt: fl.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: fl.CreatedByID,
			UpdatedByID: fl.UpdatedByID,
			CreatedBy:   ToUser(fl.CreatedBy),
			UpdatedBy:   ToUser(fl.UpdatedBy),
		},
		Name:     fl.Name,
		BranchID: fl.BranchID,
		Branch:   ToBranch(fl.Branch),
	}
}

func ToFloors(fls []*entity.FloorEntity) []*model.Floor {
	if len(fls) == 0 {
		return make([]*model.Floor, 0)
	}

	flsRes := make([]*model.Floor, 0, len(fls))
	for _, fl := range fls {
		flsRes = append(flsRes, ToFloor(fl))
	}

	return flsRes
}

func ToRoomType(rt *entity.RoomTypeEntity) *model.RoomType {
	if rt == nil {
		return nil
	}

	return &model.RoomType{
		Base: model.Base{
			ID:        rt.ID,
			CreatedAt: rt.CreatedAt,
			UpdatedAt: rt.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: rt.CreatedByID,
			UpdatedByID: rt.UpdatedByID,
			CreatedBy:   ToUser(rt.CreatedBy),
			UpdatedBy:   ToUser(rt.UpdatedBy),
		},
		Name:           rt.Name,
		Price:          rt.Price,
		SurfaceArea:    rt.SurfaceArea,
		MaxPeople:      rt.MaxPeople,
		Branches:       ToBranches(rt.Branches),
		RoomTypePrices: ToRoomTypePrices(rt.RoomTypePrices),
	}
}

func ToRooms(rooms []*entity.RoomEntity) []*model.Room {
	if len(rooms) == 0 {
		return make([]*model.Room, 0)
	}

	roomsRes := make([]*model.Room, 0, len(rooms))
	for _, room := range rooms {
		roomsRes = append(roomsRes, ToRoom(room))
	}

	return roomsRes
}

func ToRoomTypes(rts []*entity.RoomTypeEntity) []*model.RoomType {
	if len(rts) == 0 {
		return make([]*model.RoomType, 0)
	}

	rtsRes := make([]*model.RoomType, 0, len(rts))
	for _, rt := range rts {
		rtsRes = append(rtsRes, ToRoomType(rt))
	}

	return rtsRes
}
