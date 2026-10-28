package mapper

import (
	"github.com/InstaySystem/is_v2-be/internal/domain/model"
	"github.com/InstaySystem/is_v2-be/internal/infrastructure/persistence/orm/entity"
)

func ToEventEntity(event *model.Event) *entity.EventEntity {
	if event == nil {
		return nil
	}

	return &entity.EventEntity{
		BaseEntity: entity.BaseEntity{
			ID:        event.ID,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		},
		AuditEntity: entity.AuditEntity{
			CreatedByID: event.CreatedByID,
			UpdatedByID: event.UpdatedByID,
		},
		Thumbnail:          event.Thumbnail,
		Title:              event.Title,
		Slug:               event.Slug,
		Content:            event.Content,
		Excerpt:            event.Excerpt,
		IsActive:           event.IsActive,
		PublishedAt:        event.PublishedAt,
		IsNotificationSent: event.IsNotificationSent,
		Branches:           ToBranchEntities(event.Branches),
	}
}

func ToEvent(event *entity.EventEntity) *model.Event {
	if event == nil {
		return nil
	}

	return &model.Event{
		Base: model.Base{
			ID:        event.ID,
			CreatedAt: event.CreatedAt,
			UpdatedAt: event.UpdatedAt,
		},
		Audit: model.Audit{
			CreatedByID: event.CreatedByID,
			UpdatedByID: event.UpdatedByID,
			CreatedBy:   ToUser(event.CreatedBy),
			UpdatedBy:   ToUser(event.UpdatedBy),
		},
		Thumbnail:          event.Thumbnail,
		Title:              event.Title,
		Slug:               event.Slug,
		Content:            event.Content,
		Excerpt:            event.Excerpt,
		IsActive:           event.IsActive,
		PublishedAt:        event.PublishedAt,
		IsNotificationSent: event.IsNotificationSent,
		Branches:           ToBranches(event.Branches),
	}
}

func ToEvents(events []*entity.EventEntity) []*model.Event {
	if len(events) == 0 {
		return make([]*model.Event, 0)
	}

	eventsRes := make([]*model.Event, 0, len(events))
	for _, event := range events {
		eventsRes = append(eventsRes, ToEvent(event))
	}

	return eventsRes
}
