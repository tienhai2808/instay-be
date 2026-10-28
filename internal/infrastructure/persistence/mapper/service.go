package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToServiceEntity(svc *model.Service) *entity.ServiceEntity {
	if svc == nil {
		return nil
	}

	return &entity.ServiceEntity{
		BaseEntity: entity.BaseEntity{
			ID:        svc.ID,
			CreatedAt: svc.CreatedAt,
			UpdatedAt: svc.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: svc.CreatedByID,
			UpdatedByID: svc.UpdatedByID,
		},
		Name:            svc.Name,
		Slug:            svc.Slug,
		Description:     svc.Description,
		Price:           svc.Price,
		IsActive:        svc.IsActive,
		DiscountType:    (*string)(svc.DiscountType),
		DiscountValue:   svc.DiscountValue,
		DiscountStartAt: svc.DiscountStartAt,
		DiscountEndAt:   svc.DiscountEndAt,
		ServiceTypeID:   svc.ServiceTypeID,
		Branches:        ToBranchEntities(svc.Branches),
	}
}

func ToServiceTypeEntity(st *model.ServiceType) *entity.ServiceTypeEntity {
	if st == nil {
		return nil
	}

	return &entity.ServiceTypeEntity{
		BaseEntity: entity.BaseEntity{
			ID:        st.ID,
			CreatedAt: st.CreatedAt,
			UpdatedAt: st.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: st.CreatedByID,
			UpdatedByID: st.UpdatedByID,
		},
		Name:         st.Name,
		Icon:         st.Icon,
		Slug:         st.Slug,
		DepartmentID: st.DepartmentID,
	}
}

func ToServiceImageEntity(si *model.ServiceImage) *entity.ServiceImageEntity {
	if si == nil {
		return nil
	}

	return &entity.ServiceImageEntity{
		ServiceID:   si.ServiceID,
		Key:         si.Key,
		IsThumbnail: si.IsThumbnail,
		SortOrder:   si.SortOrder,
	}
}

func ToServiceImageEntities(sis []*model.ServiceImage) []*entity.ServiceImageEntity {
	if len(sis) == 0 {
		return make([]*entity.ServiceImageEntity, 0)
	}

	sisRes := make([]*entity.ServiceImageEntity, 0, len(sis))
	for _, si := range sis {
		sisRes = append(sisRes, ToServiceImageEntity(si))
	}

	return sisRes
}

func ToServiceTypes(sts []*entity.ServiceTypeEntity) []*model.ServiceType {
	if len(sts) == 0 {
		return make([]*model.ServiceType, 0)
	}

	stsRes := make([]*model.ServiceType, 0, len(sts))
	for _, st := range sts {
		stsRes = append(stsRes, ToServiceType(st))
	}

	return stsRes
}

func ToOrderServiceEntity(os *model.OrderService) *entity.OrderServiceEntity {
	if os == nil {
		return nil
	}

	return &entity.OrderServiceEntity{
		BaseEntity: entity.BaseEntity{
			ID:        os.ID,
			CreatedAt: os.CreatedAt,
			UpdatedAt: os.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: os.CreatedByID,
			UpdatedByID: os.UpdatedByID,
		},
		ServiceID:     os.ServiceID,
		BookingRoomID: os.BookingRoomID,
		Quantity:      os.Quantity,
		UnitPrice:     os.UnitPrice,
		TotalPrice:    os.TotalPrice,
		Status:        string(os.Status),
		GuestNote:     os.GuestNote,
		StaffNote:     os.StaffNote,
		CancelReason:  os.CancelReason,
		EstimatedAt:   os.EstimatedAt,
	}
}

func ToOrderService(os *entity.OrderServiceEntity) *model.OrderService {
	if os == nil {
		return nil
	}

	return &model.OrderService{
		Base: model.Base{
			ID:        os.ID,
			CreatedAt: os.CreatedAt,
			UpdatedAt: os.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: os.CreatedByID,
			UpdatedByID: os.UpdatedByID,
			CreatedBy:   ToUser(os.CreatedBy),
			UpdatedBy:   ToUser(os.UpdatedBy),
		},
		ServiceID:     os.ServiceID,
		BookingRoomID: os.BookingRoomID,
		Quantity:      os.Quantity,
		UnitPrice:     os.UnitPrice,
		TotalPrice:    os.TotalPrice,
		Status:        model.OrderServiceStatus(os.Status),
		GuestNote:     os.GuestNote,
		StaffNote:     os.StaffNote,
		CancelReason:  os.CancelReason,
		EstimatedAt:   os.EstimatedAt,
		Service:       ToService(os.Service),
		BookingRoom:   ToBookingRoom(os.BookingRoom),
	}
}

func ToService(svc *entity.ServiceEntity) *model.Service {
	if svc == nil {
		return nil
	}

	return &model.Service{
		Base: model.Base{
			ID:        svc.ID,
			CreatedAt: svc.CreatedAt,
			UpdatedAt: svc.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: svc.CreatedByID,
			UpdatedByID: svc.UpdatedByID,
			CreatedBy:   ToUser(svc.CreatedBy),
			UpdatedBy:   ToUser(svc.UpdatedBy),
		},
		Name:            svc.Name,
		Slug:            svc.Slug,
		Description:     svc.Description,
		Price:           svc.Price,
		IsActive:        svc.IsActive,
		DiscountType:    (*model.ServiceDiscountType)(svc.DiscountType),
		DiscountValue:   svc.DiscountValue,
		DiscountStartAt: svc.DiscountStartAt,
		DiscountEndAt:   svc.DiscountEndAt,
		ServiceTypeID:   svc.ServiceTypeID,
		ServiceType:     ToServiceType(svc.ServiceType),
		ServiceImages:   ToServiceImages(svc.ServiceImages),
		Branches:        ToBranches(svc.Branches),
	}
}

func ToServices(svcs []*entity.ServiceEntity) []*model.Service {
	if len(svcs) == 0 {
		return make([]*model.Service, 0)
	}

	svcsRes := make([]*model.Service, 0, len(svcs))
	for _, svc := range svcs {
		svcsRes = append(svcsRes, ToService(svc))
	}

	return svcsRes
}

func ToServiceImage(si *entity.ServiceImageEntity) *model.ServiceImage {
	if si == nil {
		return nil
	}

	return &model.ServiceImage{
		ServiceID:   si.ServiceID,
		Key:         si.Key,
		IsThumbnail: si.IsThumbnail,
		SortOrder:   si.SortOrder,
		Service:     ToService(si.Service),
	}
}

func ToServiceImages(sis []*entity.ServiceImageEntity) []*model.ServiceImage {
	if len(sis) == 0 {
		return make([]*model.ServiceImage, 0)
	}

	sisRes := make([]*model.ServiceImage, 0, len(sis))
	for _, si := range sis {
		sisRes = append(sisRes, ToServiceImage(si))
	}

	return sisRes
}

func ToServiceType(st *entity.ServiceTypeEntity) *model.ServiceType {
	if st == nil {
		return nil
	}

	return &model.ServiceType{
		Base: model.Base{
			ID:        st.ID,
			CreatedAt: st.CreatedAt,
			UpdatedAt: st.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: st.CreatedByID,
			UpdatedByID: st.UpdatedByID,
			CreatedBy:   ToUser(st.CreatedBy),
			UpdatedBy:   ToUser(st.UpdatedBy),
		},
		Name:         st.Name,
		Icon:         st.Icon,
		Slug:         st.Slug,
		DepartmentID: st.DepartmentID,
		Department:   ToDepartment(st.Department),
		Services:     ToServices(st.Services),
	}
}

func ToOrderServices(oss []*entity.OrderServiceEntity) []*model.OrderService {
	if len(oss) == 0 {
		return make([]*model.OrderService, 0)
	}

	ossRes := make([]*model.OrderService, 0, len(oss))
	for _, os := range oss {
		ossRes = append(ossRes, ToOrderService(os))
	}

	return ossRes
}
