package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/application/dto/request"
	"github.com/InstaySystem/is_v2-be/internal/application/dto/response"
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
)

func ToNotificationResponse(nt *model.Notification) *response.NotificationResponse {
	if nt == nil {
		return nil
	}

	return &response.NotificationResponse{
		ID:          nt.ID,
		ContentType: nt.ContentType,
		ContentKey:  nt.ContentKey,
		ContentVI:   nt.ContentVI,
		ContentEN:   nt.ContentEN,
		IsRead:      nt.IsRead,
		CreatedAt:   nt.CreatedAt,
	}
}

func ToNotificationsResponse(nts []*model.Notification) []*response.NotificationResponse {
	if len(nts) == 0 {
		return make([]*response.NotificationResponse, 0)
	}

	ntsRes := make([]*response.NotificationResponse, 0, len(nts))
	for _, nt := range nts {
		ntsRes = append(ntsRes, ToNotificationResponse(nt))
	}

	return ntsRes
}

func ToNotificationQuery(query request.NotificationPaginationQuery) model.NotificationQuery {
	return model.NotificationQuery{
		Page:  query.Page,
		Limit: query.Limit,
	}
}
