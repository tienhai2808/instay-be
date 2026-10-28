package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToSupportCategoryEntity(sc *model.SupportCategory) *entity.SupportCategoryEntity {
	if sc == nil {
		return nil
	}

	return &entity.SupportCategoryEntity{
		BaseEntity: entity.BaseEntity{
			ID:        sc.ID,
			CreatedAt: sc.CreatedAt,
			UpdatedAt: sc.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: sc.CreatedByID,
			UpdatedByID: sc.UpdatedByID,
		},
		Name:         sc.Name,
		DepartmentID: sc.DepartmentID,
	}
}

func ToSupportCategory(sc *entity.SupportCategoryEntity) *model.SupportCategory {
	if sc == nil {
		return nil
	}

	return &model.SupportCategory{
		Base: model.Base{
			ID:        sc.ID,
			CreatedAt: sc.CreatedAt,
			UpdatedAt: sc.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: sc.CreatedByID,
			UpdatedByID: sc.UpdatedByID,
			CreatedBy:   ToUser(sc.CreatedBy),
			UpdatedBy:   ToUser(sc.UpdatedBy),
		},
		Name:                sc.Name,
		DepartmentID:        sc.DepartmentID,
		Department:          ToDepartment(sc.Department),
		SupportTickets:      ToSupportTickets(sc.SupportTickets),
		SupportTicketsCount: sc.SupportTicketsCount,
	}
}

func ToSupportCategories(scs []*entity.SupportCategoryEntity) []*model.SupportCategory {
	if len(scs) == 0 {
		return make([]*model.SupportCategory, 0)
	}

	scsRes := make([]*model.SupportCategory, 0, len(scs))
	for _, sc := range scs {
		scsRes = append(scsRes, ToSupportCategory(sc))
	}

	return scsRes
}

func ToSupportTicketEntity(st *model.SupportTicket) *entity.SupportTicketEntity {
	if st == nil {
		return nil
	}

	return &entity.SupportTicketEntity{
		BaseEntity: entity.BaseEntity{
			ID:        st.ID,
			CreatedAt: st.CreatedAt,
			UpdatedAt: st.UpdatedAt,
		},
		BookingRoomID:     st.BookingRoomID,
		Content:           st.Content,
		Status:            string(st.Status),
		Priority:          string(st.Priority),
		ResolutionSeconds: st.ResolutionSeconds,
		SupportCategoryID: st.SupportCategoryID,
		UpdatedByID:       st.UpdatedByID,
	}
}

func ToSupportTicket(st *entity.SupportTicketEntity) *model.SupportTicket {
	if st == nil {
		return nil
	}

	return &model.SupportTicket{
		Base: model.Base{
			ID:        st.ID,
			CreatedAt: st.CreatedAt,
			UpdatedAt: st.UpdatedAt,
		},
		BookingRoomID:     st.BookingRoomID,
		Content:           st.Content,
		Status:            model.SupportTicketStatus(st.Status),
		Priority:          model.SupportTicketPriority(st.Priority),
		ResolutionSeconds: st.ResolutionSeconds,
		SupportCategoryID: st.SupportCategoryID,
		UpdatedByID:       st.UpdatedByID,
		BookingRoom:       ToBookingRoom(st.BookingRoom),
		SupportCategory:   ToSupportCategory(st.SupportCategory),
		UpdatedBy:         ToUser(st.UpdatedBy),
	}
}

func ToSupportTickets(sts []*entity.SupportTicketEntity) []*model.SupportTicket {
	if len(sts) == 0 {
		return make([]*model.SupportTicket, 0)
	}

	stsRes := make([]*model.SupportTicket, 0, len(sts))
	for _, st := range sts {
		stsRes = append(stsRes, ToSupportTicket(st))
	}

	return stsRes
}
